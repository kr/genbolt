package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) Bool() bool {
	v := o.db.Get(keyBool)
	return v[0] != 0
}

// SetBool stores x as the value of Bool.
func (o *T) SetBool(x bool) {
	v := make([]byte, 1)
	if x {
		v[0] = 1
	}
	put(o.db, keyBool, v)
}

func (o *T) Byte() byte {
	v := o.db.Get(keyByte)
	return v[0]
}

// SetByte stores x as the value of Byte.
func (o *T) SetByte(x byte) {
	v := []byte{x}
	put(o.db, keyByte, v)
}

func (o *T) Uint16() uint16 {
	v := o.db.Get(keyUint16)
	return binary.BigEndian.Uint16(v)
}

// SetUint16 stores x as the value of Uint16.
func (o *T) SetUint16(x uint16) {
	v := make([]byte, 2)
	binary.BigEndian.PutUint16(x)
	put(o.db, keyUint16, v)
}

func (o *T) Uint32() uint32 {
	v := o.db.Get(keyUint32)
	return binary.BigEndian.Uint32(v)
}

// SetUint32 stores x as the value of Uint32.
func (o *T) SetUint32(x uint32) {
	v := make([]byte, 4)
	binary.BigEndian.PutUint32(x)
	put(o.db, keyUint32, v)
}

func (o *T) Uint64() uint64 {
	v := o.db.Get(keyUint64)
	return binary.BigEndian.Uint64(v)
}

// SetUint64 stores x as the value of Uint64.
func (o *T) SetUint64(x uint64) {
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(x)
	put(o.db, keyUint64, v)
}

func (o *T) Int8() int8 {
	v := o.db.Get(keyInt8)
	return int8(v[0])
}

// SetInt8 stores x as the value of Int8.
func (o *T) SetInt8(x int8) {
	v := []byte{byte(x)}
	put(o.db, keyInt8, v)
}

func (o *T) Int16() int16 {
	v := o.db.Get(keyInt16)
	return int16(binary.BigEndian.Uint16(v))
}

// SetInt16 stores x as the value of Int16.
func (o *T) SetInt16(x int16) {
	v := make([]byte, 2)
	binary.BigEndian.PutUint16(uint16(x))
	put(o.db, keyInt16, v)
}

func (o *T) Int32() int32 {
	v := o.db.Get(keyInt32)
	return int32(binary.BigEndian.Uint32(v))
}

// SetInt32 stores x as the value of Int32.
func (o *T) SetInt32(x int32) {
	v := make([]byte, 4)
	binary.BigEndian.PutUint32(uint32(x))
	put(o.db, keyInt32, v)
}

func (o *T) Int64() int64 {
	v := o.db.Get(keyInt64)
	return int64(binary.BigEndian.Uint64(v))
}

// SetInt64 stores x as the value of Int64.
func (o *T) SetInt64(x int64) {
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(uint64(x))
	put(o.db, keyInt64, v)
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
