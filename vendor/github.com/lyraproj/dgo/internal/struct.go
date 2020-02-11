package internal

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

type (
	structVal struct {
		rs     reflect.Value
		frozen bool
	}
)

func (v *structVal) AppendTo(w dgo.Indenter) {
	appendMapTo(v, w)
}

func (v *structVal) All(predicate dgo.EntryPredicate) bool {
	rv := v.rs
	rt := rv.Type()
	for i, n := 0, rt.NumField(); i < n; i++ {
		if !predicate(&mapEntry{&hstring{s: rt.Field(i).Name}, ValueFromReflected(rv.Field(i))}) {
			return false
		}
	}
	return true
}

func (v *structVal) AllKeys(predicate dgo.Predicate) bool {
	rv := v.rs
	rt := rv.Type()
	for i, n := 0, rt.NumField(); i < n; i++ {
		if !predicate(&hstring{s: rt.Field(i).Name}) {
			return false
		}
	}
	return true
}

func (v *structVal) AllValues(predicate dgo.Predicate) bool {
	rv := v.rs
	for i, n := 0, rv.NumField(); i < n; i++ {
		if !predicate(ValueFromReflected(rv.Field(i))) {
			return false
		}
	}
	return true
}

func (v *structVal) Any(predicate dgo.EntryPredicate) bool {
	return !v.All(func(entry dgo.MapEntry) bool { return !predicate(entry) })
}

func (v *structVal) AnyKey(predicate dgo.Predicate) bool {
	return !v.AllKeys(func(entry dgo.Value) bool { return !predicate(entry) })
}

func (v *structVal) AnyValue(predicate dgo.Predicate) bool {
	return !v.AllValues(func(entry dgo.Value) bool { return !predicate(entry) })
}

func (v *structVal) ContainsKey(key interface{}) bool {
	if s, ok := stringKey(key); ok {
		return v.rs.FieldByName(s).IsValid()
	}
	return false
}

func (v *structVal) Copy(frozen bool) dgo.Map {
	if frozen && v.frozen {
		return v
	}
	if frozen {
		return v.FrozenCopy().(dgo.Map)
	}
	// Thaw: Perform a by-value copy of the frozen struct
	rs := reflect.New(v.rs.Type()).Elem() // create and dereference pointer to a zero value
	rs.Set(v.rs)                          // copy v.rs to the zero value
	return &structVal{rs: rs, frozen: false}
}

func (v *structVal) Each(actor dgo.Consumer) {
	v.All(func(entry dgo.MapEntry) bool { actor(entry); return true })
}

func (v *structVal) EachEntry(actor dgo.EntryActor) {
	v.All(func(entry dgo.MapEntry) bool { actor(entry); return true })
}

func (v *structVal) EachKey(actor dgo.Consumer) {
	v.AllKeys(func(entry dgo.Value) bool { actor(entry); return true })
}

func (v *structVal) EachValue(actor dgo.Consumer) {
	v.AllValues(func(entry dgo.Value) bool { actor(entry); return true })
}

func (v *structVal) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *structVal) GoStruct() interface{} {
	return v.rs.Addr().Interface()
}

func (v *structVal) deepEqual(seen []dgo.Value, other deepEqual) bool {
	return mapEqual(seen, v, other)
}

func (v *structVal) HashCode() int {
	return deepHashCode(nil, v)
}

func (v *structVal) deepHashCode(seen []dgo.Value) int {
	hs := make([]int, v.Len())
	i := 0
	v.EachEntry(func(e dgo.MapEntry) {
		hs[i] = deepHashCode(seen, e)
		i++
	})
	sort.Ints(hs)
	h := 1
	for i = range hs {
		h = h*31 + hs[i]
	}
	return h
}

func (v *structVal) Freeze() {
	// Perform a shallow copy of the struct
	if !v.frozen {
		v.frozen = true
		v.rs = reflect.ValueOf(v.rs.Interface())
		v.EachValue(func(e dgo.Value) {
			if f, ok := e.(dgo.Freezable); ok {
				f.Freeze()
			}
		})
	}
}

func (v *structVal) Frozen() bool {
	return v.frozen
}

func (v *structVal) FrozenCopy() dgo.Value {
	if v.frozen {
		return v
	}

	// Perform a by-value copy of the struct
	rs := reflect.New(v.rs.Type()).Elem() // create and dereference pointer to a zero value
	rs.Set(v.rs)                          // copy v.rs to the zero value

	for i, n := 0, rs.NumField(); i < n; i++ {
		ef := rs.Field(i)
		ev := ValueFromReflected(ef)
		if f, ok := ev.(dgo.Freezable); ok && !f.Frozen() {
			ReflectTo(f.FrozenCopy(), ef)
		}
	}
	return &structVal{rs: rs, frozen: true}
}

func (v *structVal) Find(predicate dgo.EntryPredicate) dgo.MapEntry {
	rv := v.rs
	rt := rv.Type()
	for i, n := 0, rt.NumField(); i < n; i++ {
		e := &mapEntry{&hstring{s: rt.Field(i).Name}, ValueFromReflected(rv.Field(i))}
		if predicate(e) {
			return e
		}
	}
	return nil
}

func stringKey(key interface{}) (string, bool) {
	if hs, ok := key.(*hstring); ok {
		return hs.s, true
	}
	if s, ok := key.(string); ok {
		return s, true
	}
	return ``, false
}

func (v *structVal) Get(key interface{}) dgo.Value {
	if s, ok := stringKey(key); ok {
		rv := v.rs
		fv := rv.FieldByName(s)
		if fv.IsValid() {
			return ValueFromReflected(fv)
		}
	}
	return nil
}

func (v *structVal) Keys() dgo.Array {
	return arrayFromIterator(v.Len(), v.EachKey)
}

func (v *structVal) Len() int {
	return v.rs.NumField()
}

func (v *structVal) Map(mapper dgo.EntryMapper) dgo.Map {
	c := v.toHashMap()
	for e := c.first; e != nil; e = e.next {
		e.value = Value(mapper(e))
	}
	c.frozen = v.frozen
	return c
}

func (v *structVal) Merge(associations dgo.Map) dgo.Map {
	if associations.Len() == 0 || v == associations {
		return v
	}
	c := v.toHashMap()
	c.PutAll(associations)
	c.frozen = v.frozen && associations.Frozen()
	return c
}

func (v *structVal) Put(key, value interface{}) dgo.Value {
	if v.frozen {
		panic(frozenMap(`Put`))
	}
	if s, ok := stringKey(key); ok {
		rv := v.rs
		fv := rv.FieldByName(s)
		if fv.IsValid() {
			old := ValueFromReflected(fv)
			ReflectTo(Value(value), fv)
			return old
		}
	}
	panic(fmt.Errorf(`%s has no field named '%s'`, v.rs.Type(), key))
}

func (v *structVal) PutAll(associations dgo.Map) {
	associations.EachEntry(func(e dgo.MapEntry) { v.Put(e.Key(), e.Value()) })
}

func (v *structVal) ReflectTo(value reflect.Value) {
	if value.Kind() == reflect.Ptr {
		if v.frozen {
			// Don't expose pointer to frozen struct
			rs := reflect.New(v.rs.Type()) // create pointer to a zero value
			rs.Elem().Set(v.rs)            // copy v.rs to the zero value
			value.Set(rs)
		} else {
			value.Set(v.rs.Addr())
		}
	} else {
		// copy struct by value
		value.Set(v.rs)
	}
}

func (v *structVal) Remove(key interface{}) dgo.Value {
	panic(errors.New(`struct fields cannot be removed`))
}

func (v *structVal) RemoveAll(keys dgo.Array) {
	panic(errors.New(`struct fields cannot be removed`))
}

func (v *structVal) SetType(t interface{}) {
	panic(errors.New(`struct type is read only`))
}

func (v *structVal) String() string {
	return util.ToStringERP(v)
}

func (v *structVal) StringKeys() bool {
	return true
}

func (v *structVal) Type() dgo.Type {
	et := &exactMapType{value: v}
	et.ExactType = et
	return et
}

func (v *structVal) Values() dgo.Array {
	return arrayFromIterator(v.Len(), v.EachValue)
}

func (v *structVal) With(key, value interface{}) dgo.Map {
	c := v.toHashMap()
	c.Put(key, value)
	c.frozen = v.frozen
	return c
}

func (v *structVal) Without(key interface{}) dgo.Map {
	if v.Get(key) == nil {
		return v
	}
	c := v.toHashMap()
	c.Remove(key)
	c.frozen = v.frozen
	return c
}

func (v *structVal) WithoutAll(keys dgo.Array) dgo.Map {
	c := v.toHashMap()
	c.RemoveAll(keys)
	c.frozen = v.frozen
	return c
}

func (v *structVal) toHashMap() *hashMap {
	c := MapWithCapacity(v.Len(), nil)
	v.EachEntry(func(entry dgo.MapEntry) {
		c.Put(entry.Key(), entry.Value())
	})
	return c.(*hashMap)
}
