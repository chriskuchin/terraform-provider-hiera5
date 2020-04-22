package internal

import (
	"fmt"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

type (
	native reflect.Value

	nativeType struct {
		rt reflect.Type
	}
)

// DefaultNativeType is the unconstrained Native type
var DefaultNativeType = &nativeType{}

// Native creates the dgo representation of a reflect.Value.
func Native(rv reflect.Value) dgo.Native {
	nv := native(rv)
	return &nv
}

func (t *nativeType) Assignable(other dgo.Type) bool {
	if ot, ok := other.(*nativeType); ok {
		if t.rt == nil {
			return true
		}
		if ot.rt == nil {
			return false
		}
		return ot.rt.AssignableTo(t.rt)
	}
	return CheckAssignableTo(nil, other, t)
}

func (t *nativeType) Equals(other interface{}) bool {
	if ot, ok := other.(*nativeType); ok {
		return t.rt == ot.rt
	}
	return false
}

func (t *nativeType) HashCode() int {
	return util.StringHash(t.rt.Name())*31 + int(dgo.TiNative)
}

func (t *nativeType) ReflectType() reflect.Type {
	return t.rt
}

func (t *nativeType) Type() dgo.Type {
	return &metaType{t}
}

func (t *nativeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNative
}

func (t *nativeType) Instance(value interface{}) bool {
	if ov, ok := toReflected(value); ok {
		if t.rt == nil {
			return true
		}
		return ov.Type().AssignableTo(t.rt)
	}
	return false
}

func (t *nativeType) String() string {
	return TypeString(t)
}

func (t *nativeType) GoType() reflect.Type {
	return t.rt
}

func (v *native) Equals(other interface{}) bool {
	if b, ok := toReflected(other); ok {
		a := (*reflect.Value)(v)
		k := a.Kind()
		if k != b.Kind() {
			return false
		}
		return reflect.DeepEqual(a.Interface(), b.Interface())
	}
	return false
}

func (v *native) Freeze() {
	if !v.Frozen() {
		panic(fmt.Errorf(`native value cannot be frozen`))
	}
}

func (v *native) Frozen() bool {
	return false
}

func (v *native) FrozenCopy() dgo.Value {
	panic(fmt.Errorf(`native value cannot be frozen`))
}

func (v *native) ThawedCopy() dgo.Value {
	panic(fmt.Errorf(`native value cannot be copied`))
}

func (v *native) HashCode() int {
	rv := (*reflect.Value)(v)
	switch rv.Kind() {
	case reflect.Ptr:
		ev := rv.Elem()
		if ev.Kind() == reflect.Struct {
			return structHash(&ev)
		}
		p := rv.Pointer()
		return int(p ^ (p >> 32))
	case reflect.Chan, reflect.Uintptr:
		p := rv.Pointer()
		return int(p ^ (p >> 32))
	case reflect.Struct:
		return structHash(rv) * 3
	}
	return 1234
}

func structHash(rv *reflect.Value) int {
	n := rv.NumField()
	h := 1
	for i := 0; i < n; i++ {
		h = h*31 + ValueFromReflected(rv.Field(i)).HashCode()
	}
	return h
}

func (v *native) ReflectTo(value reflect.Value) {
	vr := (*reflect.Value)(v)
	if value.Kind() == reflect.Ptr {
		p := reflect.New(vr.Type())
		p.Elem().Set(*vr)
		value.Set(p)
	} else {
		value.Set(*vr)
	}
}

func (v *native) String() string {
	rv := (*reflect.Value)(v)
	if rv.CanInterface() {
		if s, ok := rv.Interface().(fmt.Stringer); ok {
			return s.String()
		}
	}
	return rv.String()
}

func (v *native) Type() dgo.Type {
	return &nativeType{(*reflect.Value)(v).Type()}
}

func (v *native) GoValue() interface{} {
	return (*reflect.Value)(v).Interface()
}

func toReflected(value interface{}) (reflect.Value, bool) {
	switch value := value.(type) {
	case *native:
		return reflect.Value(*value), true
	case dgo.Value:
		return reflect.Value{}, false
	}
	return reflect.ValueOf(value), true
}
