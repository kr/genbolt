package db

type T struct {
	Int8s   []int8
	Int16s  []int16
	Int32s  []int32 // aka rune
	Int64s  []int64
	Uint8s  []uint8 // aka byte
	Uint16s []uint16
	Uint32s []uint32
	Uint64s []uint64
	Bools   []bool

	// Strings don't have a fixed size, so we don't
	// special-case them like we do for the other
	// basic types. So this is a SeqOfString, not a
	// SliceOfString.
	Strings []string
}
