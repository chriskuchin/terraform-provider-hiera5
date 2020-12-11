package internal

import (
	"github.com/lyraproj/dgo/dgo"
)

type exactType struct {
	dgo.ExactType
}

type deepExactType struct {
	exactType
}

func (t *exactType) Assignable(other dgo.Type) bool {
	return t.Equals(other) || CheckAssignableTo(nil, other, t.ExactType)
}

func (t *exactType) Equals(other interface{}) bool {
	if ot, ok := other.(dgo.ExactType); ok && t.TypeIdentifier() == ot.TypeIdentifier() {
		return t.ExactValue().Equals(ot.ExactValue())
	}
	return false
}

func (t *exactType) HashCode() int {
	return t.ExactValue().HashCode()*7 + int(t.TypeIdentifier())
}

func (t *exactType) Instance(value interface{}) bool {
	return t.ExactValue().Equals(value)
}

func (t *exactType) String() string {
	return TypeString(t.ExactType)
}

func (t *exactType) Type() dgo.Type {
	return &metaType{t.ExactType}
}

func (t *deepExactType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *deepExactType) Equals(other interface{}) bool {
	return equals(nil, t.ExactType, other)
}

func (t *deepExactType) HashCode() int {
	return t.deepHashCode(nil)
}

func (t *deepExactType) Instance(value interface{}) bool {
	return Instance(nil, t.ExactType, value)
}

func (t *deepExactType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	return t.Equals(other) || CheckAssignableTo(guard, other, t.ExactType)
}

func (t *deepExactType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(dgo.ExactType); ok && t.TypeIdentifier() == ot.TypeIdentifier() {
		return equals(seen, t.ExactValue(), ot.ExactValue())
	}
	return false
}

func (t *deepExactType) deepHashCode(seen []dgo.Value) int {
	return deepHashCode(seen, t.ExactValue())*7 + int(t.TypeIdentifier())
}

func (t *deepExactType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	return t.ExactValue().Equals(value)
}
