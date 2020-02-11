package internal

import (
	"fmt"
	"strconv"

	"github.com/lyraproj/dgo/dgo"
)

type (
	mapKeyError struct {
		mapType dgo.StructMapType
		key     dgo.Value
	}

	typeError struct {
		expected dgo.Type
		actual   dgo.Type
	}

	sizeError struct {
		sizedType     dgo.Type
		attemptedSize int
	}
)

func (v *mapKeyError) Equals(other interface{}) bool {
	if ov, ok := other.(*mapKeyError); ok {
		return v.mapType.Equals(ov.mapType) && v.key.Equals(ov.key)
	}
	return false
}

func (v *mapKeyError) Error() string {
	return fmt.Sprintf("key %s cannot added to type %s", v.key, TypeString(v.mapType))
}

func (v *mapKeyError) HashCode() int {
	return v.mapType.HashCode()*31 + v.key.HashCode()
}

func (v *mapKeyError) String() string {
	return v.Error()
}

func (v *mapKeyError) Type() dgo.Type {
	return DefaultErrorType
}

func (v *typeError) Equals(other interface{}) bool {
	if ov, ok := other.(*typeError); ok {
		return v.expected.Equals(ov.expected) && v.actual.Equals(ov.actual)
	}
	return false
}

func (v *typeError) Error() string {
	var what string
	switch actual := v.actual.(type) {
	case *exactStringType:
		what = fmt.Sprintf(`the string %s`, strconv.Quote(actual.value.s))
	case dgo.ExactType:
		what = fmt.Sprintf(`the value %s`, actual.ExactValue())
	default:
		what = fmt.Sprintf(`a value of type %s`, TypeString(actual))
	}
	return fmt.Sprintf("%s cannot be assigned to a variable of type %s", what, TypeString(v.expected))
}

func (v *typeError) HashCode() int {
	return v.expected.HashCode()*31 + v.actual.HashCode()
}

func (v *typeError) String() string {
	return v.Error()
}

func (v *typeError) Type() dgo.Type {
	return DefaultErrorType
}

func (v *sizeError) Equals(other interface{}) bool {
	if ov, ok := other.(*sizeError); ok {
		return v.sizedType.Equals(ov.sizedType) && v.attemptedSize == ov.attemptedSize
	}
	return false
}

func (v *sizeError) Error() string {
	return fmt.Sprintf(
		"size constraint violation on type %s when attempting resize to %d", TypeString(v.sizedType), v.attemptedSize)
}

func (v *sizeError) HashCode() int {
	return v.sizedType.HashCode()*7 + v.attemptedSize
}

func (v *sizeError) String() string {
	return v.Error()
}

func (v *sizeError) Type() dgo.Type {
	return DefaultErrorType
}

// IllegalAssignment returns the error that represents an assignment type constraint mismatch
func IllegalAssignment(t dgo.Type, v dgo.Value) dgo.Value {
	return &typeError{t, v.Type()}
}

// IllegalMapKey returns the error that represents an assignment map key constraint mismatch
func IllegalMapKey(t dgo.StructMapType, v dgo.Value) dgo.Value {
	return &mapKeyError{t, v.Type()}
}

// IllegalSize returns the error that represents an size constraint mismatch
func IllegalSize(t dgo.Type, sz int) dgo.Value {
	return &sizeError{t, sz}
}
