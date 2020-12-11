package internal

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// intVal is an int64 that implements the dgo.Value interface
	intVal int64

	defaultIntegerType int

	exactIntegerType struct {
		exactType
		value intVal
	}

	integerType struct {
		min       int64
		max       int64
		inclusive bool
	}
)

// DefaultIntegerType is the unconstrained Integer type
const DefaultIntegerType = defaultIntegerType(0)

var reflectIntegerType = reflect.TypeOf(int64(0))

// IntegerType returns a dgo.IntegerType that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func IntegerType(min, max int64, inclusive bool) dgo.IntegerType {
	if min == max {
		if !inclusive {
			panic(fmt.Errorf(`non inclusive range cannot have equal min and max`))
		}
		return intVal(min).Type().(dgo.IntegerType)
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	if min == math.MinInt64 && max == math.MaxInt64 {
		return DefaultIntegerType
	}
	return &integerType{min: min, max: max, inclusive: inclusive}
}

// IntEnumType returns a Type that represents any of the given integers
func IntEnumType(ints []int) dgo.Type {
	switch len(ints) {
	case 0:
		return &notType{DefaultAnyType}
	case 1:
		return intVal(ints[0]).Type()
	}
	ts := make([]dgo.Value, len(ints))
	for i := range ints {
		ts[i] = intVal(ints[i]).Type()
	}
	return &anyOfType{slice: ts, frozen: true}
}

func (t *integerType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case *exactIntegerType:
		return t.IsInstance(int64(ot.value))
	case *integerType:
		if t.min > ot.min {
			return false
		}
		mm := t.max
		if !t.inclusive {
			mm--
		}
		om := ot.max
		if !ot.inclusive {
			om--
		}
		return mm >= om
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *integerType) Equals(other interface{}) bool {
	if ot, ok := other.(*integerType); ok {
		return *t == *ot
	}
	return false
}

func (t *integerType) HashCode() int {
	h := int(dgo.TiIntegerRange)
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

func (t *integerType) Instance(value interface{}) bool {
	if ov, ok := ToInt(value); ok {
		return t.IsInstance(ov)
	}
	return false
}

func (t *integerType) IsInstance(value int64) bool {
	if t.min <= value {
		if t.inclusive {
			return value <= t.max
		}
		return value < t.max
	}
	return false
}

func (t *integerType) Inclusive() bool {
	return t.inclusive
}

func (t *integerType) Max() int64 {
	return t.max
}

func (t *integerType) Min() int64 {
	return t.min
}

func (t *integerType) New(arg dgo.Value) dgo.Value {
	return newInt(t, arg)
}

func (t *integerType) ReflectType() reflect.Type {
	return reflectIntegerType
}

func (t *integerType) String() string {
	return TypeString(t)
}

func (t *integerType) Type() dgo.Type {
	return &metaType{t}
}

func (t *integerType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiIntegerRange
}

func (t *exactIntegerType) Generic() dgo.Type {
	return DefaultIntegerType
}

func (t *exactIntegerType) Inclusive() bool {
	return true
}

func (t *exactIntegerType) IsInstance(value int64) bool {
	return int64(t.value) == value
}

func (t *exactIntegerType) Max() int64 {
	return int64(t.value)
}

func (t *exactIntegerType) Min() int64 {
	return int64(t.value)
}

func (t *exactIntegerType) New(arg dgo.Value) dgo.Value {
	return newInt(t, arg)
}

func (t *exactIntegerType) ReflectType() reflect.Type {
	return reflectIntegerType
}

func (t *exactIntegerType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiIntegerExact
}

func (t *exactIntegerType) ExactValue() dgo.Value {
	return t.value
}

func (t defaultIntegerType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case defaultIntegerType, *exactIntegerType, *integerType:
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t defaultIntegerType) Equals(other interface{}) bool {
	_, ok := other.(defaultIntegerType)
	return ok
}

func (t defaultIntegerType) HashCode() int {
	return int(dgo.TiInteger)
}

func (t defaultIntegerType) Instance(value interface{}) bool {
	switch value.(type) {
	case dgo.Integer, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return true
	}
	return false
}

func (t defaultIntegerType) Inclusive() bool {
	return true
}

func (t defaultIntegerType) IsInstance(value int64) bool {
	return true
}

func (t defaultIntegerType) Max() int64 {
	return math.MaxInt64
}

func (t defaultIntegerType) Min() int64 {
	return math.MinInt64
}

func (t defaultIntegerType) New(arg dgo.Value) dgo.Value {
	return newInt(t, arg)
}

func (t defaultIntegerType) ReflectType() reflect.Type {
	return reflectIntegerType
}

func (t defaultIntegerType) String() string {
	return TypeString(t)
}

func (t defaultIntegerType) Type() dgo.Type {
	return &metaType{t}
}

func (t defaultIntegerType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiInteger
}

// Integer returns the dgo.Integer for the given int64
func Integer(v int64) dgo.Integer {
	return intVal(v)
}

func (v intVal) CompareTo(other interface{}) (int, bool) {
	r := 0
	if oi, isInt := ToInt(other); isInt {
		mv := int64(v)
		switch {
		case mv > oi:
			r = 1
		case mv < oi:
			r = -1
		}
		return r, true
	}

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

	if other == Nil || other == nil {
		return 1, true
	}
	return 0, false
}

func (v intVal) Equals(other interface{}) bool {
	i, ok := ToInt(other)
	return ok && int64(v) == i
}

func (v intVal) GoInt() int64 {
	return int64(v)
}

func (v intVal) HashCode() int {
	return int(v ^ (v >> 32))
}

func (v intVal) ReflectTo(value reflect.Value) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(int64(v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(uint64(v))
	case reflect.Ptr:
		value.Set(v.intPointer(value.Type().Elem().Kind()))
	default:
		value.Set(reflect.ValueOf(int64(v)))
	}
}

func (v intVal) intPointer(kind reflect.Kind) reflect.Value {
	var p reflect.Value
	switch kind {
	case reflect.Int:
		gv := int(v)
		p = reflect.ValueOf(&gv)
	case reflect.Int8:
		gv := int8(v)
		p = reflect.ValueOf(&gv)
	case reflect.Int16:
		gv := int16(v)
		p = reflect.ValueOf(&gv)
	case reflect.Int32:
		gv := int32(v)
		p = reflect.ValueOf(&gv)
	case reflect.Uint:
		gv := uint(v)
		p = reflect.ValueOf(&gv)
	case reflect.Uint8:
		gv := uint8(v)
		p = reflect.ValueOf(&gv)
	case reflect.Uint16:
		gv := uint16(v)
		p = reflect.ValueOf(&gv)
	case reflect.Uint32:
		gv := uint32(v)
		p = reflect.ValueOf(&gv)
	case reflect.Uint64:
		gv := uint64(v)
		p = reflect.ValueOf(&gv)
	default:
		gv := int64(v)
		p = reflect.ValueOf(&gv)
	}
	return p
}

func (v intVal) String() string {
	return strconv.Itoa(int(v))
}

func (v intVal) ToFloat() float64 {
	return float64(v)
}

func (v intVal) ToInt() int64 {
	return int64(v)
}

func (v intVal) Type() dgo.Type {
	et := &exactIntegerType{value: v}
	et.ExactType = et
	return et
}

// ToInt returns the given value as a int64 if, and only if, the value type is one of the go int types. An
// additional boolean is returned to indicate if that was the case or not.
func ToInt(value interface{}) (int64, bool) {
	ok := true
	v := int64(0)
	switch value := value.(type) {
	case intVal:
		v = int64(value)
	case int:
		v = int64(value)
	case int64:
		v = value
	case int32:
		v = int64(value)
	case int16:
		v = int64(value)
	case int8:
		v = int64(value)
	case uint:
		if value == math.MaxUint64 {
			panic(fmt.Errorf(`value %d overflows int64`, value))
		}
		v = int64(value)
	case uint64:
		if value == math.MaxUint64 {
			panic(fmt.Errorf(`value %d overflows int64`, value))
		}
		v = int64(value)
	case uint32:
		v = int64(value)
	case uint16:
		v = int64(value)
	case uint8:
		v = int64(value)
	default:
		ok = false
	}
	return v, ok
}

var radixType = IntEnumType([]int{2, 8, 10, 16})

func newInt(t dgo.Type, arg dgo.Value) (i dgo.Integer) {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`int`, 1, 2)
		if args.Len() == 2 {
			i = Integer(intFromConvertible(args.Get(0), int(args.Arg(`int`, 1, radixType).(dgo.Integer).GoInt())))
		} else {
			i = Integer(intFromConvertible(args.Get(0), 10))
		}
	} else {
		i = Integer(intFromConvertible(arg, 10))
	}
	if !t.Instance(i) {
		panic(IllegalAssignment(t, i))
	}
	return i
}

func intFromConvertible(from dgo.Value, radix int) int64 {
	switch from := from.(type) {
	case dgo.Integer:
		return from.GoInt()
	case dgo.Float:
		return int64(from.GoFloat())
	case *timeVal:
		return from.GoTime().Unix()
	case dgo.Boolean:
		if from.GoBool() {
			return 1
		}
		return 0
	case dgo.String:
		if i, err := strconv.ParseInt(from.GoString(), radix, 64); err == nil {
			return i
		}
	}
	panic(fmt.Errorf(`the value '%s' cannot be converted to an int`, from))
}
