package db

type Root struct {
	A *A
}

type A struct {
	B *B
}

type B struct {
	N int32
}
