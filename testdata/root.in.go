package db

// Hello, this is the root.
type Root struct {
	// F, what a lovely field, F.
	F *T
	S []*T
}

// RootFoo is a root with a longer name.
type RootFoo struct {
	F *T
}

// Rootbar isn't a root!
type Rootbar struct {
	F *T
}

type T struct{}
