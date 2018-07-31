package main

import (
	"bytes"
	"fmt"
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

func TestRun(t *testing.T) {
	for _, name := range glob(t, "testdata/*.use.go") {
		t.Run(name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "gentest")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			genName := strings.TrimSuffix(name, ".use.go") + ".out.go"
			copyGo(t, genName, dir, "db.go")
			copyGo(t, name, dir, "db_test.go")

			c := exec.Command("go", "test")
			c.Dir = dir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			err = c.Run()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func copyGo(t *testing.T, src, dir, dst string) {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	s := fmt.Sprintf("//line %s:1\n", src)
	b = append([]byte(s), b...)
	err = ioutil.WriteFile(filepath.Join(dir, dst), b, 0600)
	if err != nil {
		t.Fatal(err)
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
