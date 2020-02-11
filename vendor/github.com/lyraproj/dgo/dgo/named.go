package dgo

import "reflect"

type (
	// Constructor is a function that creates a Value base on an initializer Map
	Constructor func(Value) Value

	// InitArgExtractor is a function that can extract initializer arguments from an
	// instance.
	InitArgExtractor func(Value) Value

	// AssignableChecker is a function that can determine whether or not another type is assignable
	// depending on given parameters.
	AssignableChecker func(self NamedType, typ Type) bool

	// NamedTypeExtension defines the extensions that a NamedType brings to a Type.
	NamedTypeExtension interface {
		Factory

		// AssignableType returns a reflect.Type that is either an interface that instances
		// of this type must implement, or an actual implementation. The default AssignableChecker
		// uses this type.
		AssignableType() reflect.Type

		// ExtractInitArg extracts the initializer argument from an instance.
		ExtractInitArg(Value) Value

		// Name returns the name of this type
		Name() string

		// Parameters returns the parameters for the type, or nil if the type isn't parameterized.
		Parameters() Array

		// ValueString returns the given value as a string. The Value must be an instance of this type.
		ValueString(value Value) string
	}

	// NamedType is implemented by types that are named and made available using an AliasMap
	NamedType interface {
		Type
		NamedTypeExtension
	}
)
