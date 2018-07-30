// Generated by github.com/kr/genbolt. DO NOT EDIT.

package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) S() *USeq {
	return &USeq{bucket(o.db, keyS)}
}

type U struct {
	db *bolt.Bucket
}

type USeq struct {
	db *bolt.Bucket
}

func (o *USeq) Get(n uint64) *U {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	return &U{bucket(o.db, key)}
}

func (o *USeq) Add() *U {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	return o.Get(n)
}

var (
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
