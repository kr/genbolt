package db

import (
	myjson "encoding/json"

	"github.com/kr/genbolt/testdata/sample"
)

var _ myjson.Marshaler = (*sample.JSON)(nil)

type T struct {
	J *sample.JSON
}
