package internal

import (
	"math"
	"math/big"
	"reflect"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
)

type (
	// bf type anonymizes the big.Float to avoid collisions between a Float attribute and the Float() function while
	// the bigFloatVal still inherits all functions from big.Float since the actual field is unnamed.
	_bf = *big.Float

	bigFloatVal struct {
		_bf
	}

	defaultBigFloatType struct {
		defaultFloatType
	}

	bigFloatType struct {
		floatType
	}
)

// DefaultBigFloatType is the unconstrained Integer type
var DefaultBigFloatType = &defaultBigFloatType{}

var reflectBigFloatType = reflect.TypeOf(&big.Float{})

func (t *defaultBigFloatType) New(arg dgo.Value) dgo.Value {
	return newBigFloat(t, arg)
}

func (t *defaultBigFloatType) ReflectType() reflect.Type {
	return reflectBigFloatType
}

func (t *defaultBigFloatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBigFloat
}

func (t *bigFloatType) New(arg dgo.Value) dgo.Value {
	return newBigFloat(t, arg)
}

func (t *bigFloatType) ReflectType() reflect.Type {
	return reflectBigFloatType
}

func (t *bigFloatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBigFloatRange
}

// BigFloat returns the dgo.BigFloat for the given *big.Float
func BigFloat(v *big.Float) dgo.BigFloat {
	return &bigFloatVal{v}
}

func (v *bigFloatVal) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v *bigFloatVal) CompareTo(other interface{}) (int, bool) {
	r := 0
	ok := true

	compare64 := func(ov float64) {
		r = v._bf.Cmp(big.NewFloat(ov))
	}

	switch ov := other.(type) {
	case nil, nilValue:
		r = 1
	case *bigFloatVal:
		r = v.Cmp(ov._bf)
	case floatVal:
		compare64(float64(ov))
	case *big.Float:
		r = v.Cmp(ov)
	case float64:
		compare64(ov)
	case float32:
		compare64(float64(ov))
	case *big.Int:
		r = v.Cmp(new(big.Float).SetInt(ov))
	case dgo.Number:
		r, ok = v.CompareTo(ov.Float())
	default:
		var i int64
		if i, ok = ToInt(ov); ok {
			compare64(float64(i))
		}
	}
	return r, ok
}

func (v *bigFloatVal) Equals(other interface{}) bool {
	yes := false
	switch ov := other.(type) {
	case *bigFloatVal:
		yes = v.Cmp(ov._bf) == 0
	case *big.Float:
		yes = v.Cmp(ov) == 0
	case floatVal:
		yes = v.Cmp(big.NewFloat(float64(ov))) == 0
	case float64:
		yes = v.Cmp(big.NewFloat(ov)) == 0
	case float32:
		yes = v.Cmp(big.NewFloat(float64(ov))) == 0
	}
	return yes
}

func (v *bigFloatVal) Float() dgo.Float {
	return v
}

func (v *bigFloatVal) Generic() dgo.Type {
	return DefaultBigFloatType
}

func (v *bigFloatVal) GoBigFloat() *big.Float {
	return v._bf
}

func (v *bigFloatVal) GoFloat() float64 {
	if f, ok := v.ToFloat(); ok {
		return f
	}
	panic(catch.Error(`BigFloat.ToFloat(): value %f cannot fit into a float64`, v))
}

func (v *bigFloatVal) HashCode() dgo.Hash {
	return bigFloatHash(v._bf)
}

func (v *bigFloatVal) Inclusive() bool {
	return true
}

func (v *bigFloatVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v *bigFloatVal) Integer() dgo.Integer {
	bi, _ := v.Int(nil)
	return &bigIntVal{bi}
}

func (v *bigFloatVal) Max() dgo.Float {
	return v
}

func (v *bigFloatVal) Min() dgo.Float {
	return v
}

func (v *bigFloatVal) New(arg dgo.Value) dgo.Value {
	return newBigFloat(v, arg)
}

func (v *bigFloatVal) ReflectTo(value reflect.Value) {
	rv := reflect.ValueOf(v._bf)
	k := value.Kind()
	if !(k == reflect.Ptr || k == reflect.Interface) {
		rv = rv.Elem()
	}
	value.Set(rv)
}

func (v *bigFloatVal) ReflectType() reflect.Type {
	return reflectBigFloatType
}

func (v *bigFloatVal) String() string {
	return TypeString(v)
}

func (v *bigFloatVal) ToBigFloat() *big.Float {
	return v._bf
}

func (v *bigFloatVal) ToBigInt() *big.Int {
	bi, _ := v.Int(nil)
	return bi
}

func (v *bigFloatVal) ToFloat() (float64, bool) {
	return demoteToFloat64(v._bf)
}

func (v *bigFloatVal) ToInt() (int64, bool) {
	return demoteToInt64(v._bf)
}

func (v *bigFloatVal) Type() dgo.Type {
	return v
}

func (v *bigFloatVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBigFloatExact
}

func bigFloatHash(v *big.Float) dgo.Hash {
	ge, _ := v.GobEncode()
	return bytesHash(ge)
}

func bigFloatFromConvertible(from dgo.Value, prec uint) dgo.Float {
	switch from := from.(type) {
	case dgo.Number:
		f := from.Float()
		if _, ok := f.(dgo.BigFloat); ok {
			return f
		}
		return &bigFloatVal{f.ToBigFloat()}
	case dgo.Boolean:
		if from.GoBool() {
			return &bigFloatVal{big.NewFloat(1)}
		}
		return &bigFloatVal{big.NewFloat(0)}
	case dgo.String:
		if f, _, err := big.ParseFloat(from.GoString(), 0, prec, big.ToNearestEven); err == nil {
			return &bigFloatVal{f}
		}
	}
	panic(catch.Error(`the value '%s' cannot be converted to a big float`, from))
}

var precType = Integer64Type(0, math.MaxUint32, true)

func newBigFloat(t dgo.Type, arg dgo.Value) (f dgo.Float) {
	prec := uint(0)
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`big`, 1, 2)
		arg = args.Get(0)
		if args.Len() > 1 {
			prec = uint(args.Arg(`int`, 1, precType).(dgo.Integer).GoInt())
		}
	}
	f = bigFloatFromConvertible(arg, prec)
	if !t.Instance(f) {
		panic(catch.Error(IllegalAssignment(t, f)))
	}
	return f
}
