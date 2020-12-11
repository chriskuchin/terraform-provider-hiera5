package dgo

type (
	// Doer is performs some task
	Doer func()

	// Consumer is performs some task with a value
	Consumer func(value Value)

	// Mapper maps a value to another value
	Mapper func(value Value) interface{}

	// Producer is a function that produces a value
	Producer func() Value

	// Function is implemented by data types that can be called such as named or
	// unnamed functions and methods.
	Function interface {
		Value

		// Call calls the function with the arguments listed in the given Array and returns a result.
		// If the function has no return value, the result will be the empty slice to indicate zero
		// returned values.
		Call(args Array) []Value
	}

	// GoFunction is a Function specialization that wraps a go function
	GoFunction interface {
		Function

		// GoFunc returns a value that can be cast to the actual go function that this GoFunction
		// represents.
		GoFunc() interface{}
	}

	// FunctionType describes a Function
	FunctionType interface {
		Type

		// In returns the tuple that describes the arguments
		In() TupleType

		// Out returns the tuple that describes the return values
		Out() TupleType

		// Variadic returns true if the last argument is variadic, in which case it is always an Array
		Variadic() bool
	}
)
