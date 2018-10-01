// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/binmap.in.go.

package db

import bytes "bytes"
import encoding "encoding"
import binary "encoding/binary"
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

// B gets the child bucket with key "B" from o.
//
// B creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *MapOfSampleBinary;
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *T) B() *MapOfSampleBinary {
	return &MapOfSampleBinary{bucket(o.db, keyB)}
}

// MapOfSampleBinary is a bucket with arbitrary keys,
// holding records of type *sample.Binary.
type MapOfSampleBinary struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *MapOfSampleBinary) Bucket() *bolt.Bucket {
	return o.db
}

// Get reads the record stored in o under the given key.
//
// If no record has been stored, it returns
// a pointer to
// the zero value.
func (o *MapOfSampleBinary) Get(key []byte) *sample.Binary {
	rec := get(o.db, key)
	v := new(sample.Binary)
	if rec == nil {
		return v
	}
	err := encoding.BinaryUnmarshaler(v).UnmarshalBinary(rec)
	if err != nil {
		panic(err)
	}
	return v
}

// GetByString is equivalent to o.Get([]byte(key)).
func (o *MapOfSampleBinary) GetByString(key string) *sample.Binary {
	return o.Get([]byte(key))
}

// Put stores v in o as a record under the given key.
func (o *MapOfSampleBinary) Put(key []byte, v *sample.Binary) {
	rec, err := encoding.BinaryMarshaler(v).MarshalBinary()
	if err != nil {
		panic(err)
	}
	put(o.db, key, rec)
}

// PutByString is equivalent to o.Put([]byte(key), v).
func (o *MapOfSampleBinary) PutByString(key string, v *sample.Binary) {
	o.Put([]byte(key), v)
}

var (
	keyB = []byte("B")
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
