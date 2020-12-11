package internal

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

// anyType represents all possible values
type anyType int

// DefaultAnyType is the unconstrained Any type
const DefaultAnyType = anyType(0)

var reflectAnyType = reflect.TypeOf((*interface{})(nil)).Elem()

func (t anyType) Assignable(other dgo.Type) bool {
	return true
}

func (t anyType) Equals(other interface{}) bool {
	return t == other
}

func (t anyType) HashCode() int {
	return int(dgo.TiAny)
}

func (t anyType) Instance(value interface{}) bool {
	return true
}

// ReflectType returns the reflect.Type for the given dgo.Type
func (t anyType) ReflectType() reflect.Type {
	return reflectAnyType
}

func (t anyType) String() string {
	return TypeString(t)
}

func (t anyType) Type() dgo.Type {
	return &metaType{t}
}

func (t anyType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiAny
}
