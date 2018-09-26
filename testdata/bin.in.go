package db

import (
	"encoding"

	"github.com/kr/genbolt/testdata/sample"
)

var _ encoding.BinaryMarshaler = (*sample.Binary)(nil)

type T struct {
	B *sample.Binary
}
