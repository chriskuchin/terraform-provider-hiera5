package tf

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Named returns the type with the given name from the global type registry. The function returns
// nil if no type has been registered under the given name.
func Named(name string) dgo.NamedType {
	return internal.NamedType(name)
}

// Parameterized returns the named type amended with the given parameters.
func Parameterized(named dgo.NamedType, params dgo.Array) dgo.NamedType {
	return internal.ParameterizedType(named, params)
}

// ExactNamed returns the exact NamedType that represents the given value.
//
// This is the function that the value.Type() method of a named type instance uses to
// obtain the actual type.
func ExactNamed(typ dgo.NamedType, value dgo.Value) dgo.NamedType {
	return internal.ExactNamedType(typ, value)
}

// NewNamed registers a new named and optionally parameterized type under the given name with the global type registry.
// The method panics if a type has already been registered with the same name.
//
// name: name of the type
//
// ctor: optional constructor that creates new values of this type
//
// extractor: optional extractor of the value used when serializing/deserializing this type
//
// implType optional reflected zero value type of implementation
//
// ifdType optional reflected nil value of interface type
//
// asgChecker optional function to check what other types that are assignable to this type
//
// params optional parameter Array
func NewNamed(
	name string,
	ctor dgo.Constructor,
	extractor dgo.InitArgExtractor,
	implType, ifdType reflect.Type,
	asgChecker dgo.AssignableChecker) dgo.NamedType {
	return internal.NewNamedType(name, ctor, extractor, implType, ifdType, asgChecker)
}

// NamedFromReflected returns the named type for the reflected implementation type from the global type
// registry. The function returns nil if no such type has been registered.
func NamedFromReflected(rt reflect.Type) dgo.NamedType {
	return internal.NamedTypeFromReflected(rt)
}

// RemoveNamed removes a named type from the global type registry. It is primarily intended for
// testing purposes.
func RemoveNamed(name string) {
	internal.RemoveNamedType(name)
}
