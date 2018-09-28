package db

import (
	"testing"

	bolt "github.com/coreos/bbolt"
)

func TestRoot(t *testing.T) {
	db, err := bolt.Open("db", 0600, nil)
	must(t, err)
	var got int32 = 1
	must(t, View(db, func(root *Root) error {
		got = root.A().B().N()
		return nil
	}))
	var want int32 = 0
	if got != want {
		t.Errorf("root.A().B().N() = %d, want %d", got, want)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
