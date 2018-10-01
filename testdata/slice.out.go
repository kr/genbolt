// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/slice.in.go.

package db

import bytes "bytes"
import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize
const _ = bytes.MinRead

// T is a bucket with a static set of elements.
// Accessor methods read and write records
// and open child buckets.
type T struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *T) Bucket() *bolt.Bucket {
	return o.db
}

// Strings gets the child bucket with key "Strings" from o.
//
// Strings don't have a fixed size, so we don't
// special-case them like we do for the other
// basic types. So this is a SeqOfString, not a
// SliceOfString.
//
// Strings creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *SeqOfString;
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *T) Strings() *SeqOfString {
	return &SeqOfString{bucket(o.db, keyStrings)}
}

// Int8s reads the record stored under key "Int8s".
// If no record has been stored, Int8s returns
// the zero value.
func (o *T) Int8s() []int8 {
	rec := get(o.db, keyInt8s)
	b := bytes.NewReader(rec)
	v := make([]int8, len(rec)/1)
	err := binary.Read(b, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// PutInt8s stores v as a record under the key "Int8s".
func (o *T) PutInt8s(v []int8) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	rec := b.Bytes()
	put(o.db, keyInt8s, rec)
}

// Int16s reads the record stored under key "Int16s".
// If no record has been stored, Int16s returns
// the zero value.
func (o *T) Int16s() []int16 {
	rec := get(o.db, keyInt16s)
	b := bytes.NewReader(rec)
	v := make([]int16, len(rec)/2)
	err := binary.Read(b, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// PutInt16s stores v as a record under the key "Int16s".
func (o *T) PutInt16s(v []int16) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	rec := b.Bytes()
	put(o.db, keyInt16s, rec)
}

// Int32s reads the record stored under key "Int32s".
// If no record has been stored, Int32s returns
// the zero value.
func (o *T) Int32s() []int32 {
	rec := get(o.db, keyInt32s)
	b := bytes.NewReader(rec)
	v := make([]int32, len(rec)/4)
	err := binary.Read(b, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// PutInt32s stores v as a record under the key "Int32s".
func (o *T) PutInt32s(v []int32) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	rec := b.Bytes()
	put(o.db, keyInt32s, rec)
}

// Int64s reads the record stored under key "Int64s".
// If no record has been stored, Int64s returns
// the zero value.
func (o *T) Int64s() []int64 {
	rec := get(o.db, keyInt64s)
	b := bytes.NewReader(rec)
	v := make([]int64, len(rec)/8)
	err := binary.Read(b, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// PutInt64s stores v as a record under the key "Int64s".
func (o *T) PutInt64s(v []int64) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	rec := b.Bytes()
	put(o.db, keyInt64s, rec)
}

// Uint8s reads the record stored under key "Uint8s".
// If no record has been stored, Uint8s returns
// the zero value.
func (o *T) Uint8s() []uint8 {
	rec := get(o.db, keyUint8s)
	v := make([]byte, len(rec))
	copy(v, rec)
	return v
}

// PutUint8s stores v as a record under the key "Uint8s".
func (o *T) PutUint8s(v []uint8) {
	rec := v
	put(o.db, keyUint8s, rec)
}

// Uint16s reads the record stored under key "Uint16s".
// If no record has been stored, Uint16s returns
// the zero value.
func (o *T) Uint16s() []uint16 {
	rec := get(o.db, keyUint16s)
	b := bytes.NewReader(rec)
	v := make([]uint16, len(rec)/2)
	err := binary.Read(b, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// PutUint16s stores v as a record under the key "Uint16s".
func (o *T) PutUint16s(v []uint16) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	rec := b.Bytes()
	put(o.db, keyUint16s, rec)
}

// Uint32s reads the record stored under key "Uint32s".
// If no record has been stored, Uint32s returns
// the zero value.
func (o *T) Uint32s() []uint32 {
	rec := get(o.db, keyUint32s)
	b := bytes.NewReader(rec)
	v := make([]uint32, len(rec)/4)
	err := binary.Read(b, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// PutUint32s stores v as a record under the key "Uint32s".
func (o *T) PutUint32s(v []uint32) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	rec := b.Bytes()
	put(o.db, keyUint32s, rec)
}

// Uint64s reads the record stored under key "Uint64s".
// If no record has been stored, Uint64s returns
// the zero value.
func (o *T) Uint64s() []uint64 {
	rec := get(o.db, keyUint64s)
	b := bytes.NewReader(rec)
	v := make([]uint64, len(rec)/8)
	err := binary.Read(b, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// PutUint64s stores v as a record under the key "Uint64s".
func (o *T) PutUint64s(v []uint64) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	rec := b.Bytes()
	put(o.db, keyUint64s, rec)
}

// Bools reads the record stored under key "Bools".
// If no record has been stored, Bools returns
// the zero value.
func (o *T) Bools() []bool {
	rec := get(o.db, keyBools)
	b := bytes.NewReader(rec)
	v := make([]bool, len(rec)/1)
	err := binary.Read(b, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// PutBools stores v as a record under the key "Bools".
func (o *T) PutBools(v []bool) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.BigEndian, v)
	if err != nil {
		panic(err)
	}
	rec := b.Bytes()
	put(o.db, keyBools, rec)
}

// SeqOfString is a bucket with sequential numeric keys,
// holding records of type string.
type SeqOfString struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *SeqOfString) Bucket() *bolt.Bucket {
	return o.db
}

// Get reads the record stored in o under sequence number n.
//
// If no record has been stored, it returns
// the zero value.
func (o *SeqOfString) Get(n uint64) string {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	rec := get(o.db, key)
	return string(rec)
}

// Add stores v in o under a new sequence number.
// It writes the new sequence number to *np
// before marshaling v. It is okay for
// np to point to a field inside v, to store
// the sequence number in the new record.
func (o *SeqOfString) Add(v string, np *uint64) {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	*np = n
	o.Put(n, v)
}

// Put stores v in o as a record under sequence number n.
func (o *SeqOfString) Put(n uint64, v string) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	rec := []byte(v)
	put(o.db, key, rec)
}

var (
	keyBools   = []byte("Bools")
	keyInt16s  = []byte("Int16s")
	keyInt32s  = []byte("Int32s")
	keyInt64s  = []byte("Int64s")
	keyInt8s   = []byte("Int8s")
	keyStrings = []byte("Strings")
	keyUint16s = []byte("Uint16s")
	keyUint32s = []byte("Uint32s")
	keyUint64s = []byte("Uint64s")
	keyUint8s  = []byte("Uint8s")
)

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
