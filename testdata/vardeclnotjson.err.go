package db

import (
	"fmt"

	"github.com/kr/genbolt/testdata/sample"
)

var _ fmt.Stringer = (*sample.Stringer)(nil)
