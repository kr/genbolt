package main

const schemaTemplate = `
{{- define "get" -}}
{{- if isjsontype . -}}
	v := new({{typestring .Elem}})
	if rec == nil {
		return v
	}
	err := json.Unmarshal(rec, json.Unmarshaler(v))
	if err != nil {
		panic(err)
	}
	return v
{{- else if isbintype . -}}
	v := new({{typestring .Elem}})
	if rec == nil {
		return v
	}
	err := encoding.BinaryUnmarshaler(v).UnmarshalBinary(rec)
	if err != nil {
		panic(err)
	}
	return v
{{- else if identical . (sliceof (basic "uint8")) -}}
	v := make([]byte, len(rec))
	copy(v, rec)
	return v
{{- else if identical . (basic "string") -}}
	return string(rec)
{{- else if identical . (basic "bool") -}}
	return len(rec) > 0 && rec[0] != 0
{{- else if identical . (basic "uint8") -}}
	if rec == nil {
		return 0
	}
	return rec[0]
{{- else if identical . (basic "uint16") -}}
	if rec == nil {
		return 0
	}
	return binary.BigEndian.Uint16(rec)
{{- else if identical . (basic "uint32") -}}
	if rec == nil {
		return 0
	}
	return binary.BigEndian.Uint32(rec)
{{- else if identical . (basic "uint64") -}}
	if rec == nil {
		return 0
	}
	return binary.BigEndian.Uint64(rec)
{{- else if identical . (basic "int8") -}}
	if rec == nil {
		return 0
	}
	return int8(rec[0])
{{- else if identical . (basic "int16") -}}
	if rec == nil {
		return 0
	}
	return int16(binary.BigEndian.Uint16(rec))
{{- else if identical . (basic "int32") -}}
	if rec == nil {
		return 0
	}
	return int32(binary.BigEndian.Uint32(rec))
{{- else if identical . (basic "int64") -}}
	if rec == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(rec))
{{- else -}}
	panic("internal error") {{- /* never generated */}}
{{- end -}}
{{- end -}}

{{- define "put" -}}
{{- if isjsontype . -}}
	rec, err := json.Marshal(json.Marshaler(v))
	if err != nil {
		panic(err)
	}
{{- else if isbintype . -}}
	rec, err := encoding.BinaryMarshaler(v).MarshalBinary()
	if err != nil {
		panic(err)
	}
{{- else if identical . (sliceof (basic "uint8")) -}}
	rec := v
{{- else if identical . (basic "string") -}}
	rec := []byte(v)
{{- else if identical . (basic "bool") -}}
	rec := make([]byte, 1)
	if v { rec[0] = 1 }
{{- else if identical . (basic "uint8") -}}
	rec := []byte{v}
{{- else if identical . (basic "uint16") -}}
	rec := make([]byte, 2)
	binary.BigEndian.PutUint16(rec, v)
{{- else if identical . (basic "uint32") -}}
	rec := make([]byte, 4)
	binary.BigEndian.PutUint32(rec, v)
{{- else if identical . (basic "uint64") -}}
	rec := make([]byte, 8)
	binary.BigEndian.PutUint64(rec, v)
{{- else if identical . (basic "int8") -}}
	rec := []byte{byte(v)}
{{- else if identical . (basic "int16") -}}
	rec := make([]byte, 2)
	binary.BigEndian.PutUint16(rec, uint16(v))
{{- else if identical . (basic "int32") -}}
	rec := make([]byte, 4)
	binary.BigEndian.PutUint32(rec, uint32(v))
{{- else if identical . (basic "int64") -}}
	rec := make([]byte, 8)
	binary.BigEndian.PutUint64(rec, uint64(v))
{{- else -}}
	panic("internal error") {{- /* never generated */}}
{{- end -}}
{{- end -}}

// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See {{.InputFile}}.

package {{.Package}}

{{- range $path, $name := .Imports}}
import {{$name}} {{printf "%q" $path}}
{{- end}}

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

{{range .StructTypes}}
{{ $level := or (and .IsRoot "Tx") "Bucket" }}

// {{.Name}} is a bucket with a static set of elements.
{{if .Doc -}}
//
{{range .Doc.List -}}
{{.Text}}
{{end -}}
//
{{end -}}
// Accessor methods read and write records
// and open child buckets.
{{if .IsRoot -}}
// See functions View{{trimprefix .Name "Root"}} and
// Update{{trimprefix .Name "Root"}} to open transactions.
{{end -}}
type {{.Name}} struct {
	db *bolt.{{$level}}
}

{{if .IsRoot}}
// New{{.Name}} returns a new {{.Name}} for tx.
{{if .Doc -}}
//
{{range .Doc.List -}}
{{.Text}}
{{end -}}
{{end -}}
func New{{.Name}}(tx *bolt.Tx) *{{.Name}} {
	return &{{.Name}}{tx}
}

// View{{trimprefix .Name "Root"}} opens a read-only transaction
// and calls f with an instance of {{.Name}} as the root bucket.
// It returns the error returned by f.
func View{{trimprefix .Name "Root"}}(db *bolt.DB, f func(*{{.Name}}) error) error {
	return db.View(func(tx *bolt.Tx) error {
		return f(&{{.Name}}{tx})
	})
}

// Update{{trimprefix .Name "Root"}} opens a writable transaction
// and calls f with an instance of {{.Name}} as the root bucket,
// then it commits the transaction.
// It returns the error returned by f,
// or any error committing to the database, if f was successful.
func Update{{trimprefix .Name "Root"}}(db *bolt.DB, f func(*{{.Name}}) error) error {
	return db.Update(func(tx *bolt.Tx) error {
		return f(&{{.Name}}{tx})
	})
}
{{end}}

// {{$level}} returns o's underlying *bolt.{{$level}} object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
{{if not .IsRoot -}}
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
{{end -}}
func (o *{{.Name}}) {{$level}}() *bolt.{{$level}} {
	return o.db
}
{{end}}

{{range .BucketFields}}
// {{.Name}} gets the child bucket with key {{printf "%q" .Name}} from o.
{{if .Doc -}}
//
{{range .Doc.List -}}
{{.Text}}
{{end -}}
{{end -}}
//
// {{.Name}} creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil {{typestring .Type}};
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *{{.Bucket}}) {{.Name}}() {{typestring .Type}} {
	return &{{typestring .Type.Elem}}{bucket(o.db, key{{.Name}})}
}
{{end}}

{{range .RecordFields}}
// {{.Name}} reads the record stored under key {{printf "%q" .Name}}.
{{if .Doc -}}
//
{{range .Doc.List -}}
{{.Text}}
{{end -}}
//
{{end -}}
// If no record has been stored, {{.Name}} returns
{{if ispointer .Type -}}
// a pointer to
{{end -}}
// the zero value.
func (o *{{.Bucket}}) {{.Name}}() {{typestring .Type}} {
	rec := get(o.db, key{{.Name}})
	{{template "get" .Type}}
}

// Put{{.Name}} stores v as a record under the key {{printf "%q" .Name}}.
{{if .Doc -}}
//
{{range .Doc.List -}}
{{.Text}}
{{end -}}
{{end -}}
func (o *{{.Bucket}}) Put{{.Name}}(v {{typestring .Type}}) {
	{{template "put" .Type}}
	put(o.db, key{{.Name}}, rec)
}
{{end}}

{{range $type, $elem := .MapOfBucketTypes}}
// {{$type}} is a bucket with arbitrary keys,
// holding child buckets of type {{$elem}}.
type {{$type}} struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *{{$type}}) Bucket() *bolt.Bucket {
	return o.db
}

// Get gets the child bucket with the given key from o.
//
// It creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *{{$elem}};
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *{{$type}}) Get(key []byte) *{{$elem}} {
	return &{{$elem}}{bucket(o.db, key)}
}

// GetByString is equivalent to o.Get([]byte(key)).
func (o *{{$type}}) GetByString(key string) *{{$elem}} {
	{{/* TODO(kr): consider unsafe conversion */ -}}
	return &{{$elem}}{bucket(o.db, []byte(key))}
}
{{end}}

{{range $type, $elem := .SeqOfBucketTypes}}
// {{$type}} is a bucket with sequential numeric keys,
// holding child buckets of type {{$elem}}.
type {{$type}} struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *{{$type}}) Bucket() *bolt.Bucket {
	return o.db
}

// Get gets child bucket n from o.
//
// It creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *{{$elem}};
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *{{$type}}) Get(n uint64) *{{$elem}} {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	return &{{$elem}}{bucket(o.db, key)}
}

// Add creates and returns a new, empty child bucket to o
// with a new sequence number.
//
// It panics if called in a read-only transaction.
func (o *{{$type}}) Add() (*{{$elem}}, uint64) {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	return o.Get(n), n
}
{{end}}

{{range $type, $elem := .MapOfRecordTypes}}
// {{$type}} is a bucket with arbitrary keys,
// holding records of type {{typestring $elem}}.
type {{$type}} struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *{{$type}}) Bucket() *bolt.Bucket {
	return o.db
}

// Get reads the record stored in o under the given key.
//
// If no record has been stored, it returns
// a pointer to the zero value.
func (o *{{$type}}) Get(key []byte) {{typestring $elem}} {
	rec := get(o.db, key)
	{{template "get" $elem}}
}

// GetByString is equivalent to o.Get([]byte(key)).
func (o *{{$type}}) GetByString(key string) {{typestring $elem}} {
	{{/* TODO(kr): consider unsafe conversion */ -}}
	return o.Get([]byte(key))
}

// Put stores v in o as a record under the given key.
func (o *{{$type}}) Put(key []byte, v {{typestring $elem}}) {
	{{template "put" $elem}}
	put(o.db, key, rec)
}

// PutByString is equivalent to o.Put([]byte(key), v).
func (o *{{$type}}) PutByString(key string, v {{typestring $elem}}) {
	{{/* TODO(kr): consider unsafe conversion */ -}}
	o.Put([]byte(key), v)
}
{{end}}

{{range $type, $elem := .SeqOfRecordTypes}}
// {{$type}} is a bucket with sequential numeric keys,
// holding records of type {{typestring $elem}}.
type {{$type}} struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *{{$type}}) Bucket() *bolt.Bucket {
	return o.db
}

// Get reads the record stored in o under sequence number n.
//
// If no record has been stored, it returns
// a pointer to the zero value.
func (o *{{$type}}) Get(n uint64) {{typestring $elem}} {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	rec := get(o.db, key)
	{{template "get" $elem}}
}

// Add stores v in o under a new sequence number.
// It writes the new sequence number to *np
// before marshaling v. It is okay for
// np to point to a field inside v, to store
// the sequence number in the new record.
func (o *{{$type}}) Add(v {{typestring $elem}}, np *uint64) {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	*np = n
	o.Put(n, v)
}

// Put stores v in o as a record under sequence number n.
func (o *{{$type}}) Put(n uint64, v {{typestring $elem}}) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	{{template "put" $elem}}
	put(o.db, key, rec)
}
{{end}}

var (
	{{- range $key, $_ := .Keys}}
	key{{$key}} = []byte({{printf "%q" $key}})
	{{- end}}
)

{{if .BucketFields}}
type db interface {
	Writable() bool
	CreateBucketIfNotExists([]byte) (*bolt.Bucket, error)
	Bucket([]byte) *bolt.Bucket
}

func bucket(db db, key []byte) *bolt.Bucket {
	if bu, ok := db.(*bolt.Bucket); ok && bu == nil {
		return nil // can happen in read-only txs
	}
	if !db.Writable() {
		return db.Bucket(key)
	}
	b, err := db.CreateBucketIfNotExists(key)
	if err != nil {
		panic(err)
	}
	return b
}
{{end}}

{{if or .RecordFields .MapOfRecordTypes .SeqOfRecordTypes}}
func get(b *bolt.Bucket, key []byte) []byte {
	if b == nil {
		return nil
	}
	return b.Get(key)
}

func put(b *bolt.Bucket, key, value []byte) {
	err := b.Put(key, value)
	if err != nil {
		panic(err)
	}
}
{{end}}
`
