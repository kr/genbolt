package db

import (
	"encoding/json"

	"github.com/kr/genbolt/testdata/sample"
)

var _ json.Marshaler = (*sample.JSON)(nil)

type T struct {
	J *sample.JSON
}
