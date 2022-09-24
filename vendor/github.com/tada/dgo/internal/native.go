package internal

import (
	"fmt"
	"reflect"

	"github.com/tada/catch/pio"
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/util"
)

type (
	_rv = reflect.Value

	native struct {
		_rv
	}

	_rt = reflect.Type

	nativeType struct {
		_rt
	}
)

// DefaultNativeType is the unconstrained Native type
var DefaultNativeType = &nativeType{}

// Native creates the dgo representation of a reflect.Value.
func Native(rv reflect.Value) dgo.Native {
	return &native{rv}
}

func (t *nativeType) Assignable(other dgo.Type) bool {
	var ort reflect.Type
	switch ot := other.(type) {
	case *nativeType:
		ort = ot._rt
	case *native:
		ort = ot.ReflectType()
	default:
		return CheckAssignableTo(nil, other, t)
	}
	if t._rt == nil {
		return true
	}
	if ort == nil {
		return false
	}
	return ort.AssignableTo(t._rt)
}

func (t *nativeType) Equals(other interface{}) bool {
	if ot, ok := other.(*nativeType); ok {
		return t._rt == ot._rt
	}
	return false
}

func (t *nativeType) GoType() reflect.Type {
	return t._rt
}

func (t *nativeType) HashCode() dgo.Hash {
	h := dgo.Hash(dgo.TiNative)
	if t._rt != nil {
		h += util.StringHash(t.Name()) * 31
	}
	return h
}

func (t *nativeType) Instance(value interface{}) bool {
	if ov, ok := toReflected(value); ok {
		if t._rt == nil {
			return true
		}
		return ov.Type().AssignableTo(t._rt)
	}
	return false
}

func (t *nativeType) ReflectType() reflect.Type {
	return t._rt
}

func (t *nativeType) String() string {
	return TypeString(t)
}

func (t *nativeType) Type() dgo.Type {
	return MetaType(t)
}

func (t *nativeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNative
}

func (v *native) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v *native) Equals(other interface{}) bool {
	if b, ok := toReflected(other); ok {
		k := v.Kind()
		if k != b.Kind() {
			return false
		}
		return reflect.DeepEqual(v.Interface(), b.Interface())
	}
	return false
}

func (v *native) Format(s fmt.State, format rune) {
	if v.CanInterface() {
		doFormat(v.Interface(), s, format)
	} else {
		pio.WriteString(s, v.String())
	}
}

func (v *native) Generic() dgo.Type {
	return &nativeType{_rt: v.ReflectValue().Type()}
}

func (v *native) GoType() reflect.Type {
	return v.ReflectValue().Type()
}

func (v *native) GoValue() interface{} {
	return v.Interface()
}

func (v *native) HashCode() dgo.Hash {
	if v.CanAddr() {
		return dgo.Hash(v.UnsafeAddr())
	}
	switch v.Kind() {
	case reflect.Ptr:
		ev := v.Elem()
		if ev.Kind() == reflect.Struct {
			return structHash(&ev)
		}
		p := v.Pointer()
		return dgo.Hash(p ^ (p >> 32))
	case reflect.Struct:
		return structHash(v.ReflectValue()) * 3
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return dgo.Hash(v.Int())
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return dgo.Hash(v.Uint())
	case reflect.Float64, reflect.Float32:
		return dgo.Hash(v.Float())
	case reflect.Bool:
		if v.Bool() {
			return 1231
		}
		return 1237
	default:
		p := v.Pointer()
		return dgo.Hash(p ^ (p >> 32))
	}
}

func (v *native) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v *native) ReflectTo(value reflect.Value) {
	vr := v.ReflectValue()
	if value.Kind() == reflect.Ptr {
		p := reflect.New(vr.Type())
		p.Elem().Set(*vr)
		value.Set(p)
	} else {
		value.Set(*vr)
	}
}

func (v *native) ReflectType() reflect.Type {
	return v.ReflectValue().Type()
}

func (v *native) ReflectValue() *reflect.Value {
	return &v._rv
}

func (v *native) String() string {
	return TypeString(v)
}

func (v *native) Type() dgo.Type {
	return v
}

func (v *native) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNativeExact
}

func structHash(rv *reflect.Value) dgo.Hash {
	n := rv.NumField()
	h := dgo.Hash(1)
	for i := 0; i < n; i++ {
		h = h*31 + ValueFromReflected(rv.Field(i)).HashCode()
	}
	return h
}

func toReflected(value interface{}) (*reflect.Value, bool) {
	switch value := value.(type) {
	case *native:
		return value.ReflectValue(), true
	case dgo.Value:
		return nil, false
	}
	rv := reflect.ValueOf(value)
	return &rv, true
}
