// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/binseq.in.go.

package db

import encoding "encoding"
import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"
import sample "github.com/kr/genbolt/testdata/sample"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) Bucket() *bolt.Bucket {
	return o.db
}

func (o *T) B() *SampleBinarySeq {
	return &SampleBinarySeq{bucket(o.db, keyB)}
}

type SampleBinarySeq struct {
	db *bolt.Bucket
}

func (o *SampleBinarySeq) Bucket() *bolt.Bucket {
	return o.db
}

func (o *SampleBinarySeq) Get(n uint64) *sample.Binary {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
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

// Add adds v to the sequence.
// It writes the new sequence number to *np
// before marshaling v. Thus, it is okay for
// np to point to a field inside v, to store
// the sequence number in the new record.
func (o *SampleBinarySeq) Add(v *sample.Binary, np *uint64) {
	n, err := o.db.NextSequence()
	if err != nil {
		panic(err)
	}
	*np = n
	o.Put(n, v)
}

func (o *SampleBinarySeq) Put(n uint64, v *sample.Binary) {
	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, n)
	rec, err := encoding.BinaryMarshaler(v).MarshalBinary()
	if err != nil {
		panic(err)
	}
	put(o.db, key, rec)
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