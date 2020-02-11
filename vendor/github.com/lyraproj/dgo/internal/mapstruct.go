package internal

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// structType describes each mapEntry of a map
	structType struct {
		additional bool
		keys       array
		values     array
		required   []bool
	}

	structEntry struct {
		mapEntry
		required bool
	}
)

// StructMapTypeUnresolved returns an unresolved new StructMapType type built from the given StructMapEntries. The
// fact that it is unresolved vouches for that it may have keys that are not yet exact types but might become exact
// once they are resolved.
func StructMapTypeUnresolved(additional bool, entries []dgo.StructMapEntry) dgo.StructMapType {
	l := len(entries)
	exact := !additional
	keys := make([]dgo.Value, l)
	values := make([]dgo.Value, l)
	required := make([]bool, l)
	for i := 0; i < l; i++ {
		e := entries[i]
		kt := e.Key().(dgo.Type)
		vt := e.Value().(dgo.Type)
		if exact && !(e.Required() && dgo.IsExact(kt) && dgo.IsExact(vt)) {
			exact = false
		}
		keys[i] = kt
		values[i] = vt
		required[i] = e.Required()
	}

	if exact {
		return createExactMap(keys, values)
	}

	return &structType{
		additional: additional,
		keys:       array{slice: keys, frozen: true},
		values:     array{slice: values, frozen: true},
		required:   required}
}

func createExactMap(keys, values []dgo.Value) dgo.StructMapType {
	l := len(keys)
	m := MapWithCapacity(l, nil)
	for i := 0; i < l; i++ {
		m.Put(keys[i].(dgo.ExactType).ExactValue(), values[i].(dgo.ExactType).ExactValue())
	}
	return m.Type().(dgo.StructMapType)
}

// StructMapType returns a new StructMapType type built from the given StructMapEntries.
func StructMapType(additional bool, entries []dgo.StructMapEntry) dgo.StructMapType {
	t := StructMapTypeUnresolved(additional, entries)
	if st, ok := t.(*structType); ok {
		st.checkExactKeys()
	}
	return t
}

var sfmType dgo.MapType

// StructFromMapType returns the map type used when validating the map sent to
// StructMapTypeFromMap
func StructFromMapType() dgo.MapType {
	if sfmType == nil {
		sfmType = Parse(`map[string](dgo|type|{type:dgo|type,required?:bool,...})`).(dgo.MapType)
	}
	return sfmType
}

// StructMapTypeFromMap returns a new type built from a map[string](dgo|type|{type:dgo|type,required?:bool,...})
func StructMapTypeFromMap(additional bool, entries dgo.Map) dgo.StructMapType {
	if !StructFromMapType().Instance(entries) {
		panic(IllegalAssignment(sfmType, entries))
	}
	l := entries.Len()
	keys := make([]dgo.Value, l)
	values := make([]dgo.Value, l)
	required := make([]bool, l)
	i := 0

	// turn dgo|type into type
	asType := func(v dgo.Value) dgo.Type {
		tp, ok := v.(dgo.Type)
		if !ok {
			var s dgo.String
			if s, ok = v.(dgo.String); ok {
				v = Parse(s.GoString())
				tp, ok = v.(dgo.Type)
			}
			if !ok {
				tp = v.Type()
			}
		}
		return tp
	}

	exact := !additional
	entries.EachEntry(func(e dgo.MapEntry) {
		rq := true
		kt := e.Key().Type()
		var vt dgo.Type
		if vm, ok := e.Value().(dgo.Map); ok {
			vt = asType(vm.Get(`type`))
			if rqv := vm.Get(`required`); rqv != nil {
				rq = rqv.(dgo.Boolean).GoBool()
			}
		} else {
			vt = asType(e.Value())
		}
		if exact && !(rq && dgo.IsExact(kt) && dgo.IsExact(vt)) {
			exact = false
		}
		keys[i] = kt
		values[i] = vt
		required[i] = rq
		i++
	})

	if exact {
		return createExactMap(keys, values)
	}

	t := &structType{
		additional: additional,
		keys:       array{slice: keys, frozen: true},
		values:     array{slice: values, frozen: true},
		required:   required}

	t.checkExactKeys()
	return t
}

func (t *structType) checkExactKeys() {
	ks := t.keys.slice
	for i := range ks {
		if !dgo.IsExact(ks[i].(dgo.Type)) {
			panic(`non exact key types is not yet supported`)
		}
	}
}

func (t *structType) Additional() bool {
	return t.additional
}

func (t *structType) CheckEntry(k, v dgo.Value) dgo.Value {
	ks := t.keys.slice
	for i := range ks {
		kt := ks[i].(dgo.Type)
		if kt.Instance(k) {
			vt := t.values.slice[i].(dgo.Type)
			if vt.Instance(v) {
				return nil
			}
			return IllegalAssignment(vt, v)
		}
	}
	if t.additional {
		return nil
	}
	return IllegalMapKey(t, k)
}

func (t *structType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *structType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	switch ot := other.(type) {
	case *structType:
		mrs := t.required
		mks := t.keys.slice
		mvs := t.values.slice
		ors := ot.required
		oks := ot.keys.slice
		ovs := ot.values.slice
		oc := 0

	nextKey:
		for mi := range mks {
			rq := mrs[mi]
			mk := mks[mi]
			for oi := range oks {
				ok := oks[oi]
				if mk.Equals(ok) {
					if rq && !ors[oi] {
						return false
					}
					if !Assignable(guard, mvs[mi].(dgo.Type), ovs[oi].(dgo.Type)) {
						return false
					}
					oc++
					continue nextKey
				}
			}
			if rq || ot.additional { // additional included since key is allowed with unconstrained value
				return false
			}
		}
		return t.additional || oc == len(oks)
	case *exactMapType:
		ov := ot.value
		return Instance(guard, t, ov)
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *structType) Each(actor func(dgo.StructMapEntry)) {
	ks := t.keys.slice
	vs := t.values.slice
	rs := t.required
	for i := range ks {
		actor(&structEntry{mapEntry: mapEntry{key: ks[i], value: vs[i]}, required: rs[i]})
	}
}

func (t *structType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *structType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*structType); ok {
		return t.additional == ot.additional &&
			boolsEqual(t.required, ot.required) &&
			equals(seen, &t.keys, &ot.keys) &&
			equals(seen, &t.values, &ot.values)
	}
	return false
}

func (t *structType) Generic() dgo.Type {
	return newMapType(Generic(t.KeyType()), Generic(t.ValueType()), 0, math.MaxInt64)
}

func (t *structType) HashCode() int {
	return deepHashCode(nil, t)
}

func (t *structType) deepHashCode(seen []dgo.Value) int {
	h := boolsHash(t.required)*31 + deepHashCode(seen, &t.keys)*31 + deepHashCode(seen, &t.values)
	if t.additional {
		h *= 3
	}
	return h
}

func (t *structType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *structType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	if om, ok := value.(dgo.Map); ok {
		ks := t.keys.slice
		vs := t.values.slice
		rs := t.required
		oc := 0
		for i := range ks {
			k := ks[i].(dgo.ExactType)
			if ov := om.Get(k.ExactValue()); ov != nil {
				oc++
				if !Instance(guard, vs[i].(dgo.Type), ov) {
					return false
				}
			} else if rs[i] {
				return false
			}
		}
		return t.additional || oc == om.Len()
	}
	return false
}

func (t *structType) Get(key interface{}) dgo.StructMapEntry {
	kv := Value(key)
	if _, ok := kv.(dgo.Type); !ok {
		kv = kv.Type()
	}
	i := t.keys.IndexOf(kv)
	if i >= 0 {
		return StructMapEntry(kv, t.values.slice[i], t.required[i])
	}
	return nil
}

func (t *structType) KeyType() dgo.Type {
	switch t.keys.Len() {
	case 0:
		return DefaultAnyType
	case 1:
		return t.keys.Get(0).(dgo.Type)
	default:
		return (*allOfType)(&t.keys)
	}
}

func (t *structType) Len() int {
	return len(t.required)
}

func (t *structType) Max() int {
	m := len(t.required)
	if m == 0 || t.additional {
		return math.MaxInt64
	}
	return m
}

func (t *structType) Min() int {
	min := 0
	rs := t.required
	for i := range rs {
		if rs[i] {
			min++
		}
	}
	return min
}

func (t *structType) New(arg dgo.Value) dgo.Value {
	return newMap(t, arg)
}

func (t *structType) ReflectType() reflect.Type {
	return reflect.MapOf(t.KeyType().ReflectType(), t.ValueType().ReflectType())
}

func (t *structType) Resolve(ap dgo.AliasMap) {
	ks := t.keys.slice
	vs := t.values.slice
	t.keys.slice = []dgo.Value{}
	t.values.slice = []dgo.Value{}
	for i := range ks {
		ks[i] = ap.Replace(ks[i].(dgo.Type))
		vs[i] = ap.Replace(vs[i].(dgo.Type))
	}
	t.keys.slice = ks
	t.values.slice = vs
	t.checkExactKeys()
}

func (t *structType) String() string {
	return TypeString(t)
}

func (t *structType) Type() dgo.Type {
	return &metaType{t}
}

func (t *structType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiStruct
}

func (t *structType) Unbounded() bool {
	return t.additional && t.Min() == 0
}

func parameterLabel(key dgo.Value) string {
	return fmt.Sprintf(`parameter '%s'`, key)
}

func (t *structType) Validate(keyLabel func(key dgo.Value) string, value interface{}) []error {
	return validate(t, keyLabel, value)
}

func (t *structType) ValidateVerbose(value interface{}, out dgo.Indenter) bool {
	return validateVerbose(t, value, out)
}

func validate(t dgo.StructMapType, keyLabel func(key dgo.Value) string, value interface{}) []error {
	var errs []error
	pm, ok := Value(value).(dgo.Map)
	if !ok {
		return []error{errors.New(`value is not a Map`)}
	}

	if keyLabel == nil {
		keyLabel = parameterLabel
	}
	t.Each(func(e dgo.StructMapEntry) {
		ek := e.Key().(dgo.ExactType).ExactValue()
		if v := pm.Get(ek); v != nil {
			ev := e.Value().(dgo.Type)
			if !ev.Instance(v) {
				errs = append(errs, fmt.Errorf(`%s is not an instance of type %s`, keyLabel(ek), ev))
			}
		} else if e.Required() {
			errs = append(errs, fmt.Errorf(`missing required %s`, keyLabel(ek)))
		}
	})
	pm.EachKey(func(k dgo.Value) {
		if t.Get(k) == nil {
			errs = append(errs, fmt.Errorf(`unknown %s`, keyLabel(k)))
		}
	})
	return errs
}

func validateVerbose(t dgo.StructMapType, value interface{}, out dgo.Indenter) bool {
	pm, ok := Value(value).(dgo.Map)
	if !ok {
		out.Append(`value is not a Map`)
		return false
	}

	inner := out.Indent()
	t.Each(func(e dgo.StructMapEntry) {
		ek := e.Key().(dgo.ExactType).ExactValue()
		ev := e.Value().(dgo.Type)
		out.Printf(`Validating '%s' against definition %s`, ek, ev)
		inner.NewLine()
		inner.Printf(`'%s' `, ek)
		if v := pm.Get(ek); v != nil {
			if ev.Instance(v) {
				inner.Append(`OK!`)
			} else {
				ok = false
				inner.Append(`FAILED!`)
				inner.NewLine()
				inner.Printf(`Reason: expected a value of type %s, got %s`, ev, v.Type())
			}
		} else if e.Required() {
			ok = false
			inner.Append(`FAILED!`)
			inner.NewLine()
			inner.Append(`Reason: required key not found in input`)
		}
		out.NewLine()
	})
	pm.EachKey(func(k dgo.Value) {
		if t.Get(k) == nil {
			ok = false
			out.Printf(`Validating '%s'`, k)
			inner.NewLine()
			inner.Printf(`'%s' FAILED!`, k)
			inner.NewLine()
			inner.Append(`Reason: key is not found in definition`)
			out.NewLine()
		}
	})
	return ok
}

func (t *structType) ValueType() dgo.Type {
	switch t.values.Len() {
	case 0:
		return DefaultAnyType
	case 1:
		return t.values.Get(0).(dgo.Type)
	default:
		return (*allOfType)(&t.values)
	}
}

// StructMapEntry returns a new StructMapEntry initiated with the given parameters
func StructMapEntry(key interface{}, value interface{}, required bool) dgo.StructMapEntry {
	kv := Value(key)
	if _, ok := kv.(dgo.Type); !ok {
		kv = kv.Type()
	}
	vv := Value(value)
	if _, ok := vv.(dgo.Type); !ok {
		vv = vv.Type()
	}
	return &structEntry{mapEntry: mapEntry{key: kv, value: vv}, required: required}
}

func (t *structEntry) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *structEntry) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(dgo.StructMapEntry); ok {
		return t.required == ot.Required() &&
			equals(seen, t.mapEntry.key, ot.Key()) &&
			equals(seen, t.mapEntry.value, ot.Value())
	}
	return false
}

func (t *structEntry) Required() bool {
	return t.required
}

func boolsHash(s []bool) int {
	h := 1
	for i := range s {
		m := 2
		if s[i] {
			m = 3
		}
		h *= m
	}
	return h
}

func boolsEqual(a, b []bool) bool {
	l := len(a)
	if l != len(b) {
		return false
	}
	for l--; l >= 0; l-- {
		if a[l] != b[l] {
			return false
		}
	}
	return true
}
