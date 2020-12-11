package internal

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// floatVal is a float64 that implements the dgo.Value interface
	floatVal float64

	defaultFloatType int

	exactFloatType struct {
		exactType
		value floatVal
	}

	floatType struct {
		min       float64
		max       float64
		inclusive bool
	}
)

// DefaultFloatType is the unconstrained floatVal type
const DefaultFloatType = defaultFloatType(0)

var reflectFloatType = reflect.TypeOf(float64(0))

// FloatType returns a dgo.FloatType that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func FloatType(min, max float64, inclusive bool) dgo.FloatType {
	if min == max {
		if !inclusive {
			panic(fmt.Errorf(`non inclusive range cannot have equal min and max`))
		}
		return floatVal(min).Type().(dgo.FloatType)
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	if min == -math.MaxFloat64 && max == math.MaxFloat64 {
		return DefaultFloatType
	}
	return &floatType{min: min, max: max, inclusive: inclusive}
}

func (t *floatType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case *exactFloatType:
		return t.IsInstance(float64(ot.value))
	case *floatType:
		if t.min > ot.min {
			return false
		}
		if t.inclusive || t.inclusive == ot.inclusive {
			return t.max >= ot.max
		}
		return t.max > ot.max
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *floatType) Equals(other interface{}) bool {
	if ot, ok := other.(*floatType); ok {
		return *t == *ot
	}
	return false
}

func (t *floatType) HashCode() int {
	h := int(dgo.TiFloatRange)
	if t.min > 0 {
		h = h*31 + int(t.min)
	}
	if t.max < math.MaxInt64 {
		h = h*31 + int(t.max)
	}
	if t.inclusive {
		h *= 3
	}
	return h
}

func (t *floatType) Instance(value interface{}) bool {
	f, ok := ToFloat(value)
	return ok && t.IsInstance(f)
}

func (t *floatType) IsInstance(value float64) bool {
	if t.min <= value {
		if t.inclusive {
			return value <= t.max
		}
		return value < t.max
	}
	return false
}

func (t *floatType) Max() float64 {
	return t.max
}

func (t *floatType) Inclusive() bool {
	return t.inclusive
}

func (t *floatType) Min() float64 {
	return t.min
}

func (t *floatType) New(arg dgo.Value) dgo.Value {
	return newFloat(t, arg)
}

func (t *floatType) String() string {
	return TypeString(t)
}

func (t *floatType) ReflectType() reflect.Type {
	return reflectFloatType
}

func (t *floatType) Type() dgo.Type {
	return &metaType{t}
}

func (t *floatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloatRange
}

func (t *exactFloatType) Generic() dgo.Type {
	return DefaultFloatType
}

func (t *exactFloatType) Inclusive() bool {
	return true
}

func (t *exactFloatType) IsInstance(value float64) bool {
	return float64(t.value) == value
}

func (t *exactFloatType) Max() float64 {
	return float64(t.value)
}

func (t *exactFloatType) Min() float64 {
	return float64(t.value)
}

func (t *exactFloatType) New(arg dgo.Value) dgo.Value {
	return newFloat(t, arg)
}

func (t *exactFloatType) ReflectType() reflect.Type {
	return reflectFloatType
}

func (t *exactFloatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloatExact
}

func (t *exactFloatType) ExactValue() dgo.Value {
	return t.value
}

func (t defaultFloatType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case defaultFloatType, *exactFloatType, *floatType:
		return true
	}
	return false
}

func (t defaultFloatType) Equals(other interface{}) bool {
	_, ok := other.(defaultFloatType)
	return ok
}

func (t defaultFloatType) HashCode() int {
	return int(dgo.TiFloat)
}

func (t defaultFloatType) Inclusive() bool {
	return true
}

func (t defaultFloatType) Instance(value interface{}) bool {
	_, ok := ToFloat(value)
	return ok
}

func (t defaultFloatType) IsInstance(value float64) bool {
	return true
}

func (t defaultFloatType) Max() float64 {
	return math.MaxFloat64
}

func (t defaultFloatType) Min() float64 {
	return -math.MaxFloat64
}

func (t defaultFloatType) New(arg dgo.Value) dgo.Value {
	return newFloat(t, arg)
}

func (t defaultFloatType) ReflectType() reflect.Type {
	return reflectFloatType
}

func (t defaultFloatType) String() string {
	return TypeString(t)
}

func (t defaultFloatType) Type() dgo.Type {
	return &metaType{t}
}

func (t defaultFloatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloat
}

// Float returns the dgo.Float for the given float64
func Float(f float64) dgo.Float {
	return floatVal(f)
}

func (v floatVal) Type() dgo.Type {
	et := &exactFloatType{value: v}
	et.ExactType = et
	return et
}

func (v floatVal) CompareTo(other interface{}) (int, bool) {
	r := 0
	if ov, isFloat := ToFloat(other); isFloat {
		fv := float64(v)
		switch {
		case fv > ov:
			r = 1
		case fv < ov:
			r = -1
		}
		return r, true
	}

	if oi, isInt := ToInt(other); isInt {
		fv := float64(v)
		ov := float64(oi)
		switch {
		case fv > ov:
			r = 1
		case fv < ov:
			r = -1
		}
		return r, true
	}

	if other == Nil || other == nil {
		return 1, true
	}
	return 0, false
}

func (v floatVal) Equals(other interface{}) bool {
	f, ok := ToFloat(other)
	return ok && float64(v) == f
}

func (v floatVal) GoFloat() float64 {
	return float64(v)
}

func (v floatVal) HashCode() int {
	return int(v)
}

func (v floatVal) ReflectTo(value reflect.Value) {
	switch value.Kind() {
	case reflect.Interface:
		value.Set(reflect.ValueOf(float64(v)))
	case reflect.Ptr:
		if value.Type().Elem().Kind() == reflect.Float32 {
			gv := float32(v)
			value.Set(reflect.ValueOf(&gv))
		} else {
			gv := float64(v)
			value.Set(reflect.ValueOf(&gv))
		}
	default:
		value.SetFloat(float64(v))
	}
}

func (v floatVal) String() string {
	return util.Ftoa(float64(v))
}

func (v floatVal) ToFloat() float64 {
	return float64(v)
}

func (v floatVal) ToInt() int64 {
	return int64(v)
}

// ToFloat returns the given value as a float64 if, and only if, the value is a float32 or float64. An
// additional boolean is returned to indicate if that was the case or not.
func ToFloat(value interface{}) (v float64, ok bool) {
	ok = true
	switch value := value.(type) {
	case floatVal:
		v = float64(value)
	case float64:
		v = value
	case float32:
		v = float64(value)
	default:
		ok = false
	}
	return
}

func newFloat(t dgo.Type, arg dgo.Value) (f dgo.Float) {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`float`, 1, 1)
		arg = args.Get(0)
	}
	f = Float(floatFromConvertible(arg))
	if !t.Instance(f) {
		panic(IllegalAssignment(t, f))
	}
	return f
}

func floatFromConvertible(from dgo.Value) float64 {
	switch from := from.(type) {
	case dgo.Float:
		return from.GoFloat()
	case dgo.Integer:
		return float64(from.GoInt())
	case *timeVal:
		return from.SecondsWithFraction()
	case dgo.Boolean:
		if from.GoBool() {
			return 1
		}
		return 0
	case dgo.String:
		if i, err := strconv.ParseFloat(from.GoString(), 64); err == nil {
			return i
		}
	}
	panic(fmt.Errorf(`the value '%s' cannot be converted to a float`, from))
}
