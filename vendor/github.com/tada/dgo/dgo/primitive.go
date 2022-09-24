package dgo

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"time"
)

type (
	// Hash is the type returned by the Value.HashCode() method.
	Hash = int32

	// A Value represents an immutable value of some Type
	//
	// Values should not be compared using == (depending on value type, it may result in a panic: runtime error:
	//  comparing uncomparable). Instead, the method Equal or the function Identical should be used.
	//
	// There is no Value representation of "null" other than comparing to nil. The "undefined" found in some languages
	// such as TypeScript doesn't exist either. Instead, places where it should be relevant, such as when examining
	// if a Map contains a Value (or nil) for a certain key, methods will return a Value together with a bool that
	// indicates if a mapping was present or not.
	Value interface {
		fmt.Stringer

		// Type returns the type of this value
		Type() Type

		// Equals returns true if this value if equal to the given value. For complex objects, this
		// comparison is deep.
		Equals(other interface{}) bool

		// HashCode returns the computed hash code of the value
		HashCode() Hash
	}

	// ReflectedValue is implemented by all values that can be assigned to a reflected value in
	// a go native form.
	ReflectedValue interface {
		// ReflectTo assigns the go native form of this value to the given reflect.Value. The given
		// reflect.Value must be of a type that the native form of this value is assignable to or a pointer
		// to such a type.
		ReflectTo(value reflect.Value)
	}

	// Nil is implemented by the singleton that represents nil
	Nil interface {
		Value

		GoNil() interface{}
	}

	// Number is implemented by Float and Integer implementations
	Number interface {
		// Integer returns this number as an Integer
		Integer() Integer

		// Float returns this number as a Float
		Float() Float

		// ToInt returns this number as an int64
		ToInt() (int64, bool)

		// ToFloat returns this number as an float64
		ToFloat() (float64, bool)

		// ToBigInt returns this number as an *big.Int
		ToBigInt() *big.Int

		// ToBigFloat returns this number as an *big.Float
		ToBigFloat() *big.Float
	}

	// Integer value is an int64 that implements the Value interface
	Integer interface {
		Value
		Number
		Comparable
		ReflectedValue

		// Dec returns the integer -1
		Dec() Integer

		// Inc returns the integer +1
		Inc() Integer

		// GoInt returns the Go native representation of this value
		GoInt() int64
	}

	// BigInt value is a *big.Int that implements the Value interface
	BigInt interface {
		Integer
		GoBigInt() *big.Int
	}

	// Float value is a float64 that implements the Value interface
	Float interface {
		Value
		Number
		Comparable
		ReflectedValue

		// GoFloat returns the Go native representation of this value
		GoFloat() float64
	}

	// BigFloat value is a float64 that implements the Value interface
	BigFloat interface {
		Float
		GoBigFloat() *big.Float
	}

	// String value is a string that implements the Value interface
	String interface {
		Value
		Comparable
		ReflectedValue

		// GoString returns the Go native representation of this value
		GoString() string
	}

	// Regexp value is a *regexp.Regexp that implements the Value interface
	Regexp interface {
		Value
		ReflectedValue

		// GoRegexp returns the Go native representation of this value
		GoRegexp() *regexp.Regexp
	}

	// Time value is a *time.Time that implements the Value interface
	Time interface {
		Value
		ReflectedValue

		// GoTime returns the Go native representation of this value
		GoTime() *time.Time

		// SecondsWithFraction returns the number of seconds since since January 1, 1970 UTC. The fraction
		// will have nano-second precision for values in the years 1679 to 2261 and micro second precision
		// for years that a time can represent outside of that range.
		SecondsWithFraction() float64
	}

	// Boolean value
	Boolean interface {
		Value
		Comparable
		ReflectedValue

		// GoBool returns the Go native representation of this value
		GoBool() bool
	}
	// Native is a wrapper of a runtime value such as a chan or a pointer for which there is no proper immutable Value
	// representation
	Native interface {
		Value
		ReflectedValue

		// GoValue returns the Go native representation of this value
		GoValue() interface{}

		// ReflectValue returns a pointer to the reflect.Value representation of this value
		ReflectValue() *reflect.Value
	}

	// Comparable imposes natural ordering on its implementations. A Comparable is only comparable to other
	// values of its own type with the exception of Nil which is less than everything else and the special
	// case when Integer is compared to Float. Such a comparison will convert the Integer to a Float.
	Comparable interface {
		// CompareTo compares this value with the given value for order. Returns a negative integer, zero, or a positive
		// integer as this value is less than, equal to, or greater than the specified object and a bool that indicates
		// if the comparison was at all possible.
		CompareTo(other interface{}) (int, bool)
	}

	// RecursionGuard guards against endless recursion when checking if one deep type is assignable to another.
	// A Hit is detected once both a and b has been added more than once (can be on separate calls to Append). The
	// RecursionGuard is in itself immutable.
	RecursionGuard interface {
		// Append creates a new RecursionGuard guaranteed to contain both a and b. The new instance is returned.
		Append(a, b Value) RecursionGuard

		// Hit returns true if both a and b has been appended more than once.
		Hit() bool

		// Swap returns the guard with its two internal guards for a and b swapped.
		Swap() RecursionGuard
	}
)
