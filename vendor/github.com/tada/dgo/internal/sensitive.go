package internal

import (
	"reflect"

	"github.com/tada/dgo/dgo"
)

type (
	sensitive struct {
		value dgo.Value
	}

	sensitiveType struct {
		wrapped dgo.Type
	}
)

// DefaultSensitiveType is the unconstrained Sensitive type
var DefaultSensitiveType = &sensitiveType{wrapped: DefaultAnyType}

// SensitiveType returns a Sensitive dgo.Type that wraps the given dgo.Type
func SensitiveType(args []interface{}) dgo.Type {
	switch len(args) {
	case 0:
		return DefaultSensitiveType
	case 1:
		return &sensitiveType{wrapped: AsType(Value(args[0]))}
	}
	panic(illegalArgumentCount(`SensitiveType`, 0, 1, len(args)))
}

func (t *sensitiveType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *sensitiveType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	if ot, ok := other.(*sensitiveType); ok {
		return Assignable(guard, t.wrapped, ot.wrapped)
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *sensitiveType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *sensitiveType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*sensitiveType); ok {
		return equals(seen, t.wrapped, ot.wrapped)
	}
	return false
}

func (t *sensitiveType) HashCode() dgo.Hash {
	return deepHashCode(nil, t)
}

func (t *sensitiveType) deepHashCode(seen []dgo.Value) dgo.Hash {
	return dgo.Hash(dgo.TiSensitive)*31 + deepHashCode(seen, t.wrapped)
}

func (t *sensitiveType) Instance(value interface{}) bool {
	if ov, ok := value.(*sensitive); ok {
		return t.wrapped.Instance(ov.value)
	}
	return false
}

var reflectSensitiveType = reflect.TypeOf((*dgo.Sensitive)(nil)).Elem()

func (t *sensitiveType) ReflectType() reflect.Type {
	return reflectSensitiveType
}

func (t *sensitiveType) Operand() dgo.Type {
	return t.wrapped
}

func (t *sensitiveType) Operator() dgo.TypeOp {
	return dgo.OpSensitive
}

func (t *sensitiveType) New(arg dgo.Value) dgo.Value {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`sensitive`, 1, 1)
		arg = args.Get(0)
	}
	if s, ok := arg.(dgo.Sensitive); ok {
		return s
	}
	return Sensitive(arg)
}

func (t *sensitiveType) String() string {
	return TypeString(t)
}

func (t *sensitiveType) Type() dgo.Type {
	return MetaType(t)
}

func (t *sensitiveType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiSensitive
}

// Sensitive creates a new Sensitive that wraps the given value
func Sensitive(v interface{}) dgo.Sensitive {
	return &sensitive{Value(v)}
}

func (v *sensitive) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *sensitive) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ov, ok := other.(*sensitive); ok {
		return equals(seen, v.value, ov.value)
	}
	return false
}

func (v *sensitive) Frozen() bool {
	if f, ok := v.value.(dgo.Mutability); ok {
		return f.Frozen()
	}
	return true
}

func (v *sensitive) FrozenCopy() dgo.Value {
	if f, ok := v.value.(dgo.Mutability); ok && !f.Frozen() {
		return &sensitive{f.FrozenCopy()}
	}
	return v
}

func (v *sensitive) ThawedCopy() dgo.Value {
	if f, ok := v.value.(dgo.Mutability); ok {
		return &sensitive{f.ThawedCopy()}
	}
	return v
}

func (v *sensitive) HashCode() dgo.Hash {
	return deepHashCode(nil, v)
}

func (v *sensitive) deepHashCode(seen []dgo.Value) dgo.Hash {
	return deepHashCode(seen, v.value) * 7
}

func (v *sensitive) String() string {
	return `sensitive [value redacted]`
}

func (v *sensitive) Type() dgo.Type {
	return &sensitiveType{wrapped: Generic(v.value.Type())}
}

func (v *sensitive) Unwrap() dgo.Value {
	return v.value
}
