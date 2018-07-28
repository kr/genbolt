package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) M() *VMap {
	return &VMap{bucket(o.db, keyM)}
}

type V struct {
	db *bolt.Bucket
}

type VMap struct {
	db *bolt.Bucket
}

func (o *VMap) Get(key []byte) *V {
	return &V{bucket(o.db, key)}
}

func (o *VMap) GetByString(key string) *V {
	return &V{bucket(o.db, []byte(key))}
}

var (
	keyM = []byte("M")
)

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
