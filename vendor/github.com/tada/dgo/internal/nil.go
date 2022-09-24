package internal

import (
	"fmt"
	"reflect"

	"github.com/tada/dgo/dgo"
)

type nilValue int

// Nil is the singleton dgo.Value for Nil
const Nil = nilValue(0)

func (nilValue) AppendTo(w dgo.Indenter) {
	w.Append(`nil`)
}

func (nilValue) CompareTo(other interface{}) (int, bool) {
	if Nil == other || nil == other {
		return 0, true
	}
	return -1, true
}

func (nilValue) HashCode() dgo.Hash {
	return 131
}

func (nilValue) Format(s fmt.State, format rune) {
	doFormat(nil, s, format)
}

func (nilValue) Equals(other interface{}) bool {
	return Nil == other || nil == other
}

func (nilValue) GoNil() interface{} {
	return nil
}

func (nilValue) ReflectTo(value reflect.Value) {
	value.Set(reflect.Zero(value.Type()))
}

func (nilValue) Type() dgo.Type {
	return Nil
}

func (nilValue) Assignable(ot dgo.Type) bool {
	_, ok := ot.(nilValue)
	return ok || CheckAssignableTo(nil, ot, Nil)
}

func (t nilValue) Instance(v interface{}) bool {
	return Nil == v || nil == v
}

func (t nilValue) ReflectType() reflect.Type {
	return reflectAnyType
}

func (t nilValue) String() string {
	return `nil`
}

func (t nilValue) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNil
}
