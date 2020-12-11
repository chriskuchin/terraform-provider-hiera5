package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/lyraproj/dgo/dgo"
)

var reflectValueType = reflect.TypeOf((*dgo.Value)(nil)).Elem()

// New creates an instance of the given type from the given argument
func New(typ dgo.Type, argument dgo.Value) dgo.Value {
	if nt, ok := typ.(dgo.Factory); ok {
		return nt.New(argument)
	}
	if args, ok := argument.(dgo.Arguments); ok {
		if args.Len() != 1 {
			panic(fmt.Errorf(`unable to create a %s from arguments %s`, typ, args))
		}
		argument = args.Get(0)
	}
	if typ.Instance(argument) {
		return argument
	}
	panic(fmt.Errorf(`unable to create a %s from %s`, typ, argument))
}

// Value returns the dgo.Value representation of its argument. If the argument type
// is known, it will be more efficient to use explicit methods such as Float(), String(),
// Map(), etc.
func Value(v interface{}) dgo.Value {
	// This goFunc is kept very small to enable inlining so this
	// if statement should not be baked in to the grand switch
	// in the value goFunc
	if gv, ok := v.(dgo.Value); ok {
		return gv
	}
	if dv := value(v); dv != nil {
		return dv
	}
	return ValueFromReflected(reflect.ValueOf(v))
}

func value(v interface{}) dgo.Value {
	var dv dgo.Value
	switch v := v.(type) {
	case nil:
		dv = Nil
	case bool:
		dv = boolean(v)
	case string:
		dv = String(v)
	case []byte:
		dv = Binary(v, true)
	case []string:
		dv = Strings(v)
	case []int:
		dv = Integers(v)
	case *regexp.Regexp:
		dv = Regexp(v)
	case time.Time:
		dv = (*timeVal)(&v)
	case error:
		dv = &errw{v}
	case json.Number:
		dv = valueFromJSONNumber(v)
	case reflect.Value:
		dv = ValueFromReflected(v)
	default:
		if i, ok := ToInt(v); ok {
			dv = intVal(i)
		} else {
			var f float64
			if f, ok = ToFloat(v); ok {
				dv = floatVal(f)
			}
		}
	}
	return dv
}

func valueFromJSONNumber(v json.Number) dgo.Value {
	if i, err := v.Int64(); err == nil {
		return Integer(i)
	}
	f, err := v.Float64()
	if err != nil {
		panic(err)
	}
	return Float(f)
}

// ValueFromReflected converts the given reflected value into an immutable dgo.Value
func ValueFromReflected(vr reflect.Value) dgo.Value {
	// Invalid shouldn't happen, but needs a check
	if !vr.IsValid() {
		return Nil
	}

	isPtr := false
	switch vr.Kind() {
	case reflect.Slice:
		return ArrayFromReflected(vr, true)
	case reflect.Map:
		return FromReflectedMap(vr, true)
	case reflect.Interface:
		if vr.Type().NumMethod() == 0 {
			return ValueFromReflected(vr.Elem())
		}
	case reflect.Ptr:
		if vr.IsNil() {
			return Nil
		}
		isPtr = true
	case reflect.Func:
		return (*goFunc)(&vr)
	}

	if vr.CanInterface() {
		vi := vr.Interface()
		if v, ok := vi.(dgo.Value); ok {
			return v
		}
		if v := value(vi); v != nil {
			return v
		}
	}

	if isPtr {
		er := vr.Elem()
		// Pointer to struct should have been handled at this point or it is a pointer to
		// an unknown struct and should be a native
		if er.Kind() != reflect.Struct {
			return ValueFromReflected(er)
		}
	}
	// Value as unsafe. Immutability is not guaranteed
	return Native(vr)
}

// FromValue converts a dgo.Value into a go native value. The given `dest` must be a pointer
// to the expected native value.
func FromValue(src dgo.Value, dest interface{}) {
	dp := reflect.ValueOf(dest)
	if reflect.Ptr != dp.Kind() {
		panic("destination is not a pointer")
	}
	ReflectTo(src, dp.Elem())
}

// ReflectTo assigns the given dgo.Value to the given reflect.Value
func ReflectTo(src dgo.Value, dest reflect.Value) {
	if !dest.Type().AssignableTo(reflectValueType) {
		if rv, ok := src.(dgo.ReflectedValue); ok {
			rv.ReflectTo(dest)
			return
		}
	}
	dest.Set(reflect.ValueOf(src))
}

// Add well known types like regexp, time, etc. here
var wellKnownTypes = map[reflect.Type]dgo.Type{
	reflect.TypeOf(&regexp.Regexp{}): DefaultRegexpType,
	reflect.TypeOf(time.Time{}):      DefaultTimeType,
}
