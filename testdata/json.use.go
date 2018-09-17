package db

import (
	"testing"

	bolt "github.com/coreos/bbolt"
)

func TestJSON(t *testing.T) {
	db, err := bolt.Open("db", 0600, nil)
	must(t, err)
	must(t, db.Update(func(tx *bolt.Tx) error {
		bu, err := tx.CreateBucket([]byte("x"))
		must(t, err)
		t := &T{db: bu}
		t.J()
		return nil
	}))
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
