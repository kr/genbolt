// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/jsonseq.in.go.

package db

import bytes "bytes"
import binary "encoding/binary"
import json "encoding/json"
import bolt "github.com/coreos/bbolt"
import sample "github.com/kr/genbolt/testdata/sample"

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

// J gets the child bucket with key "J" from o.
//
// J creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *SeqOfSampleJSON;
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *T) J() *SeqOfSampleJSON {
	return &SeqOfSampleJSON{bucket(o.db, keyJ)}
}

// SeqOfSampleJSON is a bucket with sequential numeric keys,
// holding records of type *sample.JSON.
type SeqOfSampleJSON struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *SeqOfSampleJSON) Bucket() *bolt.Bucket {
	return o.db
}

// Get reads the record stored in o under sequence number n.
//
// If no record has been stored, it returns
// a pointer to
// the zero value.
func (o *SeqOfSampleJSON) Get(n uint64) *sample.JSON {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	rec := get(o.db, key)
	v := new(sample.JSON)
	if rec == nil {
		return v
	}
	err := json.Unmarshal(rec, json.Unmarshaler(v))
	if err != nil {
		panic(err)
	}
	return v
}

// Add stores v in o under a new sequence number.
// It writes the new sequence number to *np
// before marshaling v. It is okay for
// np to point to a field inside v, to store
// the sequence number in the new record.
func (o *SeqOfSampleJSON) Add(v *sample.JSON, np *uint64) {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	*np = n
	o.Put(n, v)
}

// Put stores v in o as a record under sequence number n.
func (o *SeqOfSampleJSON) Put(n uint64, v *sample.JSON) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	rec, err := json.Marshal(json.Marshaler(v))
	if err != nil {
		panic(err)
	}
	put(o.db, key, rec)
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
