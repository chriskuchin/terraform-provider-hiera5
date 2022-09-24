package dgo

type (
	// Sensitive wraps another value and prevents that its string representation is used in logs.
	Sensitive interface {
		Value
		Mutability

		// Unwrap returns the wrapped value
		Unwrap() Value
	}
)
