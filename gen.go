package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io"
	"strings"
	"text/template"

	"golang.org/x/tools/go/loader"
)

var prog *loader.Program

func gen(name string) (code []byte, err error) {
	cfg := loader.Config{
		ParserMode: parser.ParseComments,
	}
	cfg.CreateFromFilenames("", name)
	prog, err = cfg.Load()
	if err != nil {
		return nil, err
	}
	if n := len(prog.AllPackages); n != 1 {
		return nil, fmt.Errorf("got %d packages, want 1", n)
	}
	for _, pi := range prog.AllPackages {
		for _, file := range pi.Files {
			var b bytes.Buffer
			err = genFile(&b, file, pi.Pkg, prog)
			if err != nil {
				return nil, err
			}
			return format.Source(b.Bytes())
		}
	}
	panic("unreached")
}

type rootInfo struct {
	S string // name suffix, e.g. "Foo" in "RootFoo"
	C *ast.CommentGroup
}

type fieldInfo struct {
	B string // bucket (struct type name)
	F string // field name
	T string // type name of field
	C *ast.CommentGroup
}

func genFile(w io.Writer, file *ast.File, pkg *types.Package, prog *loader.Program) error {
	scope := pkg.Scope()
	keys := make(map[string]bool)
	fmt.Fprintln(w, "package", pkg.Name())
	fmt.Fprintln(w)
	fmt.Fprintln(w, imports)
	type container struct {
		Type string
		Elem string
	}
	mapTypes := make(map[string]*container)
	seqTypes := make(map[string]*container)
	needBucket := false
	needPut := false
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			return fmt.Errorf("unexpected decl: %v", decl)
		}
		if genDecl.Tok != token.TYPE {
			return fmt.Errorf("unexpected decl: %v", decl)
		}
		for _, spec := range genDecl.Specs {
			spec := spec.(*ast.TypeSpec)
			if spec.Assign != 0 {
				return fmt.Errorf("unexpected decl: %v", decl)
			}
			structType, ok := spec.Type.(*ast.StructType)
			if !ok {
				return fmt.Errorf("need struct type")
			}

			fmt.Fprintln(w)

			if s := spec.Name.Name; strings.HasPrefix(s, "Root") {
				suf := s[4:]
				if suf == "" || ast.IsExported(suf) {
					templateRoot.Execute(w, rootInfo{
						S: suf,
						C: genDecl.Doc,
					})

					for _, field := range structType.Fields.List {
						for _, name := range field.Names {
							if !name.IsExported() {
								return fmt.Errorf("all fields must be exported")
							}

							keys[name.Name] = true
							switch fieldType := field.Type.(type) {
							case *ast.StarExpr:
								typeName, ok := fieldType.X.(*ast.Ident)
								if !ok {
									return fmt.Errorf("cannot have pointer to non-struct type")
								}
								templatePointerField.Execute(w, fieldInfo{
									B: spec.Name.Name,
									F: name.Name,
									T: typeName.Name,
									C: field.Doc,
								})
								needBucket = true
							default:
								return fmt.Errorf("unsupported root field type %s", esprint(field.Type))
							}
						}
					}
					continue
				}
			}

			if genDecl.Doc != nil {
				for _, c := range genDecl.Doc.List {
					fmt.Fprintln(w, c.Text)
				}
			}
			fmt.Fprintln(w, "type", spec.Name, "struct {")
			fmt.Fprintln(w, "\tdb *bolt.Bucket")
			fmt.Fprintln(w, "}")

			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					if !name.IsExported() {
						return fmt.Errorf("all fields must be exported")
					}
					keys[name.Name] = true

					switch fieldType := field.Type.(type) {
					case *ast.Ident:
						if !isBasic(scope, fieldType.Name) {
							return fmt.Errorf("unsupported type %s (try *%s instead?)", fieldType.Name, fieldType.Name)
						}
						templateField.Execute(w, fieldInfo{
							B: spec.Name.Name,
							F: name.Name,
							T: fieldType.Name,
							C: field.Doc,
						})
						needPut = true
					case *ast.StarExpr:
						typeName, ok := fieldType.X.(*ast.Ident)
						if !ok {
							return fmt.Errorf("cannot have pointer to non-named type")
						}
						// TODO(kr): look up typeName.Name and make sure
						// it's a struct.
						templatePointerField.Execute(w, fieldInfo{
							B: spec.Name.Name,
							F: name.Name,
							T: typeName.Name,
							C: field.Doc,
						})
						needBucket = true
					case *ast.ArrayType:
						if fieldType.Len != nil {
							return fmt.Errorf("cannot have array type (use a slice)")
						}

						switch elemType := fieldType.Elt.(type) {
						case *ast.Ident:
							if !isBasic(scope, elemType.Name) {
								return fmt.Errorf("unsupported type %s (try *%s instead?)", elemType.Name, elemType.Name)
							}
							templateField.Execute(w, fieldInfo{
								B: spec.Name.Name,
								F: name.Name,
								T: "[]" + elemType.Name,
								C: field.Doc,
							})
							needPut = true
						case *ast.StarExpr:
							typeName, ok := elemType.X.(*ast.Ident)
							if !ok {
								return fmt.Errorf("cannot have pointer to non-named type")
							}
							// TODO(kr): look up typeName.Name and make sure
							// it's a struct.
							seqType := typeName.Name + "Seq"
							seq := &container{
								Type: seqType,
								Elem: typeName.Name,
							}
							if prev := seqTypes[seqType]; prev != nil && *prev != *seq {
								return fmt.Errorf("conflicting seq types map[string]%s and map[string]%s", prev.Elem, seq.Elem)
							}
							seqTypes[seqType] = seq
							templatePointerField.Execute(w, fieldInfo{
								B: spec.Name.Name,
								F: name.Name,
								T: seqType,
								C: field.Doc,
							})
							needBucket = true
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
						typeName, ok := valueType.X.(*ast.Ident)
						if !ok {
							return fmt.Errorf("map value must be pointer to named struct type")
						}
						mapType := typeName.Name + "Map"
						mp := &container{
							Type: mapType,
							Elem: typeName.Name,
						}
						if prev := mapTypes[mapType]; prev != nil && *prev != *mp {
							return fmt.Errorf("conflicting map types map[string]%s and map[string]%s", prev.Elem, mp.Elem)
						}
						mapTypes[mapType] = mp

						templatePointerField.Execute(w, fieldInfo{
							B: spec.Name.Name,
							F: name.Name,
							T: mapType,
							C: field.Doc,
						})
						needBucket = true
					default:
						return fmt.Errorf("unsupported type %v", esprint(field.Type))
					}
				}
			}
		}
	}
	for _, mp := range mapTypes {
		templateMapType.Execute(w, mp)
	}
	for _, seq := range seqTypes {
		templateSeqType.Execute(w, seq)
	}
	templateKeys.Execute(w, keys)
	if needBucket {
		fmt.Fprintln(w, bucket)
	}
	if needPut {
		fmt.Fprintln(w, put)
	}
	return nil
}

func esprint(node interface{}) string {
	var b bytes.Buffer
	printer.Fprint(&b, prog.Fset, node)
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

var tlib = template.Must(template.New("lib").Parse(`
{{define "get"}}
{{- if eq . "[]byte" -}}
	return v
{{- else if eq . "string" -}}
	return string(v)
{{- else if eq . "bool" -}}
	return v[0] != 0
{{- else if or (eq . "byte") (eq . "uint8") -}}
	return v[0]
{{- else if eq . "uint16" -}}
	return binary.BigEndian.Uint16(v)
{{- else if eq . "uint32" -}}
	return binary.BigEndian.Uint32(v)
{{- else if eq . "uint64" -}}
	return binary.BigEndian.Uint64(v)
{{- else if eq . "int8" -}}
	return int8(v[0])
{{- else if eq . "int16" -}}
	return int16(binary.BigEndian.Uint16(v))
{{- else if eq . "int32" -}}
	return int32(binary.BigEndian.Uint32(v))
{{- else if eq . "int64" -}}
	return int64(binary.BigEndian.Uint64(v))
{{- else -}}
	panic("internal error") {{- /* never generated */}}
{{- end -}}
{{end}}

{{define "put"}}
{{- if eq . "[]byte" -}}
	v := x
{{- else if eq . "string" -}}
	v := []byte(x)
{{- else if eq . "bool" -}}
	v := make([]byte, 1)
	if x { v[0] = 1 }
{{- else if or (eq . "byte") (eq . "uint8") -}}
	v := []byte{x}
{{- else if eq . "uint16" -}}
	v := make([]byte, 2)
	binary.BigEndian.PutUint16(x)
{{- else if eq . "uint32" -}}
	v := make([]byte, 4)
	binary.BigEndian.PutUint32(x)
{{- else if eq . "uint64" -}}
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(x)
{{- else if eq . "int8" -}}
	v := []byte{byte(x)}
{{- else if eq . "int16" -}}
	v := make([]byte, 2)
	binary.BigEndian.PutUint16(uint16(x))
{{- else if eq . "int32" -}}
	v := make([]byte, 4)
	binary.BigEndian.PutUint32(uint32(x))
{{- else if eq . "int64" -}}
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(uint64(x))
{{- else -}}
	panic("internal error") {{- /* never generated */}}
{{- end -}}
{{end}}
`))

var templateField = template.Must(template.Must(tlib.Clone()).Parse(`
{{if .C.Text -}}
{{range .C.List -}}
{{.Text}}
{{end -}}
{{end -}}
func (o *{{.B}}) {{.F}}() {{.T}} {
	v := o.db.Get(key{{.F}})
	{{template "get" .T}}
}

// Put{{.F}} stores x as the value of {{.F}}.
{{if .C.Text -}}
//
{{range .C.List -}}
{{.Text}}
{{end -}}
{{end -}}
func (o *{{.B}}) Put{{.F}}(x {{.T}}) {
	{{template "put" .T}}
	put(o.db, key{{.F}}, v)
}
`))

var templatePointerField = template.Must(template.Must(tlib.Clone()).Parse(`
{{if .C.Text -}}
{{range .C.List -}}
{{.Text}}
{{end -}}
{{end -}}
func (o *{{.B}}) {{.F}}() *{{.T}} {
	return &{{.T}}{bucket(o.db, key{{.F}})}
}
`))

var templateMapType = template.Must(template.Must(tlib.Clone()).Parse(`
type {{.Type}} struct {
	db *bolt.Bucket
}

func (o *{{.Type}}) Get(key []byte) *{{.Elem}} {
	return &{{.Elem}}{bucket(o.db, key)}
}

func (o *{{.Type}}) GetByString(key string) *{{.Elem}} {
	{{/* TODO(kr): consider unsafe conversion */ -}}
	return &{{.Elem}}{bucket(o.db, []byte(key))}
}
`))

var templateSeqType = template.Must(template.Must(tlib.Clone()).Parse(`
type {{.Type}} struct {
	db *bolt.Bucket
}

func (o *{{.Type}}) Get(n uint64) *{{.Elem}} {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(n)
	return &{{.Elem}}{bucket(o.db, key)}
}

func (o *{{.Type}}) Add() *{{.Elem}} {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	return o.Get(n)
}
`))

var templateKeys = template.Must(template.Must(tlib.Clone()).Parse(`
var (
{{- range $name, $_ := .}}
	key{{$name}} = []byte({{printf "%q" $name}})
{{- end}}
)
`))

var templateRoot = template.Must(template.Must(tlib.Clone()).Parse(`
{{if .C.Text -}}
{{range .C.List -}}
{{.Text}}
{{end -}}
{{end -}}
type Root{{.S}} struct {
	db *bolt.Tx
}

// NewRoot{{.S}} returns a new Root{{.S}} for tx.
{{if .C.Text -}}
//
{{range .C.List -}}
{{.Text}}
{{end -}}
{{end -}}
func NewRoot{{.S}}(tx *bolt.Tx) *Root{{.S}} {
	return &Root{{.S}}{tx}
}

func View{{.S}}(db *bolt.DB, f func(*Root{{.S}}, *bolt.Tx) error) error {
	return db.View(func(tx *bolt.Tx) error {
		return f(&Root{{.S}}{tx}, tx)
	})
}

func Update{{.S}}(db *bolt.DB, f func(*Root{{.S}}, *bolt.Tx) error) error {
	return db.Update(func(tx *bolt.Tx) error {
		return f(&Root{{.S}}{tx}, tx)
	})
}
`))

const imports = `
import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize
`

const bucket = `
type db interface {
	Writable() bool
	CreateBucketIfNotExists([]byte) *bolt.Bucket
	Bucket([]byte) *bolt.Bucket
}

func bucket(db db, key []byte) *bolt.Bucket {
	if db.Writable() {
		return db.CreateBucketIfNotExists(key)
	} else {
		return db.Bucket(key)
	}
}
`

const put = `
func put(b *bolt.Bucket, key, value []byte) {
	err := b.Put(key, value)
	if err != nil {
		panic(err)
	}
}
`
