package vf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Array returns a frozen dgo.Array that represents a copy of the given value. The value can be
// a slice or an Iterable
func Array(value interface{}) dgo.Array {
	return internal.Array(value)
}

// ArrayWithCapacity creates a new mutable array of the given type and initial capacity. The type can be nil, the
// zero value of a go slice, a dgo.ArrayType, or a dgo string that parses to a dgo.ArrayType.
func ArrayWithCapacity(typ interface{}, capacity int) dgo.Array {
	return internal.ArrayWithCapacity(capacity, typ)
}

// WrapSlice wraps the given slice in an array. Unset entries in the slice will be replaced by Nil.
func WrapSlice(slice []dgo.Value) dgo.Array {
	return internal.WrapSlice(slice)
}

// Values returns a frozen dgo.Array that represents the given values. All values
// are guaranteed to be frozen.
func Values(values ...interface{}) dgo.Array {
	return internal.Values(values)
}

// MutableValues returns a dgo.Array that represents the given values
func MutableValues(values ...interface{}) dgo.Array {
	return internal.MutableValues(values)
}

// Strings returns a frozen dgo.Array that represents the given strings
func Strings(values ...string) dgo.Array {
	return internal.Strings(values)
}

// Integers returns a frozen dgo.Array that represents the given ints
func Integers(values ...int) dgo.Array {
	return internal.Integers(values)
}

// Arguments returns an immutable Arguments instance that represents the given values
func Arguments(values ...interface{}) dgo.Arguments {
	return internal.Arguments(values)
}

// ArgumentsFromArray returns an Arguments instance backed by the given array
func ArgumentsFromArray(values dgo.Array) dgo.Arguments {
	return internal.ArgumentsFromArray(values)
}
