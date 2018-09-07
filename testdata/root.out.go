// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/root.in.go

package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

// NewRoot returns a new Root for tx.
//
// Hello, this is the root.
func NewRoot(tx *bolt.Tx) *Root {
	return &Root{tx}
}

func View(db *bolt.DB, f func(*Root) error) error {
	return db.View(func(tx *bolt.Tx) error {
		return f(&Root{tx})
	})
}

func Update(db *bolt.DB, f func(*Root) error) error {
	return db.Update(func(tx *bolt.Tx) error {
		return f(&Root{tx})
	})
}

// Hello, this is the root.
type Root struct {
	db *bolt.Tx
}

func (o *Root) Tx() *bolt.Tx {
	return o.db
}

// F, what a lovely field, F.
func (o *Root) F() *T {
	return &T{bucket(o.db, keyF)}
}

func (o *Root) S() *TSeq {
	return &TSeq{bucket(o.db, keyS)}
}

// NewRootFoo returns a new RootFoo for tx.
//
// RootFoo is a root with a longer name.
func NewRootFoo(tx *bolt.Tx) *RootFoo {
	return &RootFoo{tx}
}

func ViewFoo(db *bolt.DB, f func(*RootFoo) error) error {
	return db.View(func(tx *bolt.Tx) error {
		return f(&RootFoo{tx})
	})
}

func UpdateFoo(db *bolt.DB, f func(*RootFoo) error) error {
	return db.Update(func(tx *bolt.Tx) error {
		return f(&RootFoo{tx})
	})
}

// RootFoo is a root with a longer name.
type RootFoo struct {
	db *bolt.Tx
}

func (o *RootFoo) Tx() *bolt.Tx {
	return o.db
}

func (o *RootFoo) F() *T {
	return &T{bucket(o.db, keyF)}
}

// Rootbar isn't a root!
type Rootbar struct {
	db *bolt.Bucket
}

func (o *Rootbar) Bucket() *bolt.Bucket {
	return o.db
}

func (o *Rootbar) F() *T {
	return &T{bucket(o.db, keyF)}
}

type T struct {
	db *bolt.Bucket
}

func (o *T) Bucket() *bolt.Bucket {
	return o.db
}

type TSeq struct {
	db *bolt.Bucket
}

func (o *TSeq) Get(n uint64) *T {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	return &T{bucket(o.db, key)}
}

func (o *TSeq) Add() (*T, uint64) {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	return o.Get(n), n
}

var (
	keyF = []byte("F")
	keyS = []byte("S")
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
