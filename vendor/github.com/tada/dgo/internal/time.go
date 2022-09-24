package internal

import (
	"math"
	"math/big"
	"reflect"
	"time"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
)

type (
	timeType int

	timeVal struct {
		time.Time
	}
)

// DefaultTimeType is the unconstrainted Time type
const DefaultTimeType = timeType(0)

var reflectTimeType = reflect.TypeOf(time.Time{})

func (t timeType) Assignable(ot dgo.Type) bool {
	switch ot.(type) {
	case *timeVal, timeType:
		return true
	}
	return CheckAssignableTo(nil, ot, t)
}

func (t timeType) Equals(v interface{}) bool {
	return t == v
}

func (t timeType) HashCode() dgo.Hash {
	return dgo.Hash(dgo.TiTime)
}

func (t timeType) Instance(v interface{}) bool {
	switch v.(type) {
	case *timeVal, *time.Time, time.Time:
		return true
	}
	return false
}

func (t timeType) New(arg dgo.Value) dgo.Value {
	return newTime(t, arg)
}

func (t timeType) ReflectType() reflect.Type {
	return reflectTimeType
}

func (t timeType) String() string {
	return TypeString(t)
}

func (t timeType) Type() dgo.Type {
	return MetaType(t)
}

func (t timeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiTime
}

func newTime(t dgo.Type, arg dgo.Value) dgo.Time {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`time`, 1, 1)
		arg = args.Get(0)
	}
	var tv dgo.Time
	switch arg := arg.(type) {
	case dgo.Time:
		tv = arg
	case dgo.Integer:
		tv = Time(time.Unix(arg.GoInt(), 0))
	case dgo.Float:
		s, f := math.Modf(arg.GoFloat())
		tv = Time(time.Unix(int64(s), int64(f*1000000000.0)))
	case dgo.String:
		tv = TimeFromString(arg.GoString())
	default:
		panic(illegalArgument(`time`, `time|string`, []interface{}{arg}, 0))
	}
	if !t.Instance(tv) {
		panic(catch.Error(IllegalAssignment(t, tv)))
	}
	return tv
}

// Time returns the given timestamp as a dgo.Time
func Time(ts time.Time) dgo.Time {
	return &timeVal{ts}
}

// TimeFromString returns the given time string as a dgo.Time. The string must conform to
// the time.RFC3339 or time.RFC3339Nano format. The goFunc will panic if the given string
// cannot be parsed.
func TimeFromString(s string) dgo.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(catch.Error(err))
	}
	return &timeVal{ts}
}

func (v *timeVal) Equals(other interface{}) bool {
	switch ov := other.(type) {
	case *timeVal:
		return v.Equal(ov.Time)
	case time.Time:
		return v.Equal(ov)
	case *time.Time:
		return v.Equal(*ov)
	}
	return false
}

func (v *timeVal) SecondsWithFraction() float64 {
	y := v.Year()
	// Timestamps that represent a date before the year 1678 or after 2262 can
	// be represented as nanoseconds in an int64.
	if 1678 < y && y < 2262 {
		return float64(v.UnixNano()) / 1000000000.0
	}
	// Fall back to microsecond precision
	us := v.Unix()*1000000 + int64(v.Nanosecond())/1000
	return float64(us) / 1000000.0
}

func (v *timeVal) Float() dgo.Float {
	return floatVal(v.SecondsWithFraction())
}

func (v *timeVal) GoTime() *time.Time {
	return &v.Time
}

func (v *timeVal) HashCode() dgo.Hash {
	return dgo.Hash(v.UnixNano())
}

func (v *timeVal) Integer() dgo.Integer {
	return intVal(v.Unix())
}

func (v *timeVal) ReflectTo(value reflect.Value) {
	rv := reflect.ValueOf(&v.Time)
	k := value.Kind()
	if !(k == reflect.Ptr || k == reflect.Interface) {
		rv = rv.Elem()
	}
	value.Set(rv)
}

func (v *timeVal) String() string {
	return TypeString(v)
}

func (v *timeVal) ToBigFloat() *big.Float {
	return big.NewFloat(v.SecondsWithFraction())
}

func (v *timeVal) ToBigInt() *big.Int {
	return big.NewInt(v.Unix())
}

func (v *timeVal) ToFloat() (float64, bool) {
	return v.SecondsWithFraction(), true
}

func (v *timeVal) ToInt() (int64, bool) {
	return v.Unix(), true
}

func (v *timeVal) Type() dgo.Type {
	return v
}

// Time exact type implementation

func (v *timeVal) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v *timeVal) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v *timeVal) Generic() dgo.Type {
	return DefaultTimeType
}

func (v *timeVal) IsInstance(tv time.Time) bool {
	return v.Equal(tv)
}

func (v *timeVal) New(arg dgo.Value) dgo.Value {
	return newTime(v, arg)
}

func (v *timeVal) ReflectType() reflect.Type {
	return reflectTimeType
}

func (v *timeVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiTimeExact
}
