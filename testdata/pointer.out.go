package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) U() *U {
	return &U{bucket(o.db, keyU)}
}

type U struct {
	db *bolt.Bucket
}

var (
	keyU = []byte("U")
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
