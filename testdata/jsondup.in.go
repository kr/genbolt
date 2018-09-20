package db

import (
	"encoding/json"

	"github.com/kr/genbolt/testdata/sample"
)

var (
	_ json.Marshaler = (*sample.JSON)(nil)
	_ json.Marshaler = (*sample.JSON2)(nil)
)

type T struct {
	J *sample.JSON
	H *sample.JSON2
}
