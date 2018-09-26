// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/map.in.go.

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

type V struct {
	db *bolt.Bucket
}

func (o *V) Bucket() *bolt.Bucket {
	return o.db
}

func (o *T) M() *MapOfV {
	return &MapOfV{bucket(o.db, keyM)}
}

type MapOfV struct {
	db *bolt.Bucket
}

func (o *MapOfV) Bucket() *bolt.Bucket {
	return o.db
}

func (o *MapOfV) Get(key []byte) *V {
	return &V{bucket(o.db, key)}
}

func (o *MapOfV) GetByString(key string) *V {
	return &V{bucket(o.db, []byte(key))}
}

var (
	keyM = []byte("M")
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
