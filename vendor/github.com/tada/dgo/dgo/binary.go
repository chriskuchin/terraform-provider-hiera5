package dgo

type (
	// BinaryType is the type that represents a Binary value
	BinaryType interface {
		sizedType
		Factory

		IsInstance([]byte) bool
	}

	// Binary represents a sequence of bytes. Its string representation is a strictly base64 encoded string
	Binary interface {
		Value
		Comparable
		Mutability
		ReflectedValue

		// Copy returns a copy of the Binary. The copy is frozen or mutable depending on
		// the given argument. A request to create a frozen copy of an already frozen Binary
		// is a no-op that returns the receiver.
		Copy(frozen bool) Binary

		// Encode returns the strict base64 encoding of the given bytes
		Encode() string

		// GoBytes returns a copy of the internal array to ensure immutability
		GoBytes() []byte
	}
)
