package db

import (
	"encoding/json"

	"github.com/kr/genbolt/testdata/sample"
)

var _, a json.Marshaler = (*sample.JSON)(nil), (*sample.JSON)(nil)
