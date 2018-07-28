package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type Root struct {
	db *bolt.Tx
}

func NewRoot(tx *bolt.Tx) *Root {
	return &Root{tx}
}

func View(db *bolt.DB, f func(*Root, *bolt.Tx) error) error {
	return db.View(func(tx *bolt.Tx) error {
		return f(&Root{tx}, tx)
	})
}

func Update(db *bolt.DB, f func(*Root, *bolt.Tx) error) error {
	return db.Update(func(tx *bolt.Tx) error {
		return f(&Root{tx}, tx)
	})
}

func (o *Root) F() *T {
	return &T{bucket(o.db, keyF)}
}

type T struct {
	db *bolt.Bucket
}

var (
	keyF = []byte("F")
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
