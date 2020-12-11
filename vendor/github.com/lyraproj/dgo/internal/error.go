package internal

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

type (
	errw struct {
		error
	}

	errType int

	exactErrorType struct {
		exactType
		value *errw
	}
)

// DefaultErrorType is the unconstrained Error type
const DefaultErrorType = errType(0)

var reflectErrorType = reflect.TypeOf((*error)(nil)).Elem()

func (t errType) Type() dgo.Type {
	return &metaType{t}
}

func (t errType) Equals(other interface{}) bool {
	return t == other
}

func (t errType) HashCode() int {
	return int(t.TypeIdentifier())
}

func (t errType) Assignable(other dgo.Type) bool {
	_, ok := other.(errType)
	if !ok {
		_, ok = other.(*exactErrorType)
	}
	return ok || CheckAssignableTo(nil, other, t)
}

func (t errType) Instance(value interface{}) bool {
	_, ok := value.(error)
	return ok
}

func (t errType) IsInstance(err error) bool {
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

func (t *exactErrorType) Generic() dgo.Type {
	return DefaultErrorType
}

func (t *exactErrorType) IsInstance(err error) bool {
	return t.value.Equals(err)
}

func (t *exactErrorType) ReflectType() reflect.Type {
	return reflectErrorType
}

func (t *exactErrorType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiErrorExact
}

func (t *exactErrorType) ExactValue() dgo.Value {
	return t.value
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

func (e *errw) HashCode() int {
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
	return e.error.Error()
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
	ea := &exactErrorType{value: e}
	ea.ExactType = ea
	return ea
}
