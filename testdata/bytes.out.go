package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) B() []byte {
	v := o.db.Get(keyB)
	return v
}

func (o *T) SetB(x []byte) {
	v := x
	put(o.db, keyB, v)
}

var (
	keyB = []byte("B")
)

func put(b *bolt.Bucket, key, value []byte) {
	err := b.Put(key, value)
	if err != nil {
		panic(err)
	}
}
