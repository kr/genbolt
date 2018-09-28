// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/mapstruct.in.go.

package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

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

// V is a bucket with a static set of elements.
// Accessor methods read and write records
// and open child buckets.
type V struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *V) Bucket() *bolt.Bucket {
	return o.db
}

// M gets the child bucket with key "M" from o.
//
// M creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *MapOfV;
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *T) M() *MapOfV {
	return &MapOfV{bucket(o.db, keyM)}
}

// MapOfV is a bucket with arbitrary keys,
// holding child buckets of type V.
type MapOfV struct {
	db *bolt.Bucket
}

// Bucket returns o's underlying *bolt.Bucket object.
// This can be useful to access low-level database functions
// or other features not exposed by this generated code.
//
// Note, if o's transaction is read-only and the underlying
// bucket has not previously been created in a writable
// transaction, Bucket returns nil.
func (o *MapOfV) Bucket() *bolt.Bucket {
	return o.db
}

// Get gets the child bucket with the given key from o.
//
// It creates a new bucket if none exists
// and o's transaction is writable.
// Regardless, it always returns a non-nil *V;
// if the bucket doesn't exist
// and o's transaction is read-only, the returned value
// represents an empty bucket.
func (o *MapOfV) Get(key []byte) *V {
	return &V{bucket(o.db, key)}
}

// GetByString is equivalent to o.Get([]byte(key)).
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