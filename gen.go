package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"go/types"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/packages"
)

var fset = new(token.FileSet)

func gen(name string) (code []byte, err error) {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Fset: fset,
	}
	pkgs, err := packages.Load(cfg, name)
	if err != nil {
		return nil, err
	}

	if len(pkgs[0].Syntax) != 1 {
		// TODO(kr): remove this limitation
		return nil, errors.New("genbolt: schema must be a single file")
	}

	sch := &schema{
		Imports:      make(map[string]string),
		Keys:         make(map[string]bool),
		MapTypes:     make(map[string]bool),
		SeqTypes:     make(map[string]bool),
		JSONMapTypes: make(map[string]*types.Pointer),
		JSONSeqTypes: make(map[string]*types.Pointer),
		funcs:        make(template.FuncMap),
	}
	err = genFile(sch, pkgs[0].Syntax[0], pkgs[0])
	if err != nil {
		return nil, err
	}
	sch.InputFile = name

	var b bytes.Buffer
	tmpl, err := template.New("").
		Funcs(sch.funcs).
		Funcs(template.FuncMap{
			"trimprefix": strings.TrimPrefix,
		}).
		Parse(schemaTemplate)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(&b, sch)
	if err != nil {
		return nil, err
	}
	return format.Source(b.Bytes())
}

func genFile(sch *schema, file *ast.File, schemaPkg *packages.Package) error {
	typesInfo := schemaPkg.TypesInfo

	scope := schemaPkg.Types.Scope()
	sch.funcs["typestring"] = func(v interface{}) string {
		if s, ok := v.(string); ok {
			return s
		}
		t := v.(types.Type)
		return types.TypeString(t, (*types.Package).Name)
	}
	jsonTypes := make(map[*types.Named]bool)
	sch.funcs["isjsontype"] = func(v interface{}) bool {
		p, ok := v.(*types.Pointer)
		if !ok {
			return false
		}
		t, _ := p.Elem().(*types.Named)
		return jsonTypes[t]
	}
	sch.Package = schemaPkg.Types.Name()

	var userImports []*types.Package
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			vs := spec.(*ast.ValueSpec)
			iface := typesInfo.Types[vs.Type].Type
			if iface == nil {
				return fmt.Errorf("interface assertion has no interface: %v", esprint(vs))
			}
			if n, ok := iface.(*types.Named); !ok || n.String() != "encoding/json.Marshaler" {
				return fmt.Errorf("unsupported interface: %v", n)
			}

			for i, value := range vs.Values {
				if vs.Names[i].Name != "_" {
					return fmt.Errorf("interface assertion has non-_ name: %v", esprint(vs))
				}
				convType := typesInfo.Types[value].Type
				ptr, ok := convType.(*types.Pointer)
				if !ok {
					return fmt.Errorf("interface assertion has bad expression (must be pointer to named type): %v", esprint(convType))
				}
				named, ok := ptr.Elem().(*types.Named)
				if !ok {
					return fmt.Errorf("interface assertion has bad expression (must be pointer to named type): %v", esprint(convType))
				}
				jsonTypes[named] = true
				userImports = append(userImports, named.Obj().Pkg())
			}
		}
	}

	if len(jsonTypes) > 0 {
		sch.Imports["encoding/json"] = "json"
	}
	sch.Imports["encoding/binary"] = "binary"
	sch.Imports["github.com/coreos/bbolt"] = "bolt"
	seenNames := make(map[string]string) // name -> import path
	for _, pkg := range userImports {
		if p := seenNames[pkg.Name()]; p != "" && p != pkg.Path() {
			return fmt.Errorf("duplicate package name %s from import paths %s and %s", pkg.Name(), pkg.Path(), p)
		}
		seenNames[pkg.Name()] = pkg.Path()
		// TODO(kr): use package label from input file (also gives uniqueness)
		sch.Imports[pkg.Path()] = pkg.Name()
	}

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			return fmt.Errorf("unexpected decl: %v", esprint(decl))
		}
		switch genDecl.Tok {
		case token.VAR, token.IMPORT:
			continue
		case token.TYPE: // ok, proceed
		default:
			return fmt.Errorf("unexpected decl: %v", esprint(decl))
		}
		for _, spec := range genDecl.Specs {
			spec := spec.(*ast.TypeSpec)
			if spec.Assign != 0 {
				return fmt.Errorf("unexpected decl: %v", esprint(decl))
			}
			structType, ok := spec.Type.(*ast.StructType)
			if !ok {
				return fmt.Errorf("need struct type")
			}

			isRoot, rootSuffix := false, ""
			if s := spec.Name.Name; strings.HasPrefix(s, "Root") {
				rootSuffix = s[4:]
				isRoot = s == "Root" || ast.IsExported(rootSuffix)
			}

			sch.StructTypes = append(sch.StructTypes, &schemaStruct{
				Name:   spec.Name.Name,
				IsRoot: isRoot,
				Doc:    genDecl.Doc,
			})

			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					if !name.IsExported() {
						return fmt.Errorf("all fields must be exported")
					}
					if isReserved(name) {
						return fmt.Errorf("field name %s is reserved (sorry)", name.Name)
					}
					sch.Keys[name.Name] = true

					// TODO(kr): change this switch to use go/types
					switch fieldType := field.Type.(type) {
					case *ast.Ident:
						if !isBasic(scope, fieldType.Name) {
							return fmt.Errorf("unsupported type %s (try *%s instead?)", fieldType.Name, fieldType.Name)
						}
						if isRoot {
							return fmt.Errorf("unsupported root field type %s", esprint(field.Type))
						}
						sch.RecordFields = append(sch.RecordFields, &schemaField{
							Name:   name.Name,
							Type:   fieldType.Name,
							Bucket: spec.Name.Name,
							Doc:    field.Doc,
						})
					case *ast.StarExpr:
						ptr := typesInfo.Types[fieldType].Type.(*types.Pointer)
						named, ok := ptr.Elem().(*types.Named)
						if !ok {
							return fmt.Errorf("unknown type %s", esprint(field.Type))
						}

						if _, ok := jsonTypes[named]; ok {
							sch.RecordFields = append(sch.RecordFields, &schemaField{
								Name:   name.Name,
								Type:   types.NewPointer(named),
								Bucket: spec.Name.Name,
								Doc:    field.Doc,
							})
						} else if typ, ok := fieldType.X.(*ast.Ident); ok {
							// TODO(kr): look up typeName.Name and make sure it's a struct.
							sch.BucketFields = append(sch.BucketFields, &schemaField{
								Name:   name.Name,
								Type:   typ.Name,
								Bucket: spec.Name.Name,
								Doc:    field.Doc,
							})
						} else {
							return fmt.Errorf("unknown type %s", esprint(field.Type))
						}
					case *ast.ArrayType:
						if fieldType.Len != nil {
							return fmt.Errorf("cannot have array type (use a slice)")
						}

						switch elemType := fieldType.Elt.(type) {
						case *ast.Ident:
							if !isBasic(scope, elemType.Name) {
								return fmt.Errorf("unsupported type %s (try *%s instead?)", elemType.Name, elemType.Name)
							}
							sch.RecordFields = append(sch.RecordFields, &schemaField{
								Name:   name.Name,
								Type:   "[]" + elemType.Name,
								Bucket: spec.Name.Name,
								Doc:    field.Doc,
							})
						case *ast.StarExpr:
							switch typ := elemType.X.(type) {
							case *ast.Ident:
								// TODO(kr): look up typ.Name and make sure
								// it's a struct.
								sch.SeqTypes[typ.Name] = true
								sch.BucketFields = append(sch.BucketFields, &schemaField{
									Name:   name.Name,
									Type:   typ.Name + "Seq",
									Bucket: spec.Name.Name,
									Doc:    field.Doc,
								})
							default:
								t, _ := typesInfo.Types[typ].Type.(*types.Named)
								_, ok := jsonTypes[t]
								if !ok {
									return fmt.Errorf("cannot marshal %v; please implement json.Marshaler", esprint(elemType.X))
								}

								pkg := t.Obj().Pkg().Name()
								ru, n := utf8.DecodeRuneInString(pkg)
								seqType := string(unicode.ToUpper(ru)) + pkg[n:] + t.Obj().Name() + "Seq"
								sch.JSONSeqTypes[seqType] = types.NewPointer(t)

								sch.BucketFields = append(sch.BucketFields, &schemaField{
									Name:   name.Name,
									Type:   seqType,
									Bucket: spec.Name.Name,
									Doc:    field.Doc,
								})
							}
						default:
							return fmt.Errorf("slice element must be pointer to struct or basic type")
						}
					case *ast.MapType:
						keyType, ok := fieldType.Key.(*ast.Ident)
						if !ok || keyType.Name != "string" {
							return fmt.Errorf("map key must be string")
						}
						valueType, ok := fieldType.Value.(*ast.StarExpr)
						if !ok {
							return fmt.Errorf("map value must be pointer to named struct type")
						}

						switch typ := valueType.X.(type) {
						case *ast.Ident:
							sch.MapTypes[typ.Name] = true
							sch.BucketFields = append(sch.BucketFields, &schemaField{
								Name:   name.Name,
								Type:   typ.Name + "Map",
								Bucket: spec.Name.Name,
								Doc:    field.Doc,
							})
						default:
							t, _ := typesInfo.Types[typ].Type.(*types.Named)
							_, ok := jsonTypes[t]
							if !ok {
								return fmt.Errorf("cannot marshal %v; please implement json.Marshaler", esprint(valueType.X))
							}

							pkg := t.Obj().Pkg().Name()
							ru, n := utf8.DecodeRuneInString(pkg)
							mapType := string(unicode.ToUpper(ru)) + pkg[n:] + t.Obj().Name() + "Map"
							sch.JSONMapTypes[mapType] = types.NewPointer(t)

							sch.BucketFields = append(sch.BucketFields, &schemaField{
								Name:   name.Name,
								Type:   mapType,
								Bucket: spec.Name.Name,
								Doc:    field.Doc,
							})
						}
					default:
						return fmt.Errorf("unsupported type %v", esprint(field.Type))
					}
				}
			}
		}
	}
	return nil
}

func esprint(node interface{}) string {
	var b bytes.Buffer
	printer.Fprint(&b, fset, node)
	return b.String()
}

func isBasic(scope *types.Scope, name string) bool {
	obj := scope.Lookup(name)
	if obj == nil && scope == types.Universe {
		return false
	}
	if obj == nil {
		return isBasic(types.Universe, name)
	}
	_, ok := obj.Type().(*types.Basic)
	return ok
}

func isReserved(name *ast.Ident) bool {
	switch name.Name {
	case "Tx", "Bucket":
		return true
	}
	return false
}
