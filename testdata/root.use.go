package db

import (
	"testing"

	bolt "github.com/coreos/bbolt"
)

func TestRoot(t *testing.T) {
	db, err := bolt.Open("db", 0600, nil)
	must(t, err)
	must(t, View(db, func(root *Root) error {
		root.F()
		return nil
	}))
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
