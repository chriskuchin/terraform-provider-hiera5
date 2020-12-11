package internal

import (
	"fmt"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// booleanType represents an boolean without constraints (-1), constrained to false (0) or constrained to true(1)
	booleanType int

	exactBooleanType struct {
		exactType
		value boolean
	}

	boolean bool
)

// DefaultBooleanType is the unconstrained Boolean type
const DefaultBooleanType = booleanType(0)

// True is the dgo.Value for true
const True = boolean(true)

// False is the dgo.Value for false
const False = boolean(false)

// FalseType is the Boolean type that represents false
var FalseType dgo.BooleanType

// TrueType is the Boolean type that represents false
var TrueType dgo.BooleanType

func (t booleanType) Assignable(ot dgo.Type) bool {
	_, ok := ot.(dgo.BooleanType)
	return ok || CheckAssignableTo(nil, ot, t)
}

func (t booleanType) Equals(v interface{}) bool {
	return t == v
}

func (t booleanType) HashCode() int {
	return int(t.TypeIdentifier())
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
	return &metaType{t}
}

func (t booleanType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBoolean
}

func (t *exactBooleanType) ExactValue() dgo.Value {
	return t.value
}

func (t *exactBooleanType) Generic() dgo.Type {
	return DefaultBooleanType
}

func (t *exactBooleanType) IsInstance(value bool) bool {
	return bool(t.value) == value
}

func (t *exactBooleanType) New(arg dgo.Value) dgo.Value {
	return newBoolean(t, arg)
}

func (t *exactBooleanType) ReflectType() reflect.Type {
	return reflectBooleanType
}

func (t *exactBooleanType) TypeIdentifier() dgo.TypeIdentifier {
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

func (v boolean) GoBool() bool {
	return bool(v)
}

func (v boolean) HashCode() int {
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
	if v {
		return `true`
	}
	return `false`
}

func (v boolean) Type() dgo.Type {
	if v {
		return TrueType
	}
	return FalseType
}

func init() {
	et := &exactBooleanType{value: boolean(true)}
	et.ExactType = et
	TrueType = et

	et = &exactBooleanType{value: boolean(false)}
	et.ExactType = et
	FalseType = et
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
			panic(fmt.Errorf(`unable to create a bool from %s`, arg))
		}
	}
	if !t.Instance(v) {
		panic(IllegalAssignment(t, v))
	}
	return v
}
