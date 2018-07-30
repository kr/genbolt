package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

// Hello, this is the root.
type Root struct {
	db *bolt.Tx
}

// NewRoot returns a new Root for tx.
//
// Hello, this is the root.
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

// F, what a lovely field, F.
func (o *Root) F() *T {
	return &T{bucket(o.db, keyF)}
}

// RootFoo is a root with a longer name.
type RootFoo struct {
	db *bolt.Tx
}

// NewRootFoo returns a new RootFoo for tx.
//
// RootFoo is a root with a longer name.
func NewRootFoo(tx *bolt.Tx) *RootFoo {
	return &RootFoo{tx}
}

func ViewFoo(db *bolt.DB, f func(*RootFoo, *bolt.Tx) error) error {
	return db.View(func(tx *bolt.Tx) error {
		return f(&RootFoo{tx}, tx)
	})
}

func UpdateFoo(db *bolt.DB, f func(*RootFoo, *bolt.Tx) error) error {
	return db.Update(func(tx *bolt.Tx) error {
		return f(&RootFoo{tx}, tx)
	})
}

func (o *RootFoo) F() *T {
	return &T{bucket(o.db, keyF)}
}

// Rootbar isn't a root!
type Rootbar struct {
	db *bolt.Bucket
}

func (o *Rootbar) F() *T {
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
