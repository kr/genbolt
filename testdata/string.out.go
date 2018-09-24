// Generated by github.com/kr/genbolt. DO NOT EDIT.
// See testdata/string.in.go.

package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) Bucket() *bolt.Bucket {
	return o.db
}

func (o *T) S() string {
	rec := get(o.db, keyS)
	return string(rec)
}

// PutS stores v as the value of S.
func (o *T) PutS(v string) {
	rec := []byte(v)
	put(o.db, keyS, rec)
}

var (
	keyS = []byte("S")
)

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
