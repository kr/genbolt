package db

import (
	"encoding/json"

	"github.com/kr/genbolt/testdata/sample"
)

var a json.Marshaler = (*sample.JSON)(nil)
