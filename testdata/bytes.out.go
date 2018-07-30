package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

// T is a nice thing.
type T struct {
	db *bolt.Bucket
}

// B is a byte slice.
// It is useful.
func (o *T) B() []byte {
	v := o.db.Get(keyB)
	return v
}

// PutB stores x as the value of B.
//
// B is a byte slice.
// It is useful.
func (o *T) PutB(x []byte) {
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
