// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/seqmap.in.go.

package db

import bytes "bytes"
import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize
const _ = bytes.MinRead

// T is a bucket with a static set of elements.
// Accessor methods read and write records
// and open child buckets.
type T struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *T) Bucket() *bolt.Bucket {
	return o.db
}

// S gets the child bucket with key "S" from o.
//
// S creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *SeqOfMapOfBool;
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *T) S() *SeqOfMapOfBool {
	return &SeqOfMapOfBool{bucket(o.db, keyS)}
}

// SeqOfMapOfBool is a bucket with sequential numeric keys,
// holding child buckets of type MapOfBool.
type SeqOfMapOfBool struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *SeqOfMapOfBool) Bucket() *bolt.Bucket {
	return o.db
}

// Get gets child bucket n from o.
//
// It creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *MapOfBool;
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *SeqOfMapOfBool) Get(n uint64) *MapOfBool {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	return &MapOfBool{bucket(o.db, key)}
}

// Add creates and returns a new, empty child bucket to o
// with a new sequence number.
//
// It panics if called in a read-only transaction.
func (o *SeqOfMapOfBool) Add() (*MapOfBool, uint64) {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	return o.Get(n), n
}

// MapOfBool is a bucket with arbitrary keys,
// holding records of type bool.
type MapOfBool struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *MapOfBool) Bucket() *bolt.Bucket {
	return o.db
}

// Get reads the record stored in o under the given key.
//
// If no record has been stored, it returns
// the zero value.
func (o *MapOfBool) Get(key []byte) bool {
	rec := get(o.db, key)
	return len(rec) > 0 && rec[0] != 0
}

// GetByString is equivalent to o.Get([]byte(key)).
func (o *MapOfBool) GetByString(key string) bool {
	return o.Get([]byte(key))
}

// Put stores v in o as a record under the given key.
func (o *MapOfBool) Put(key []byte, v bool) {
	rec := make([]byte, 1)
	if v {
		rec[0] = 1
	}
	put(o.db, key, rec)
}

// PutByString is equivalent to o.Put([]byte(key), v).
func (o *MapOfBool) PutByString(key string, v bool) {
	o.Put([]byte(key), v)
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
	if bu, ok := db.(*bolt.Bucket); ok && bu == nil {
		return nil // can happen in read-only txs
	}
	if !db.Writable() {
		return db.Bucket(key)
	}
	b, err := db.CreateBucketIfNotExists(key)
	if err != nil {
		panic(err)
	}
	return b
}

func get(b *bolt.Bucket, key []byte) []byte {
	if b == nil {
		return nil
	}
	return b.Get(key)
}

func put(b *bolt.Bucket, key, value []byte) {
	err := b.Put(key, value)
	if err != nil {
		panic(err)
	}
}
