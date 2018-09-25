package db

import (
	"encoding/json"

	foo "github.com/kr/genbolt/testdata/sample"
)

var _ json.Marshaler = (*foo.JSON)(nil)

type T struct {
	J *foo.JSON
}
