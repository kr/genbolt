package main

import (
	"io/ioutil"
	"os"
)

func main() {
	in, out := os.Args[1], os.Args[2]
	b, err := gen(in)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(out, b, 0644)
}
