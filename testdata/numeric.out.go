// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/numeric.in.go.

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

// Bool reads the record stored under key "Bool".
// If no record has been stored, Bool returns
// the zero value.
func (o *T) Bool() bool {
	rec := get(o.db, keyBool)
	return len(rec) > 0 && rec[0] != 0
}

// PutBool stores v as a record under the key "Bool".
func (o *T) PutBool(v bool) {
	rec := make([]byte, 1)
	if v {
		rec[0] = 1
	}
	put(o.db, keyBool, rec)
}

// Byte reads the record stored under key "Byte".
// If no record has been stored, Byte returns
// the zero value.
func (o *T) Byte() byte {
	rec := get(o.db, keyByte)
	if rec == nil {
		return 0
	}
	return rec[0]
}

// PutByte stores v as a record under the key "Byte".
func (o *T) PutByte(v byte) {
	rec := []byte{v}
	put(o.db, keyByte, rec)
}

// Uint16 reads the record stored under key "Uint16".
// If no record has been stored, Uint16 returns
// the zero value.
func (o *T) Uint16() uint16 {
	rec := get(o.db, keyUint16)
	if rec == nil {
		return 0
	}
	return binary.BigEndian.Uint16(rec)
}

// PutUint16 stores v as a record under the key "Uint16".
func (o *T) PutUint16(v uint16) {
	rec := make([]byte, 2)
	binary.BigEndian.PutUint16(rec, v)
	put(o.db, keyUint16, rec)
}

// Uint32 reads the record stored under key "Uint32".
// If no record has been stored, Uint32 returns
// the zero value.
func (o *T) Uint32() uint32 {
	rec := get(o.db, keyUint32)
	if rec == nil {
		return 0
	}
	return binary.BigEndian.Uint32(rec)
}

// PutUint32 stores v as a record under the key "Uint32".
func (o *T) PutUint32(v uint32) {
	rec := make([]byte, 4)
	binary.BigEndian.PutUint32(rec, v)
	put(o.db, keyUint32, rec)
}

// Uint64 reads the record stored under key "Uint64".
// If no record has been stored, Uint64 returns
// the zero value.
func (o *T) Uint64() uint64 {
	rec := get(o.db, keyUint64)
	if rec == nil {
		return 0
	}
	return binary.BigEndian.Uint64(rec)
}

// PutUint64 stores v as a record under the key "Uint64".
func (o *T) PutUint64(v uint64) {
	rec := make([]byte, 8)
	binary.BigEndian.PutUint64(rec, v)
	put(o.db, keyUint64, rec)
}

// Int8 reads the record stored under key "Int8".
// If no record has been stored, Int8 returns
// the zero value.
func (o *T) Int8() int8 {
	rec := get(o.db, keyInt8)
	if rec == nil {
		return 0
	}
	return int8(rec[0])
}

// PutInt8 stores v as a record under the key "Int8".
func (o *T) PutInt8(v int8) {
	rec := []byte{byte(v)}
	put(o.db, keyInt8, rec)
}

// Int16 reads the record stored under key "Int16".
// If no record has been stored, Int16 returns
// the zero value.
func (o *T) Int16() int16 {
	rec := get(o.db, keyInt16)
	if rec == nil {
		return 0
	}
	return int16(binary.BigEndian.Uint16(rec))
}

// PutInt16 stores v as a record under the key "Int16".
func (o *T) PutInt16(v int16) {
	rec := make([]byte, 2)
	binary.BigEndian.PutUint16(rec, uint16(v))
	put(o.db, keyInt16, rec)
}

// Int32 reads the record stored under key "Int32".
// If no record has been stored, Int32 returns
// the zero value.
func (o *T) Int32() int32 {
	rec := get(o.db, keyInt32)
	if rec == nil {
		return 0
	}
	return int32(binary.BigEndian.Uint32(rec))
}

// PutInt32 stores v as a record under the key "Int32".
func (o *T) PutInt32(v int32) {
	rec := make([]byte, 4)
	binary.BigEndian.PutUint32(rec, uint32(v))
	put(o.db, keyInt32, rec)
}

// Int64 reads the record stored under key "Int64".
// If no record has been stored, Int64 returns
// the zero value.
func (o *T) Int64() int64 {
	rec := get(o.db, keyInt64)
	if rec == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(rec))
}

// PutInt64 stores v as a record under the key "Int64".
func (o *T) PutInt64(v int64) {
	rec := make([]byte, 8)
	binary.BigEndian.PutUint64(rec, uint64(v))
	put(o.db, keyInt64, rec)
}

var (
	keyBool   = []byte("Bool")
	keyByte   = []byte("Byte")
	keyInt16  = []byte("Int16")
	keyInt32  = []byte("Int32")
	keyInt64  = []byte("Int64")
	keyInt8   = []byte("Int8")
	keyUint16 = []byte("Uint16")
	keyUint32 = []byte("Uint32")
	keyUint64 = []byte("Uint64")
)

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
