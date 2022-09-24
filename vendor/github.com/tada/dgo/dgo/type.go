package dgo

import (
	"math"
	"reflect"
	"regexp"
	"time"
)

// UnboundedSize is returned by the Max() function of types that have a size constraint but no upper bound
const UnboundedSize = math.MaxUint32

type (
	// A Type describes an immutable Value. The Type is in itself also a Value
	Type interface {
		Value

		// Assignable returns true if a variable or parameter of this type can be hold a value of the other type
		Assignable(other Type) bool

		// Instance returns true if the value is an instance of this type
		Instance(value interface{}) bool

		// TypeIdentifier returns a unique identifier for this type. The TypeIdentifier is intended to be used by
		// decorators providing string representation of the type
		TypeIdentifier() TypeIdentifier

		// ReflectType returns the reflect.Type that corresponds to the receiver
		ReflectType() reflect.Type
	}

	// Meta is the description of a Type.
	Meta interface {
		Type

		// Describes returns the type that the meta type describes.
		Describes() Type
	}

	// IntegerType describes integers that are within an inclusive or exclusive range
	IntegerType interface {
		Type

		// Inclusive returns true if this range has an inclusive end
		Inclusive() bool

		// Max returns the maximum constraint
		Max() Integer

		// Min returns the minimum constraint
		Min() Integer
	}

	// FloatType describes floating point numbers that are within an inclusive or exclusive range
	FloatType interface {
		Type

		// Inclusive returns true if this range has an inclusive end
		Inclusive() bool

		// Max returns the maximum constraint
		Max() Float

		// Min returns the minimum constraint
		Min() Float
	}

	// BooleanType matches the true and false literals
	BooleanType interface {
		Type

		// IsInstance returns true if the Go native value is represented by this type
		IsInstance(value bool) bool
	}

	// RegexpType matches regular expressions
	RegexpType interface {
		Type

		// IsInstance returns true if the Go native value is represented by this type
		IsInstance(regexp *regexp.Regexp) bool
	}

	// TimeType matches time values
	TimeType interface {
		Type

		// IsInstance returns true if the Go native value is represented by this type
		IsInstance(tm time.Time) bool
	}

	// sizedType is implemented by types that may have a size constraint
	// such as String, Array, or Map
	sizedType interface {
		Type

		// Max returns the maximum size for instances of this type or
		// the constant UnboundedSize if here is no max limit
		Max() int

		// Min returns the minimum size for instances of this type
		Min() int

		// Unbounded returns true when the type has no size constraint
		Unbounded() bool
	}

	// StringType is a sizedType.
	StringType interface {
		sizedType
	}

	// PatternType is represents all strings that matches a regular expression pattern.
	PatternType interface {
		StringType

		GoRegexp() *regexp.Regexp
	}

	// NativeType is the type for all Native values
	NativeType interface {
		Type

		// GoType returns the reflect.Type
		GoType() reflect.Type
	}

	// ErrorType is the type for all error values
	ErrorType interface {
		Type

		// IsInstance returns true if the Go native value is represented by this type
		IsInstance(error) bool
	}

	// AliasContainer is implemented by types and values that can contain other types.
	//
	// The parser uses this interface to perform in-place replacement of aliases
	AliasContainer interface {
		Resolve(AliasAdder)
	}

	// Alias is a named reference of another type which can be resolved using an AliasAdder
	Alias interface {
		Type

		// Reference returns the name of the aliased type.
		Reference() String
	}

	// An AliasAdder maintains mappings of names to types
	AliasAdder interface {
		// Add adds the type t with the given name to this map
		Add(t Type, name String)

		// GetType returns the type with the given name or nil if the type isn't found
		GetType(n String) Type

		// Replace replaces aliases with their concrete value.
		Replace(Value) Value
	}

	// An AliasMap maps names to types and vice versa.
	AliasMap interface {
		// Collect is used to collect new aliases in a manner that is safe from a concurrency standpoint. If
		// new aliases were added by the given function argument, it returns a new fully resolved AliasAdder,
		// otherwise it returns itself.
		Collect(func(AliasAdder)) AliasMap

		// GetName returns the name for the given type or nil if the type isn't found
		GetName(t Type) String

		// GetType returns the type with the given name or nil if the type isn't found
		GetType(n String) Type
	}

	// GenericType is implemented by types that represent themselves stripped from
	// range and size constraints.
	GenericType interface {
		// Generic returns the generic type that this type represents stripped
		// from range and size constraints
		Generic() Type
	}

	// ExactType is implemented by types that match exactly one value but
	// isn't the same instance as the value.
	ExactType interface {
		Type

		// ExactValue returns the value that this type represents
		ExactValue() Value
	}

	// Factory provides the New method that types use to create new instances
	Factory interface {
		// New creates instances of this type based on arguments.
		//
		// The arguments can either be passed as one single value or as an instance
		// of Arguments (an special purpose Array).
		New(Value) Value
	}

	// DeepAssignable is implemented by values that need deep Assignable comparisons.
	DeepAssignable interface {
		DeepAssignable(guard RecursionGuard, other Type) bool
	}

	// DeepInstance is implemented by values that need deep Intance comparisons.
	DeepInstance interface {
		DeepInstance(guard RecursionGuard, value interface{}) bool
	}

	// ReverseAssignable indicates that the check for assignable must continue by delegating to the
	// type passed as an argument to the Assignable method. The reason is that types like AllOf, AnyOf
	// OneOf or types representing exact slices or maps, might need to check if individual types are
	// assignable.
	//
	// All implementations of Assignable must take into account the argument may implement this interface
	// do a reverse by calling the CheckAssignableTo function
	ReverseAssignable interface {
		// AssignableTo returns true if a variable or parameter of the other type can be hold a value of this type.
		// All implementations of Assignable must take into account that the given type might implement this method
		// do a reverse check before returning false.
		//
		// The guard is part of the internal endless recursion mechanism and should be passed as nil unless provided
		// by a DeepAssignable caller.
		AssignableTo(guard RecursionGuard, other Type) bool
	}
)
