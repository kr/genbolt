// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/jsonmap.in.go.

package db

import binary "encoding/binary"
import json "encoding/json"
import bolt "github.com/coreos/bbolt"
import sample "github.com/kr/genbolt/testdata/sample"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) Bucket() *bolt.Bucket {
	return o.db
}

func (o *T) J() *SampleJSONMap {
	return &SampleJSONMap{bucket(o.db, keyJ)}
}

type SampleJSONMap struct {
	db *bolt.Bucket
}

func (o *SampleJSONMap) Bucket() *bolt.Bucket {
	return o.db
}

func (o *SampleJSONMap) Get(key []byte) *sample.JSON {
	rec := o.db.Get(key)
	if rec == nil {
		return nil
	}
	v := new(sample.JSON)
	err := json.Unmarshal(rec, json.Unmarshaler(v))
	if err != nil {
		panic(err)
	}
	return v
}

func (o *SampleJSONMap) GetByString(key string) *sample.JSON {
	return o.Get([]byte(key))
}

func (o *SampleJSONMap) Put(key []byte, v *sample.JSON) {
	rec, err := json.Marshal(json.Marshaler(v))
	if err != nil {
		panic(err)
	}
	put(o.db, key, rec)
}

func (o *SampleJSONMap) PutByString(key string, v *sample.JSON) {
	o.Put([]byte(key), v)
}

var (
	keyJ = []byte("J")
)

type db interface {
	Writable() bool
	CreateBucketIfNotExists([]byte) (*bolt.Bucket, error)
	Bucket([]byte) *bolt.Bucket
}

func bucket(db db, key []byte) *bolt.Bucket {
	if !db.Writable() {
		return db.Bucket(key)
	}
	b, err := db.CreateBucketIfNotExists(key)
	if err != nil {
		panic(err)
	}
	return b
}

func put(b *bolt.Bucket, key, value []byte) {
	err := b.Put(key, value)
	if err != nil {
		panic(err)
	}
}
