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
	"strconv"
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

	ctx := &context{
		pkg: pkgs[0],
		sch: &schema{
			Imports:      make(map[string]string),
			Keys:         make(map[string]bool),
			MapTypes:     make(map[string]bool),
			SeqTypes:     make(map[string]bool),
			JSONMapTypes: make(map[string]*types.Pointer),
			JSONSeqTypes: make(map[string]*types.Pointer),
			funcs:        make(template.FuncMap),
		},
		jsonTypes: make(map[*types.Named]bool),
	}
	err = genFile(ctx, pkgs[0].Syntax[0])
	if err != nil {
		return nil, err
	}
	ctx.sch.InputFile = name

	var b bytes.Buffer
	tmpl, err := template.New("").
		Funcs(ctx.sch.funcs).
		Funcs(template.FuncMap{
			"trimprefix": strings.TrimPrefix,
			"identical":  types.Identical,
			"basic":      basicType,
			"sliceof":    types.NewSlice,
		}).
		Parse(schemaTemplate)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(&b, ctx.sch)
	if err != nil {
		return nil, err
	}
	code, err = format.Source(b.Bytes())
	if err != nil {
		return b.Bytes(), err
	}
	return code, nil
}

type context struct {
	sch *schema
	pkg *packages.Package

	jsonTypes map[*types.Named]bool
}

func genFile(ctx *context, file *ast.File) error {
	sch := ctx.sch
	typesInfo := ctx.pkg.TypesInfo

	sch.Package = file.Name.Name

	sch.funcs["typestring"] = func(t types.Type) string {
		return types.TypeString(t, func(p *types.Package) string {
			return sch.Imports[p.Path()]
		})
	}
	sch.funcs["isjsontype"] = func(v interface{}) bool {
		p, ok := v.(*types.Pointer)
		if !ok {
			return false
		}
		t, _ := p.Elem().(*types.Named)
		return ctx.jsonTypes[t]
	}

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
				ctx.jsonTypes[named] = true
			}
		}
	}

	for _, imp := range ctx.pkg.Types.Imports() {
		sch.Imports[imp.Path()] = imp.Name()
	}
	for _, imp := range file.Imports {
		path, _ := strconv.Unquote(imp.Path.Value)
		if imp.Name != nil {
			sch.Imports[path] = imp.Name.Name
		}
	}

	if len(ctx.jsonTypes) > 0 {
		sch.Imports["encoding/json"] = "json"
	}
	sch.Imports["encoding/binary"] = "binary"
	sch.Imports["github.com/coreos/bbolt"] = "bolt"

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			return fmt.Errorf("unexpected decl: %v", esprint(decl))
		}
		switch genDecl.Tok {
		default:
			return fmt.Errorf("unexpected decl: %v", esprint(decl))
		case token.VAR, token.IMPORT:
			continue
		case token.TYPE: // ok, proceed
		}
		for _, spec := range genDecl.Specs {
			spec := spec.(*ast.TypeSpec)
			if spec.Assign != 0 {
				return fmt.Errorf("unexpected decl: %v", esprint(decl))
			}
			doc := spec.Doc
			if doc == nil {
				doc = genDecl.Doc
			}
			err := genStruct(ctx, spec.Name.Name, spec.Type, doc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genStruct(ctx *context, name string, typ ast.Expr, doc *ast.CommentGroup) error {
	sch := ctx.sch

	structType, ok := typ.(*ast.StructType)
	if !ok {
		return fmt.Errorf("need struct type")
	}

	isRoot := false
	if strings.HasPrefix(name, "Root") {
		isRoot = name == "Root" || ast.IsExported(name[4:])
	}

	sch.StructTypes = append(sch.StructTypes, &schemaStruct{
		Name:   name,
		IsRoot: isRoot,
		Doc:    doc,
	})

	for _, field := range structType.Fields.List {
		for _, fieldIdent := range field.Names {
			if !fieldIdent.IsExported() {
				return fmt.Errorf("all fields must be exported")
			}
			if isReserved(fieldIdent.Name) {
				return fmt.Errorf("field name %s is reserved (sorry)", fieldIdent)
			}
			sch.Keys[fieldIdent.Name] = true

			switch fieldType := ctx.pkg.TypesInfo.Defs[fieldIdent].Type().(type) {
			default:
				return fmt.Errorf("unsupported type %v", esprint(field.Type))
			case *types.Basic:
				if isRoot {
					return fmt.Errorf("unsupported root field type %s", esprint(field.Type))
				}
				sch.RecordFields = append(sch.RecordFields, &schemaField{
					Name:   fieldIdent.Name,
					Type:   fieldType,
					Bucket: name,
					Doc:    field.Doc,
				})
			case *types.Named:
				return fmt.Errorf("unsupported type %s (try *%s instead?)", fieldIdent.Name, fieldType.String())
			case *types.Pointer:
				named, ok := fieldType.Elem().(*types.Named)
				if !ok {
					return fmt.Errorf("unknown type %s", esprint(field.Type))
				}

				if _, ok := ctx.jsonTypes[named]; ok {
					sch.RecordFields = append(sch.RecordFields, &schemaField{
						Name:   fieldIdent.Name,
						Type:   fieldType,
						Bucket: name,
						Doc:    field.Doc,
					})
				} else if _, ok := named.Underlying().(*types.Struct); ok {
					sch.BucketFields = append(sch.BucketFields, &schemaField{
						Name:   fieldIdent.Name,
						Type:   fieldType,
						Bucket: name,
						Doc:    field.Doc,
					})
				} else {
					return fmt.Errorf("unknown type %s", esprint(field.Type))
				}
			case *types.Array:
				return fmt.Errorf("cannot have array type (use a slice)")
			case *types.Slice:
				switch elemType := fieldType.Elem().(type) {
				default:
					return fmt.Errorf("slice element must be pointer to struct or basic type")
				case *types.Basic:
					sch.RecordFields = append(sch.RecordFields, &schemaField{
						Name:   fieldIdent.Name,
						Type:   fieldType,
						Bucket: name,
						Doc:    field.Doc,
					})
				case *types.Pointer:
					named, ok := elemType.Elem().(*types.Named)
					if !ok {
						return fmt.Errorf("unknown type %s", esprint(field.Type))
					}

					var seqTypeName string

					if _, ok := ctx.jsonTypes[named]; ok {
						pkgName := named.Obj().Pkg().Name()
						ru, n := utf8.DecodeRuneInString(pkgName)
						seqTypeName = string(unicode.ToUpper(ru)) + pkgName[n:] + named.Obj().Name() + "Seq"
						sch.JSONSeqTypes[seqTypeName] = elemType
					} else if _, ok := named.Underlying().(*types.Struct); ok {
						seqTypeName = named.Obj().Name() + "Seq"
						sch.SeqTypes[named.Obj().Name()] = true
					} else {
						return fmt.Errorf("unknown type %s", esprint(field.Type))
					}

					sch.BucketFields = append(sch.BucketFields, &schemaField{
						Name: fieldIdent.Name,
						Type: types.NewPointer(types.NewNamed(
							types.NewTypeName(0, ctx.pkg.Types, seqTypeName, nil),
							types.Typ[types.Invalid],
							nil,
						)),
						Bucket: name,
						Doc:    field.Doc,
					})
				}
			case *types.Map:
				// TODO(kr): allow numeric types as map keys too
				keyType, ok := fieldType.Key().(*types.Basic)
				if !ok || keyType.Kind() != types.String {
					return fmt.Errorf("map key must be string")
				}
				ptr, ok := fieldType.Elem().(*types.Pointer)
				if !ok {
					return fmt.Errorf("map value must be pointer to named struct type")
				}

				named, ok := ptr.Elem().(*types.Named)
				if !ok {
					return fmt.Errorf("unknown type %s", esprint(field.Type))
				}

				var mapTypeName string

				if _, ok := ctx.jsonTypes[named]; ok {
					pkgName := named.Obj().Pkg().Name()
					ru, n := utf8.DecodeRuneInString(pkgName)
					mapTypeName = string(unicode.ToUpper(ru)) + pkgName[n:] + named.Obj().Name() + "Map"
					sch.JSONMapTypes[mapTypeName] = types.NewPointer(named)
				} else if _, ok := named.Underlying().(*types.Struct); ok {
					mapTypeName = named.Obj().Name() + "Map"
					sch.MapTypes[named.Obj().Name()] = true
				} else {
					return fmt.Errorf("unknown type %s", esprint(field.Type))
				}

				sch.BucketFields = append(sch.BucketFields, &schemaField{
					Name: fieldIdent.Name,
					Type: types.NewPointer(types.NewNamed(
						types.NewTypeName(0, ctx.pkg.Types, mapTypeName, nil),
						types.Typ[types.Invalid],
						nil,
					)),
					Bucket: name,
					Doc:    field.Doc,
				})
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

func isReserved(name string) bool {
	switch name {
	case "Tx", "Bucket":
		return true
	}
	return false
}

var basicTypes = make(map[string]*types.Basic)

// basicType returns the named basic type.
func basicType(name string) *types.Basic {
	return basicTypes[name]
}

func init() {
	for _, t := range types.Typ {
		basicTypes[t.Name()] = t
	}
}
