package db

import (
	"reflect"
	"testing"

	bolt "github.com/coreos/bbolt"
)

func TestRoot(t *testing.T) {
	db, err := bolt.Open("db", 0600, nil)
	must(t, err)
	must(t, db.Update(func(tx *bolt.Tx) error {
		bu, err := tx.CreateBucket([]byte("x"))
		must(t, err)
		dbT := &T{db: bu}
		want := []int64{1, 2, 3}
		dbT.PutInt64s(want)
		got := dbT.Int64s()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("round trip got = %v, want %v", got, want)
		}
		return nil
	}))
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
