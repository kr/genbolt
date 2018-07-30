package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) S() string {
	v := o.db.Get(keyS)
	return string(v)
}

// PutS stores x as the value of S.
func (o *T) PutS(x string) {
	v := []byte(x)
	put(o.db, keyS, v)
}

var (
	keyS = []byte("S")
)

func put(b *bolt.Bucket, key, value []byte) {
	err := b.Put(key, value)
	if err != nil {
		panic(err)
	}
}
