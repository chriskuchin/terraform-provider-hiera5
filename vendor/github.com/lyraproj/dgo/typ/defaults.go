// Package typ contains the static dgo types such as typ.String and typ.Any
package typ

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/loader"
)

// AllOf is the default AllOf type. Since it contains no types, everything is
// assignable to it.
var AllOf dgo.TernaryType = internal.DefaultAllOfType

// AnyOf is the default AnyOf type. Since it contains no types, nothing is
// assignable to it except all AnyOf types.
var AnyOf dgo.TernaryType = internal.DefaultAnyOfType

// OneOf is the default OneOf type. Since it contains no types, nothing is
// assignable to it except all OneOf types
var OneOf dgo.TernaryType = internal.DefaultOneOfType

// Array represents all array values
var Array dgo.ArrayType = internal.DefaultArrayType

// Tuple is represents all arrays since it's the tuple with one ellipsis argument of type any
var Tuple dgo.TupleType = internal.DefaultTupleType

// EmptyTuple is represents an empty arrays
var EmptyTuple dgo.TupleType = internal.EmptyTupleType

// Not represents a negated type. The default Not is negated Any so no other type
// is assignable to it.
var Not dgo.UnaryType = internal.DefaultNotType

// Map is the unconstrained type. It represents all Map values
var Map dgo.MapType = internal.DefaultMapType

// Any is a type that represents all values
var Any dgo.Type = internal.DefaultAnyType

// Nil is a type that represents the nil Value
var Nil dgo.Type = internal.DefaultNilType

// Boolean is a type that represents both true and false
var Boolean dgo.BooleanType = internal.DefaultBooleanType

// False is a type that only represents the value false
var False dgo.BooleanType = internal.FalseType

// True is a type that only represents the value true
var True dgo.BooleanType = internal.TrueType

// Number is a type that represents all numbers
var Number = internal.AnyOfType([]interface{}{internal.DefaultIntegerType, internal.DefaultFloatType})

// Float is a type that represents all floats
var Float dgo.FloatType = internal.DefaultFloatType

// Function is the type that represents all functions
var Function dgo.FunctionType = internal.DefaultFunctionType

// Integer is a type that represents all integers
var Integer dgo.IntegerType = internal.DefaultIntegerType

// Regexp is a type that represents all regexps
var Regexp dgo.RegexpType = internal.DefaultRegexpType

// Time is a type that represents all timestamps
var Time dgo.Type = internal.DefaultTimeType

// Binary is a type that represents all Binary values
var Binary dgo.BinaryType = internal.DefaultBinaryType

// String is a type that represents all strings
var String dgo.StringType = internal.DefaultStringType

// DgoString is a type that represents all strings with Dgo syntax
var DgoString dgo.StringType = internal.DefaultDgoStringType

// Error is a type that represents all implementation of error
var Error dgo.ErrorType = internal.DefaultErrorType

// Native is a type that represents all Native values
var Native dgo.Type = internal.DefaultNativeType

// Sensitive is a type that represents Sensitive values
var Sensitive dgo.UnaryType = internal.DefaultSensitiveType

// Type is a type that represents all types
var Type dgo.Type = internal.DefaultMetaType

// Loader represents all Loaders
var Loader dgo.Type = loader.Type

// DefiningLoader represents all Loaders capable of defining new entries
var DefiningLoader dgo.Type = loader.MutableType

// ParentedLoader represents all Loaders capable of defining new entries
var ParentedLoader dgo.Type = loader.ChildType
