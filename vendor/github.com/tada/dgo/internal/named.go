package internal

import (
	"reflect"
	"sync"

	"github.com/tada/catch"
	"github.com/tada/catch/pio"
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/util"
)

type (
	named struct {
		name       string
		ctor       dgo.Constructor
		extractor  dgo.InitArgExtractor
		implType   reflect.Type
		ifdType    reflect.Type
		asgChecker dgo.AssignableChecker
	}

	parameterized struct {
		dgo.NamedType
		params dgo.Array
	}

	exactNamed struct {
		dgo.NamedTypeExtension
		value dgo.Value
	}
)

var namedTypes = sync.Map{}

func defaultAsgChecker(t dgo.NamedType, other dgo.Type) bool {
	if ot, ok := other.(dgo.NamedType); ok {
		return ot.ReflectType().AssignableTo(t.AssignableType())
	}
	return false
}

// RemoveNamedType removes a named type from the global type registry. It is primarily intended for
// testing purposes.
func RemoveNamedType(name string) {
	namedTypes.Delete(name)
}

// NewNamedType registers a new named and optionally parameterized type under the given name with the global type registry.
// The method panics if a type has already been registered with the same name.
//
// name: name of the type
//
// ctor: optional constructor that creates new values of this type
//
// extractor: optional extractor of the value used when serializing/deserializing this type
//
// implType optional reflected zero value type of implementation
//
// ifdType optional reflected nil value of interface type
//
// asgChecker optional function to check what other types that are assignable to this type
func NewNamedType(
	name string,
	ctor dgo.Constructor,
	extractor dgo.InitArgExtractor,
	implType,
	ifdType reflect.Type,
	asgChecker dgo.AssignableChecker) dgo.NamedType {
	t, loaded := namedTypes.LoadOrStore(name,
		&named{name: name, ctor: ctor, extractor: extractor, implType: implType, ifdType: ifdType, asgChecker: asgChecker})
	if loaded {
		panic(catch.Error(`attempt to redefine named type '%s'`, name))
	}
	return t.(*named)
}

// ExactNamedType returns the exact NamedType that represents the given value.
func ExactNamedType(namedType dgo.NamedType, value dgo.Value) dgo.NamedType {
	return &exactNamed{NamedTypeExtension: namedType, value: value}
}

// NamedType returns the type with the given name from the global type registry. The function returns
// nil if no type has been registered under the given name.
func NamedType(name string) dgo.NamedType {
	if t, ok := namedTypes.Load(name); ok {
		return t.(*named)
	}
	return nil
}

// ParameterizedType returns the named type amended with the given parameters.
func ParameterizedType(named dgo.NamedType, params dgo.Array) dgo.NamedType {
	return &parameterized{named, params}
}

// NamedTypeFromReflected returns the named type for the reflected implementation type from the global type
// registry. The function returns nil if no such type has been registered.
func NamedTypeFromReflected(rt reflect.Type) dgo.NamedType {
	var t dgo.NamedType
	namedTypes.Range(func(k, v interface{}) bool {
		nt := v.(*named)
		if rt == nt.implType {
			t = nt
			return false
		}
		return true
	})
	return t
}

func (t *named) AssignableType() reflect.Type {
	if t.ifdType != nil {
		return t.ifdType
	}
	return t.implType
}

func (t *named) Assignable(other dgo.Type) bool {
	f := t.asgChecker
	if f == nil {
		f = defaultAsgChecker
	}
	return f(t, other) || CheckAssignableTo(nil, other, t)
}

func (t *named) New(arg dgo.Value) dgo.Value {
	return newNamed(t, arg)
}

func (t *named) Equals(other interface{}) bool {
	if ot, ok := other.(*named); ok {
		return t.name == ot.name
	}
	return false
}

func (t *named) HashCode() dgo.Hash {
	return util.StringHash(t.name)*7 + dgo.Hash(dgo.TiNamed)
}

func (t *named) ExtractInitArg(value dgo.Value) dgo.Value {
	if t.extractor != nil {
		return t.extractor(value)
	}
	panic(catch.Error(`creating new instances of %s is not possible`, t.name))
}

func (t *named) Instance(value interface{}) bool {
	return t.Assignable(Value(value).Type())
}

func (t *named) Name() string {
	return t.name
}

func (t *named) Parameters() dgo.Array {
	return nil
}

func (t *named) ReflectType() reflect.Type {
	return t.implType
}

func (t *named) String() string {
	return TypeString(t)
}

func (t *named) Type() dgo.Type {
	return MetaType(t)
}

func (t *named) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNamed
}

func (t *named) ValueString(v dgo.Value) string {
	b := util.NewERPIndenter(``)
	pio.WriteString(b, t.name)
	if t.extractor != nil {
		switch t.extractor(v).(type) {
		case dgo.Array, dgo.Map, dgo.String:
		default:
			b.AppendRune(' ')
		}
		b.AppendValue(t.extractor(v))
	}
	return b.String()
}

func (t *parameterized) Assignable(other dgo.Type) bool {
	f := t.NamedType.(*named).asgChecker
	if f == nil {
		f = defaultAsgChecker
	}
	return f(t, other) || CheckAssignableTo(nil, other, t)
}

func (t *parameterized) Equals(other interface{}) bool {
	if ot, ok := other.(*parameterized); ok {
		return t.NamedType.Equals(ot.NamedType) && t.params.Equals(ot.params)
	}
	return false
}

func (t *parameterized) Generic() dgo.Type {
	return t.NamedType
}

func (t *parameterized) HashCode() dgo.Hash {
	return t.NamedType.HashCode()*31 + t.params.HashCode()
}

func (t *parameterized) Instance(value interface{}) bool {
	return t.Assignable(Value(value).Type())
}

func (t *parameterized) New(arg dgo.Value) dgo.Value {
	return newNamed(t, arg)
}

func (t *parameterized) Parameters() dgo.Array {
	return t.params
}

func (t *parameterized) String() string {
	return TypeString(t)
}

func (t *parameterized) Type() dgo.Type {
	return MetaType(t)
}

func (t *exactNamed) Assignable(other dgo.Type) bool {
	return t.Equals(other) || CheckAssignableTo(nil, other, t)
}

func (t *exactNamed) ExactValue() dgo.Value {
	return t.value
}

func (t *exactNamed) Equals(other interface{}) bool {
	if ot, ok := other.(dgo.ExactType); ok && t.TypeIdentifier() == ot.TypeIdentifier() {
		return t.ExactValue().Equals(ot.ExactValue())
	}
	return false
}

func (t *exactNamed) Generic() dgo.Type {
	return t.NamedTypeExtension.(dgo.Type)
}

func (t *exactNamed) HashCode() dgo.Hash {
	return t.ExactValue().HashCode()*7 + dgo.Hash(t.TypeIdentifier())
}

func (t *exactNamed) Instance(value interface{}) bool {
	return t.ExactValue().Equals(value)
}

func (t *exactNamed) New(arg dgo.Value) dgo.Value {
	return newNamed(t, arg)
}

func (t *exactNamed) ReflectType() reflect.Type {
	return t.Generic().ReflectType()
}

func (t *exactNamed) String() string {
	return TypeString(t)
}

func (t *exactNamed) Type() dgo.Type {
	return MetaType(t)
}

func (t *exactNamed) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNamedExact
}

func newNamed(rt dgo.NamedType, arg dgo.Value) dgo.Value {
	t := Generic(rt).(*named)
	if t.ctor == nil {
		panic(catch.Error(`creating new instances of %s is not possible`, rt.Name()))
	}
	v := t.ctor(arg)
	if t.asgChecker == nil || t.asgChecker(rt, v.Type()) {
		return v
	}
	panic(catch.Error(IllegalAssignment(rt, v)))
}
