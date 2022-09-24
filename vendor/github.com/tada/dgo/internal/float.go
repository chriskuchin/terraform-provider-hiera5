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
	// floatVal is a float64 that implements the dgo.Value interface
	floatVal float64

	defaultFloatType int

	floatType struct {
		min       dgo.Float
		max       dgo.Float
		inclusive bool
	}
)

// DefaultFloatType is the unconstrained floatVal type
const DefaultFloatType = defaultFloatType(0)

var reflectFloatType = reflect.TypeOf(float64(0))

// Float64Type returns a dgo.Float64Type that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func Float64Type(min, max float64, inclusive bool) dgo.FloatType {
	if min == max {
		if !inclusive {
			panic(catch.Error(`non inclusive range cannot have equal min and max`))
		}
		return floatVal(min).Type().(dgo.FloatType)
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	var minV dgo.Float
	var maxV dgo.Float
	if min != -math.MaxFloat64 {
		minV = floatVal(min)
	}
	if max != math.MaxFloat64 {
		maxV = floatVal(max)
	}
	if minV == nil && maxV == nil {
		return DefaultFloatType
	}
	return &floatType{min: minV, max: maxV, inclusive: inclusive}
}

// FloatType returns a dgo.Float64Type that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end. The FloatType.ReflectType() returns
// the *big.Float type.
func FloatType(min, max dgo.Float, inclusive bool) dgo.FloatType {
	if min != nil && max != nil {
		cmp, _ := min.CompareTo(max)
		if cmp == 0 {
			if !inclusive {
				panic(catch.Error(`non inclusive range cannot have equal min and max`))
			}
			return min.(dgo.FloatType)
		}
		if cmp > 0 {
			t := max
			max = min
			min = t
		}
	} else if min == nil && max == nil {
		return DefaultFloatType
	}
	_, useBig := min.(dgo.BigFloat)
	if !useBig {
		_, useBig = max.(dgo.BigFloat)
	}
	if useBig {
		return &bigFloatType{floatType{min: min, max: max, inclusive: inclusive}}
	}
	return &floatType{min: min, max: max, inclusive: inclusive}
}

func (t *floatType) Assignable(other dgo.Type) bool {
	switch ot := other.(type) {
	case defaultFloatType:
		return false
	case dgo.FloatType:
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
		if t.max != nil {
			om := ot.Max()
			if om == nil {
				return false
			}
			cmp, _ := t.max.CompareTo(om)
			if cmp < 0 || cmp == 0 && !(t.inclusive || !ot.Inclusive()) {
				return false
			}
		}
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *floatType) Equals(other interface{}) bool {
	ot, ok := other.(dgo.FloatType)
	return ok && t.inclusive == ot.Inclusive() && equals(nil, t.min, ot.Min()) && equals(nil, t.max, ot.Max())
}

func (t *floatType) HashCode() dgo.Hash {
	h := dgo.Hash(dgo.TiFloatRange)
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

func (t *floatType) Inclusive() bool {
	return t.inclusive
}

func (t *floatType) Instance(value interface{}) bool {
	yes := false
	switch ov := value.(type) {
	case floatVal:
		yes = t.isInstance(ov)
	case float64:
		yes = t.isInstance(floatVal(ov))
	case float32:
		yes = t.isInstance(floatVal(ov))
	case *bigFloatVal:
		yes = t.isInstance(ov)
	case *big.Float:
		yes = t.isInstance(&bigFloatVal{ov})
	}
	return yes
}

func (t *floatType) isInstance(f dgo.Float) bool {
	if t.min != nil {
		cmp, ok := t.min.CompareTo(f)
		if !ok || cmp > 0 {
			return false
		}
	}
	if t.max != nil {
		cmp, ok := t.max.CompareTo(f)
		if !ok || cmp < 0 || cmp == 0 && !t.inclusive {
			return false
		}
	}
	return true
}

func (t *floatType) Max() dgo.Float {
	return t.max
}

func (t *floatType) Min() dgo.Float {
	return t.min
}

func (t *floatType) New(arg dgo.Value) dgo.Value {
	return newFloat(t, arg)
}

func (t *floatType) ReflectType() reflect.Type {
	return reflectFloatType
}

func (t *floatType) String() string {
	return TypeString(t)
}

func (t *floatType) Type() dgo.Type {
	return MetaType(t)
}

func (t *floatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloatRange
}

func (t defaultFloatType) Assignable(other dgo.Type) bool {
	_, ok := other.(dgo.FloatType)
	return ok || CheckAssignableTo(nil, other, t)
}

func (t defaultFloatType) Equals(other interface{}) bool {
	_, ok := other.(defaultFloatType)
	return ok
}

func (t defaultFloatType) HashCode() dgo.Hash {
	return dgo.Hash(dgo.TiFloat)
}

func (t defaultFloatType) Instance(value interface{}) bool {
	switch value.(type) {
	case dgo.Float, *big.Float, float64, float32:
		return true
	}
	return false
}

func (t defaultFloatType) Inclusive() bool {
	return true
}

func (t defaultFloatType) Max() dgo.Float {
	return nil
}

func (t defaultFloatType) Min() dgo.Float {
	return nil
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
	return MetaType(t)
}

func (t defaultFloatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloat
}

// Float returns the dgo.Float for the given float64
func Float(f float64) dgo.Float {
	return floatVal(f)
}

func (v floatVal) Type() dgo.Type {
	return v
}

func (v floatVal) CompareTo(other interface{}) (int, bool) {
	r := 0
	ok := true

	mv := float64(v)
	compare64 := func(fv float64) {
		switch {
		case mv > fv:
			r = 1
		case mv < fv:
			r = -1
		}
	}

	compareBig := func(ov *big.Float) {
		fv, a := ov.Float64()
		switch {
		case a == big.Above && math.IsInf(fv, 1):
			// number is too big and positive
			r = 1
		case a == big.Below && math.IsInf(fv, -1):
			// number is too big and negative
			r = -1
		case a == big.Below && fv == 0:
			// number is too small and positive
			if mv < 0 {
				r = -1
			} else {
				r = 1
			}
		case a == big.Above && fv == -0:
			// number is too small and negative
			if mv > 0 {
				r = 1
			} else {
				r = -1
			}
		default:
			compare64(fv)
		}
	}

	switch ov := other.(type) {
	case nil, nilValue:
		r = 1
	case floatVal:
		compare64(float64(ov))
	case *bigFloatVal:
		compareBig(ov._bf)
	case *bigIntVal:
		compareBig(new(big.Float).SetInt(ov.GoBigInt()))
	case dgo.Number:
		r, ok = v.CompareTo(ov.Float())
	default:
		var n dgo.Number
		if n, ok = Value(ov).(dgo.Number); ok {
			r, ok = v.CompareTo(n)
		}
	}
	return r, ok
}

func (v floatVal) Equals(other interface{}) bool {
	f, ok := ToFloat(other)
	return ok && float64(v) == f
}

func (v floatVal) Float() dgo.Float {
	return v
}

func (v floatVal) Format(s fmt.State, format rune) {
	doFormat(float64(v), s, format)
}

func (v floatVal) GoFloat() float64 {
	return float64(v)
}

func (v floatVal) HashCode() dgo.Hash {
	return dgo.Hash(v)
}

func (v floatVal) Integer() dgo.Integer {
	return intVal(v)
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
	return TypeString(v)
}

func (v floatVal) ToFloat() (float64, bool) {
	return float64(v), true
}

func (v floatVal) ToInt() (int64, bool) {
	return int64(v), true
}

func (v floatVal) ToBigInt() *big.Int {
	return big.NewInt(int64(v))
}

func (v floatVal) ToBigFloat() *big.Float {
	return big.NewFloat(float64(v))
}

func (v floatVal) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v floatVal) Generic() dgo.Type {
	return DefaultFloatType
}

func (v floatVal) Inclusive() bool {
	return true
}

func (v floatVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v floatVal) Max() dgo.Float {
	return v
}

func (v floatVal) Min() dgo.Float {
	return v
}

func (v floatVal) New(arg dgo.Value) dgo.Value {
	return newFloat(v, arg)
}

func (v floatVal) ReflectType() reflect.Type {
	return reflectFloatType
}

func (v floatVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloatExact
}

// ToFloat returns the given value as a float64 if the value is a *big.Float, float32, float64, dgo.Float,
// or dgo.BigFloat and the value can be converted to a float64 with exact accuracy.
// An additional boolean is returned to indicate if that was the case or not.
func ToFloat(value interface{}) (v float64, ok bool) {
	ok = true
	switch value := value.(type) {
	case floatVal:
		v = float64(value)
	case float64:
		v = value
	case float32:
		v = float64(value)
	case *bigFloatVal:
		v, ok = demoteToFloat64(value._bf)
	case *big.Float:
		v, ok = demoteToFloat64(value)
	default:
		ok = false
	}
	return
}

func demoteToFloat64(bf *big.Float) (float64, bool) {
	f, a := bf.Float64()
	if a == big.Below && (f == 0 || math.IsInf(f, -1)) || a == big.Above && (f == -0 || math.IsInf(f, 1)) {
		return 0, false
	}
	return f, true
}

func demoteToInt64(bf *big.Float) (int64, bool) {
	i, a := bf.Int64()
	if a == big.Below && i == math.MaxInt64 || a == big.Above && i == math.MinInt64 {
		return 0, false
	}
	return i, true
}

func newFloat(t dgo.Type, arg dgo.Value) (f dgo.Float) {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`float`, 1, 1)
		arg = args.Get(0)
	}
	f = floatFromConvertible(arg)
	if !t.Instance(f) {
		panic(catch.Error(IllegalAssignment(t, f)))
	}
	return f
}

func floatFromConvertible(from dgo.Value) dgo.Float {
	switch from := from.(type) {
	case dgo.Number:
		return from.Float()
	case dgo.Boolean:
		if from.GoBool() {
			return floatVal(1)
		}
		return floatVal(0)
	case dgo.String:
		if f, err := strconv.ParseFloat(from.GoString(), 64); err == nil {
			return floatVal(f)
		}
	}
	panic(fmt.Errorf(`the value '%s' cannot be converted to a float`, from))
}
