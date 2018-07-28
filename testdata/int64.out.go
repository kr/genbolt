package db

import binary "encoding/binary"
import bolt "github.com/coreos/bbolt"

const _ = binary.MaxVarintLen16
const _ = bolt.MaxKeySize

type T struct {
	db *bolt.Bucket
}

func (o *T) N() int64 {
	v := o.db.Get(keyN)
	return int64(binary.BigEndian.Uint64(v))
}

func (o *T) SetN(x int64) {
	v := make([]byte, 8)
	binary.BigEndian.PutUint64(uint64(x))
	put(o.db, keyN, v)
}

var (
	keyN = []byte("N")
)

func put(b *bolt.Bucket, key, value []byte) {
	err := b.Put(key, value)
	if err != nil {
		panic(err)
	}
}
