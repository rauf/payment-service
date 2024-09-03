package serde

import "io"

// Serde is an interface that defines the methods that a serializer/deserializer should implement.
type Serde interface {
	Serializer
	Deserializer
}

// Serializer is an interface that defines the method that a serializer should implement.
type Serializer interface {
	Serialize(w io.Writer, data any) error
}

// Deserializer is an interface that defines the method that a deserializer should implement.
type Deserializer interface {
	Deserialize(r io.Reader, v any) error
}
