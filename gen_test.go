package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGen(t *testing.T) {
	for _, name := range glob(t, "testdata/*.in.go") {
		stem := strings.TrimSuffix(name, ".in.go")
		want, err := ioutil.ReadFile(stem + ".out.go")
		if err != nil {
			t.Error(err)
			continue
		}

		got, err := gen(name)
		if err != nil {
			t.Errorf("gen(%q) err = %v", name, err)
			continue
		}
		d := diff(got, want)
		if len(d) > 0 {
			t.Errorf("gen(%q): %s", name, d)
		}
	}
}

func TestGenError(t *testing.T) {
	for _, name := range glob(t, "testdata/*.err.go") {
		got, err := gen(name)
		if err == nil {
			t.Errorf("gen(%q) = [output], want error\n%s", name, got)
		}
	}
}

// diff returns a description of
// the difference between got and want.
// If got and want have the same contents,
// diff returns nil.
func diff(got, want []byte) []byte {
	if bytes.Equal(got, want) {
		return nil
	}
	dir, err := ioutil.TempDir("", "gentest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	err = ioutil.WriteFile(filepath.Join(dir, "got "), got, 0400)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(dir, "want"), want, 0400)
	if err != nil {
		panic(err)
	}
	c := exec.Command("diff", "-u", "got ", "want")
	c.Dir = dir
	c.Stderr = os.Stderr
	out, _ := c.Output()
	cmd := []byte("diff -u got want\n")
	return append(cmd, out...)
}

func glob(t *testing.T, pattern string) []string {
	a, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatal(err)
	}
	return a
}
