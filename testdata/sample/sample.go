package sample

// JSON satisfies the json.Marshaler and json.Unmarshaler interfaces.
type JSON struct{}

func (s *JSON) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}

func (s *JSON) UnmarshalJSON([]byte) error {
	return nil
}

// JSON2 satisfies the json.Marshaler and json.Unmarshaler interfaces.
type JSON2 struct{}

func (s *JSON2) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}

func (s *JSON2) UnmarshalJSON([]byte) error {
	return nil
}

type JSONPointer = *JSON

type Stringer struct{}

func (s *Stringer) String() string {
	return ""
}

// Binary satisfies the encoding.BinaryMarshaler and
// encoding.BinaryUnmarshaler interfaces.
type Binary struct{}

func (s *Binary) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func (s *Binary) UnmarshalBinary([]byte) error {
	return nil
}
