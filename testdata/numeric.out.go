// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/numeric.in.go

package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) Bucket() *bolt.Bucket {
	return o.db
}

func (o *T) Bool() bool {
	rec := o.db.Get(keyBool)
	return rec[0] != 0
}

// PutBool stores v as the value of Bool.
func (o *T) PutBool(v bool) {
	rec := make([]byte, 1)
	if v {
		rec[0] = 1
	}
	put(o.db, keyBool, rec)
}

func (o *T) Byte() byte {
	rec := o.db.Get(keyByte)
	return rec[0]
}

// PutByte stores v as the value of Byte.
func (o *T) PutByte(v byte) {
	rec := []byte{v}
	put(o.db, keyByte, rec)
}

func (o *T) Uint16() uint16 {
	rec := o.db.Get(keyUint16)
	return binary.BigEndian.Uint16(rec)
}

// PutUint16 stores v as the value of Uint16.
func (o *T) PutUint16(v uint16) {
	rec := make([]byte, 2)
	binary.BigEndian.PutUint16(rec, v)
	put(o.db, keyUint16, rec)
}

func (o *T) Uint32() uint32 {
	rec := o.db.Get(keyUint32)
	return binary.BigEndian.Uint32(rec)
}

// PutUint32 stores v as the value of Uint32.
func (o *T) PutUint32(v uint32) {
	rec := make([]byte, 4)
	binary.BigEndian.PutUint32(rec, v)
	put(o.db, keyUint32, rec)
}

func (o *T) Uint64() uint64 {
	rec := o.db.Get(keyUint64)
	return binary.BigEndian.Uint64(rec)
}

// PutUint64 stores v as the value of Uint64.
func (o *T) PutUint64(v uint64) {
	rec := make([]byte, 8)
	binary.BigEndian.PutUint64(rec, v)
	put(o.db, keyUint64, rec)
}

func (o *T) Int8() int8 {
	rec := o.db.Get(keyInt8)
	return int8(rec[0])
}

// PutInt8 stores v as the value of Int8.
func (o *T) PutInt8(v int8) {
	rec := []byte{byte(v)}
	put(o.db, keyInt8, rec)
}

func (o *T) Int16() int16 {
	rec := o.db.Get(keyInt16)
	return int16(binary.BigEndian.Uint16(rec))
}

// PutInt16 stores v as the value of Int16.
func (o *T) PutInt16(v int16) {
	rec := make([]byte, 2)
	binary.BigEndian.PutUint16(rec, uint16(v))
	put(o.db, keyInt16, rec)
}

func (o *T) Int32() int32 {
	rec := o.db.Get(keyInt32)
	return int32(binary.BigEndian.Uint32(rec))
}

// PutInt32 stores v as the value of Int32.
func (o *T) PutInt32(v int32) {
	rec := make([]byte, 4)
	binary.BigEndian.PutUint32(rec, uint32(v))
	put(o.db, keyInt32, rec)
}

func (o *T) Int64() int64 {
	rec := o.db.Get(keyInt64)
	return int64(binary.BigEndian.Uint64(rec))
}

// PutInt64 stores v as the value of Int64.
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

func put(b *bolt.Bucket, key, value []byte) {
	err := b.Put(key, value)
	if err != nil {
		panic(err)
	}
}
