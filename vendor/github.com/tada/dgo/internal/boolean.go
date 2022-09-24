package internal

import (
	"fmt"
	"reflect"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
)

type (
	// booleanType represents an boolean without constraints (-1), constrained to false (0) or constrained to true(1)
	booleanType int

	boolean bool
)

// DefaultBooleanType is the unconstrained Boolean type
const DefaultBooleanType = booleanType(0)

// True is the dgo.Value for true
const True = boolean(true)

// False is the dgo.Value for false
const False = boolean(false)

func (t booleanType) Assignable(ot dgo.Type) bool {
	_, ok := ot.(dgo.BooleanType)
	return ok || CheckAssignableTo(nil, ot, t)
}

func (t booleanType) Equals(v interface{}) bool {
	return t == v
}

func (t booleanType) HashCode() dgo.Hash {
	return dgo.Hash(t.TypeIdentifier())
}

func (t booleanType) Instance(v interface{}) bool {
	switch v.(type) {
	case boolean, bool:
		return true
	default:
		return false
	}
}

func (t booleanType) IsInstance(v bool) bool {
	return true
}

func (t booleanType) New(arg dgo.Value) dgo.Value {
	return newBoolean(t, arg)
}

var reflectBooleanType = reflect.TypeOf(true)

func (t booleanType) ReflectType() reflect.Type {
	return reflectBooleanType
}

func (t booleanType) String() string {
	return TypeString(t)
}

func (t booleanType) Type() dgo.Type {
	return MetaType(t)
}

func (t booleanType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBoolean
}

func (v boolean) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v boolean) Generic() dgo.Type {
	return DefaultBooleanType
}

func (v boolean) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v boolean) IsInstance(value bool) bool {
	return bool(v) == value
}

func (v boolean) New(arg dgo.Value) dgo.Value {
	return newBoolean(v, arg)
}

func (v boolean) ReflectType() reflect.Type {
	return reflectBooleanType
}

func (v boolean) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBooleanExact
}

func (v boolean) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	switch ov := other.(type) {
	case boolean:
		r = 0
		if v {
			if !ov {
				r = 1
			}
		} else if ov {
			r = -1
		}
	case nilValue:
		r = 1
	default:
		ok = false
	}
	return
}

func (v boolean) Equals(other interface{}) bool {
	if ov, ok := other.(boolean); ok {
		return v == ov
	}
	if ov, ok := other.(bool); ok {
		return bool(v) == ov
	}
	return false
}

func (v boolean) Format(s fmt.State, format rune) {
	doFormat(bool(v), s, format)
}

func (v boolean) GoBool() bool {
	return bool(v)
}

func (v boolean) HashCode() dgo.Hash {
	if v {
		return 1231
	}
	return 1237
}

func (v boolean) ReflectTo(value reflect.Value) {
	b := bool(v)
	switch value.Kind() {
	case reflect.Interface:
		value.Set(reflect.ValueOf(b))
	case reflect.Ptr:
		value.Set(reflect.ValueOf(&b))
	default:
		value.SetBool(b)
	}
}

func (v boolean) String() string {
	return TypeString(v)
}

func (v boolean) Type() dgo.Type {
	return v
}

var boolStringType = CiEnumType([]string{`true`, `false`, `yes`, `no`, `y`, `n`})
var trueStringType = CiEnumType([]string{`true`, `yes`, `y`})

func newBoolean(t dgo.Type, arg dgo.Value) dgo.Value {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`bool`, 1, 1)
		arg = args.Get(0)
	}
	var v dgo.Value
	switch arg := arg.(type) {
	case boolean:
		v = arg
	case intVal:
		v = boolean(arg != 0)
	case floatVal:
		v = boolean(arg != 0)
	default:
		if boolStringType.Instance(arg) {
			v = boolean(trueStringType.Instance(arg))
		} else {
			panic(catch.Error(`unable to create a bool from %s`, arg))
		}
	}
	if !t.Instance(v) {
		panic(catch.Error(IllegalAssignment(t, v)))
	}
	return v
}
