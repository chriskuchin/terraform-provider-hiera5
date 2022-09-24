package internal

import (
	"reflect"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/util"
)

type (
	errw struct {
		error
	}

	errType int
)

// DefaultErrorType is the unconstrained Error type
const DefaultErrorType = errType(0)

var reflectErrorType = reflect.TypeOf((*error)(nil)).Elem()

func (t errType) Type() dgo.Type {
	return MetaType(t)
}

func (t errType) Equals(other interface{}) bool {
	return t == other
}

func (t errType) HashCode() dgo.Hash {
	return dgo.Hash(t.TypeIdentifier())
}

func (t errType) Assignable(other dgo.Type) bool {
	_, ok := other.(errType)
	if !ok {
		_, ok = other.(*errw)
	}
	return ok || CheckAssignableTo(nil, other, t)
}

func (t errType) Instance(value interface{}) bool {
	_, ok := value.(error)
	return ok
}

func (t errType) IsInstance(_ error) bool {
	return true
}

func (t errType) ReflectType() reflect.Type {
	return reflectErrorType
}

func (t errType) String() string {
	return TypeString(t)
}

func (t errType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiError
}

func (e *errw) Assignable(other dgo.Type) bool {
	return e.Equals(other) || CheckAssignableTo(nil, other, e)
}

func (e *errw) Generic() dgo.Type {
	return DefaultErrorType
}

func (e *errw) Instance(value interface{}) bool {
	return e.Equals(value)
}

func (e *errw) IsInstance(err error) bool {
	return e.Equals(err)
}

func (e *errw) ReflectType() reflect.Type {
	return reflectErrorType
}

func (e *errw) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiErrorExact
}

func (e *errw) Equals(other interface{}) bool {
	if oe, ok := other.(*errw); ok {
		return e.error.Error() == oe.error.Error()
	}
	if oe, ok := other.(error); ok {
		return e.error.Error() == oe.Error()
	}
	return false
}

func (e *errw) HashCode() dgo.Hash {
	return util.StringHash(e.error.Error())
}

func (e *errw) Error() string {
	return e.error.Error()
}

func (e *errw) ReflectTo(value reflect.Value) {
	if value.Kind() == reflect.Ptr {
		value.Set(reflect.ValueOf(&e.error))
	} else {
		value.Set(reflect.ValueOf(e.error))
	}
}

func (e *errw) String() string {
	return TypeString(e)
}

func (e *errw) Unwrap() error {
	if u, ok := e.error.(interface {
		Unwrap() error
	}); ok {
		return u.Unwrap()
	}
	return nil
}

func (e *errw) Type() dgo.Type {
	return e
}
