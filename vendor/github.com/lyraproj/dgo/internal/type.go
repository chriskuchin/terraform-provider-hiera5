package internal

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/lyraproj/dgo/dgo"
)

// CheckAssignableTo checks if the given type t implements the ReverseAssignable interface and if
// so, returns the value of calling its AssignableTo method with the other type as an
// argument. Otherwise, this goFunc returns false
func CheckAssignableTo(guard dgo.RecursionGuard, t, other dgo.Type) bool {
	if lt, ok := t.(dgo.ReverseAssignable); ok {
		if guard != nil {
			guard = guard.Swap()
		}
		return lt.AssignableTo(guard, other)
	}
	return false
}

// TypeFromReflected returns the dgo.Type that represents the given reflected type
func TypeFromReflected(vt reflect.Type) dgo.Type {
	if pt, ok := wellKnownTypes[vt]; ok {
		return pt
	}

	kind := vt.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		return ArrayType([]interface{}{TypeFromReflected(vt.Elem()), 0, math.MaxInt64})
	case reflect.Map:
		return MapType([]interface{}{TypeFromReflected(vt.Key()), TypeFromReflected(vt.Elem()), 0, math.MaxInt64})
	case reflect.Ptr:
		return OneOfType([]interface{}{TypeFromReflected(vt.Elem()), DefaultNilType})
	case reflect.Func:
		return exactFunctionType{vt}
	case reflect.Interface:
		switch vt.Name() {
		case ``:
			return DefaultAnyType
		case `error`:
			return DefaultErrorType
		}
	default:
		if pt, ok := primitivePTypes[vt.Kind()]; ok {
			return pt
		}
	}
	return &nativeType{vt}
}

// AsType returns the value as a type. If the value already is a type, it is returned. Otherwise the
// exact type of the value is returned.
func AsType(value dgo.Value) dgo.Type {
	if tp, ok := value.(dgo.Type); ok {
		return tp
	}
	return value.Type()
}

// ExactValue returns the "exact value" that a value represents. If the given value is a dgo.ExactType, then the value
// that it represents is the exact value. For all other cases, the exact value is the value itself.
func ExactValue(value dgo.Value) dgo.Value {
	if et, ok := value.(dgo.ExactType); ok {
		value = et.ExactValue()
	}
	return value
}

// Generic returns the generic form of the given type. All non exact types are considered generic
// and will be returned directly. Exact types will loose information about what instance they represent
// and also range and size information. Nested types will return a generic version of the contained
// types as well.
func Generic(t dgo.Type) dgo.Type {
	if et, ok := t.(dgo.GenericType); ok {
		return et.Generic()
	}
	return t
}

func illegalArgument(name, expected interface{}, args []interface{}, argno int) error {
	if len(args) == 1 {
		return fmt.Errorf(`illegal argument for %s. Expected %s, got %s`, name, expected, Value(args[argno]))
	}
	return fmt.Errorf(
		`illegal argument %d for %s with %d arguments. Expected %s, got %s`,
		argno+1, name, len(args), expected, Value(args[argno]))
}

func illegalArgumentCount(name string, min, max, actual int) error {
	var exp string
	switch {
	case max == math.MaxInt64:
		exp = fmt.Sprintf(`at least %d`, min)
	case min == max:
		exp = strconv.Itoa(min)
	case max-min == 1:
		exp = fmt.Sprintf(`%d or %d`, min, max)
	default:
		exp = fmt.Sprintf(`%d to %d`, min, max)
	}
	if name != `` {
		name = ` for ` + name
	}
	return fmt.Errorf(`illegal number of arguments%s. Expected %s, got %d`, name, exp, actual)
}

var primitivePTypes = map[reflect.Kind]dgo.Type{
	reflect.String:  DefaultStringType,
	reflect.Int:     DefaultIntegerType,
	reflect.Int8:    IntegerType(math.MinInt8, math.MaxInt8, true),
	reflect.Int16:   IntegerType(math.MinInt16, math.MaxInt16, true),
	reflect.Int32:   IntegerType(math.MinInt32, math.MaxInt32, true),
	reflect.Int64:   DefaultIntegerType,
	reflect.Uint:    IntegerType(0, math.MaxInt64, true),
	reflect.Uint8:   IntegerType(0, math.MaxUint8, true),
	reflect.Uint16:  IntegerType(0, math.MaxUint16, true),
	reflect.Uint32:  IntegerType(0, math.MaxUint32, true),
	reflect.Uint64:  IntegerType(0, math.MaxInt64, true),
	reflect.Float32: FloatType(-math.MaxFloat32, math.MaxFloat32, true),
	reflect.Float64: DefaultFloatType,
	reflect.Bool:    DefaultBooleanType,
}

// Parse parses the given string into a dgo.Value
var Parse func(s string) dgo.Value

// ParseFile parses the given content into a dgo.Type. The filename is used in error messages.
//
// The alias map is optional. If given, the parser will recognize the type aliases provided in the map
// and also add any new aliases declared within the parsed content to that map.
var ParseFile func(am dgo.AliasAdder, fileName, content string) dgo.Value

// TypeString produces the string that represents the given type
var TypeString func(dgo.Type) string
