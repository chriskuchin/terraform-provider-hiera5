package internal

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
)

type (
	// intVal is an int64 that implements the dgo.Value interface
	intVal int64

	defaultIntegerType int

	integerType struct {
		min       dgo.Integer
		max       dgo.Integer
		inclusive bool
	}
)

// DefaultIntegerType is the unconstrained Integer type
const DefaultIntegerType = defaultIntegerType(0)

var reflectIntegerType = reflect.TypeOf(int64(0))

// Integer64Type returns a dgo.Integer64Type that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func Integer64Type(min, max int64, inclusive bool) dgo.IntegerType {
	if min == max {
		if !inclusive {
			panic(catch.Error(`non inclusive range cannot have equal min and max`))
		}
		return intVal(min).Type().(dgo.IntegerType)
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	var minV dgo.Integer
	var maxV dgo.Integer
	if min != math.MinInt64 {
		minV = intVal(min)
	}
	if max != math.MaxInt64 {
		maxV = intVal(max)
	}
	if minV == nil && maxV == nil {
		return DefaultIntegerType
	}
	return &integerType{min: minV, max: maxV, inclusive: inclusive}
}

// IntegerType returns a dgo.Integer64Type that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end. The IntegerType.ReflectType() returns
// the *big.Int type.
func IntegerType(min, max dgo.Integer, inclusive bool) dgo.IntegerType {
	if min != nil && max != nil {
		cmp, _ := min.CompareTo(max)
		if cmp == 0 {
			if !inclusive {
				panic(catch.Error(`non inclusive range cannot have equal min and max`))
			}
			return min.(dgo.IntegerType)
		}
		if cmp > 0 {
			t := max
			max = min
			min = t
		}
	} else if min == nil && max == nil {
		return DefaultIntegerType
	}
	_, useBig := min.(dgo.BigInt)
	if !useBig {
		_, useBig = max.(dgo.BigInt)
	}
	if useBig {
		return &bigIntType{integerType{min: min, max: max, inclusive: inclusive}}
	}
	return &integerType{min: min, max: max, inclusive: inclusive}
}

func (t *integerType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case defaultIntegerType:
		return false
	case dgo.IntegerType:
		if t.min != nil {
			om := ot.Min()
			if om == nil {
				return false
			}
			cmp, _ := t.min.CompareTo(om)
			if cmp > 0 {
				return false
			}
		}
		if mm := t.max; mm != nil {
			om := ot.Max()
			if om == nil {
				return false
			}
			if t.Inclusive() {
				mm = mm.Inc()
			}
			if ot.Inclusive() {
				om = om.Inc()
			}
			cmp, _ := mm.CompareTo(om)
			if cmp < 0 {
				return false
			}
		}
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *integerType) Equals(other interface{}) bool {
	ot, ok := other.(dgo.IntegerType)
	return ok && t.inclusive == ot.Inclusive() && equals(nil, t.min, ot.Min()) && equals(nil, t.max, ot.Max())
}

func (t *integerType) HashCode() dgo.Hash {
	h := dgo.Hash(dgo.TiIntegerRange)
	if t.min != nil {
		h = h*31 + t.min.HashCode()
	}
	if t.max != nil {
		h = h*31 + t.max.HashCode()
	}
	if t.inclusive {
		h *= 3
	}
	return h
}

func (t *integerType) Inclusive() bool {
	return t.inclusive
}

func (t *integerType) Instance(value interface{}) bool {
	yes := false
	switch ov := value.(type) {
	case intVal:
		yes = t.isInstance(ov)
	case int:
		yes = t.isInstance(intVal(ov))
	case *bigIntVal:
		yes = t.isInstance(ov)
	case *big.Int:
		yes = t.isInstance(&bigIntVal{ov})
	default:
		var iv int64
		iv, yes = ToInt(value)
		yes = yes && t.isInstance(intVal(iv))
	}
	return yes
}

func (t *integerType) isInstance(i dgo.Integer) bool {
	if t.min != nil {
		cmp, ok := t.min.CompareTo(i)
		if !ok || cmp > 0 {
			return false
		}
	}
	if t.max != nil {
		cmp, ok := t.max.CompareTo(i)
		if !ok || cmp < 0 || cmp == 0 && !t.inclusive {
			return false
		}
	}
	return true
}

func (t *integerType) Max() dgo.Integer {
	return t.max
}

func (t *integerType) Min() dgo.Integer {
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
	return MetaType(t)
}

func (t *integerType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiIntegerRange
}

func (t defaultIntegerType) Assignable(other dgo.Type) bool {
	_, ok := other.(dgo.IntegerType)
	return ok || CheckAssignableTo(nil, other, t)
}

func (t defaultIntegerType) Equals(other interface{}) bool {
	_, ok := other.(defaultIntegerType)
	return ok
}

func (t defaultIntegerType) HashCode() dgo.Hash {
	return dgo.Hash(dgo.TiInteger)
}

func (t defaultIntegerType) Instance(value interface{}) bool {
	switch value.(type) {
	case dgo.Integer, *big.Int, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		return true
	}
	return false
}

func (t defaultIntegerType) Inclusive() bool {
	return true
}

func (t defaultIntegerType) Max() dgo.Integer {
	return nil
}

func (t defaultIntegerType) Min() dgo.Integer {
	return nil
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
	return MetaType(t)
}

func (t defaultIntegerType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiInteger
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
	return &anyOfType{slice: ts}
}

// Integer returns the dgo.Integer for the given int64
func Integer(v int64) dgo.Integer {
	return intVal(v)
}

func (v intVal) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v intVal) CompareTo(other interface{}) (int, bool) {
	r := 0
	ok := true

	mv := int64(v)
	compare64 := func(fv int64) {
		switch {
		case mv > fv:
			r = 1
		case mv < fv:
			r = -1
		}
	}

	compareBig := func(ov *big.Int) {
		if ov.IsInt64() {
			compare64(ov.Int64())
		} else {
			r = -ov.Sign()
		}
	}

	switch ov := other.(type) {
	case nil, nilValue:
		r = 1
	case intVal:
		compare64(int64(ov))
	case int:
		compare64(int64(ov))
	case int64:
		compare64(ov)
	case uint:
		if ov > math.MaxInt64 {
			r = -1
		} else {
			compare64(int64(ov))
		}
	case uint64:
		if ov > math.MaxInt64 {
			r = -1
		} else {
			compare64(int64(ov))
		}
	case *bigIntVal:
		compareBig(ov.Int)
	case *big.Int:
		compareBig(ov)
	case dgo.Float:
		r, ok = v.Float().CompareTo(ov)
	default: // all other int types
		var iv int64
		iv, ok = ToInt(other)
		if ok {
			compare64(iv)
		}
	}
	return r, ok
}

func (v intVal) Dec() dgo.Integer {
	return v - 1
}

func (v intVal) Equals(other interface{}) bool {
	i, ok := ToInt(other)
	return ok && int64(v) == i
}

func (v intVal) Float() dgo.Float {
	return floatVal(v)
}

func (v intVal) Format(s fmt.State, format rune) {
	doFormat(int64(v), s, format)
}

func (v intVal) Generic() dgo.Type {
	return DefaultIntegerType
}

func (v intVal) GoInt() int64 {
	return int64(v)
}

func (v intVal) HashCode() dgo.Hash {
	return dgo.Hash(v ^ (v >> 32))
}

func (v intVal) Inc() dgo.Integer {
	return v + 1
}

func (v intVal) Inclusive() bool {
	return true
}

func (v intVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v intVal) Integer() dgo.Integer {
	return v
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

func (v intVal) Max() dgo.Integer {
	return v
}

func (v intVal) Min() dgo.Integer {
	return v
}

func (v intVal) New(arg dgo.Value) dgo.Value {
	return newInt(v, arg)
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

func (v intVal) ReflectType() reflect.Type {
	return reflectIntegerType
}

func (v intVal) String() string {
	return TypeString(v)
}

func (v intVal) ToBigInt() *big.Int {
	return big.NewInt(int64(v))
}

func (v intVal) ToBigFloat() *big.Float {
	return new(big.Float).SetInt64(int64(v))
}

func (v intVal) ToFloat() (float64, bool) {
	return float64(v), true
}

func (v intVal) ToInt() (int64, bool) {
	return int64(v), true
}

func (v intVal) Type() dgo.Type {
	return v
}

func (v intVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiIntegerExact
}

// ToInt returns the given value as a *big.Int if the value type is one of the go int types, a *big.Int,
// a dgo.Integer or a dgo.BigInt, and the value fits an int64 value. An additional boolean is returned
// to indicate if that was the case or not.
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
		if value <= math.MaxInt64 {
			v = int64(value)
		} else {
			ok = false
		}
	case uint64:
		if value <= math.MaxInt64 {
			v = int64(value)
		} else {
			ok = false
		}
	case uint32:
		v = int64(value)
	case uint16:
		v = int64(value)
	case uint8:
		v = int64(value)
	case *big.Int:
		if value.IsInt64() {
			v = value.Int64()
		} else {
			ok = false
		}
	case dgo.BigInt:
		v, ok = value.ToInt()
	default:
		ok = false
	}
	return v, ok
}

func unsignedToInteger(v uint64) dgo.Integer {
	if v <= math.MaxInt64 {
		return intVal(int64(v))
	}
	return &bigIntVal{new(big.Int).SetUint64(v)}
}

var radixType = IntEnumType([]int{0, 2, 8, 10, 16})

func newInt(t dgo.Type, arg dgo.Value) (i dgo.Integer) {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`int`, 1, 2)
		if args.Len() == 2 {
			i = intFromConvertible(args.Get(0), int(args.Arg(`int`, 1, radixType).(dgo.Integer).GoInt()))
		} else {
			i = intFromConvertible(args.Get(0), 0)
		}
	} else {
		i = intFromConvertible(arg, 0)
	}
	if !t.Instance(i) {
		panic(catch.Error(IllegalAssignment(t, i)))
	}
	return i
}

func intFromConvertible(from dgo.Value, radix int) dgo.Integer {
	switch from := from.(type) {
	case dgo.Number:
		return from.Integer()
	case dgo.Boolean:
		if from.GoBool() {
			return intVal(1)
		}
		return intVal(0)
	case dgo.String:
		s := from.GoString()
		i, err := strconv.ParseInt(s, radix, 64)
		if err == nil {
			return Integer(i)
		}
		numErr, ok := err.(*strconv.NumError)
		if ok && numErr.Err == strconv.ErrRange {
			var bi *big.Int
			if bi, ok = new(big.Int).SetString(s, radix); ok {
				return BigInt(bi)
			}
		}
	}
	panic(catch.Error(`the value '%v' cannot be converted to an int`, from))
}
