package internal

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/util"
)

type (
	arraySlice interface {
		_slice() []dgo.Value
		_deepContainsAll(seen []dgo.Value, other dgo.Iterable) bool
	}

	array struct {
		slice []dgo.Value
	}

	arrayFrozen struct {
		array
	}

	// defaultArrayType is the unconstrained array type
	defaultArrayType int

	// sizedArrayType represents array with element type constraint and a size constraint
	sizedArrayType struct {
		sizeRange
		elementType dgo.Type
	}
)

// DefaultArrayType is the unconstrained Array type
const DefaultArrayType = defaultArrayType(0)

func arrayTypeOne(args []interface{}) dgo.ArrayType {
	switch a0 := Value(args[0]).(type) {
	case dgo.Integer:
		return newArrayType(nil, a0.GoInt(), dgo.UnboundedSize)
	case dgo.Type:
		return newArrayType(a0, 0, dgo.UnboundedSize)
	default:
		panic(illegalArgument(`Array`, `Type or Integer`, args, 0))
	}
}

func arrayTypeTwo(args []interface{}) dgo.ArrayType {
	a1, ok := Value(args[1]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Array`, `Integer`, args, 1))
	}
	switch a0 := Value(args[0]).(type) {
	case dgo.Integer:
		return newArrayType(nil, a0.GoInt(), a1.GoInt())
	case dgo.Type:
		return newArrayType(a0, a1.GoInt(), dgo.UnboundedSize)
	default:
		panic(illegalArgument(`Array`, `Type or Integer`, args, 0))
	}
}

func arrayTypeThree(args []interface{}) dgo.ArrayType {
	a0, ok := Value(args[0]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Array`, `Type`, args, 0))
	}
	a1, ok := Value(args[1]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Array`, `Integer`, args, 1))
	}
	a2, ok := Value(args[2]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`ArrayType`, `Integer`, args, 2))
	}
	return newArrayType(a0, a1.GoInt(), a2.GoInt())
}

// ArrayType returns a type that represents an Array value
func ArrayType(args []interface{}) dgo.ArrayType {
	switch len(args) {
	case 0:
		return DefaultArrayType
	case 1:
		return arrayTypeOne(args)
	case 2:
		return arrayTypeTwo(args)
	case 3:
		return arrayTypeThree(args)
	default:
		panic(illegalArgumentCount(`Array`, 0, 3, len(args)))
	}
}

func makeFrozenArray(slice []dgo.Value) dgo.Array {
	return &arrayFrozen{array{slice: slice}}
}

func newArrayType(elementType dgo.Type, min, max int64) dgo.ArrayType {
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	if elementType == nil {
		elementType = DefaultAnyType
	}
	if min == 0 && max == dgo.UnboundedSize && elementType == DefaultAnyType {
		// Unbounded
		return DefaultArrayType
	}
	if min == 1 && min == max && dgo.IsExact(elementType) {
		return makeFrozenArray([]dgo.Value{elementType}).(dgo.ArrayType)
	}
	return &sizedArrayType{sizeRange: sizeRange{min: uint32(min), max: uint32(max)}, elementType: elementType}
}

func (t defaultArrayType) Assignable(other dgo.Type) bool {
	_, ok := other.(dgo.ArrayType)
	return ok || CheckAssignableTo(nil, other, t)
}

func (t defaultArrayType) ElementType() dgo.Type {
	return DefaultAnyType
}

func (t defaultArrayType) Equals(other interface{}) bool {
	return t == other
}

func (t defaultArrayType) HashCode() dgo.Hash {
	return dgo.Hash(dgo.TiArray)
}

func (t defaultArrayType) Instance(value interface{}) bool {
	_, ok := value.(dgo.Array)
	return ok
}

func (t defaultArrayType) Max() int {
	return dgo.UnboundedSize
}

func (t defaultArrayType) Min() int {
	return 0
}

func (t defaultArrayType) New(arg dgo.Value) dgo.Value {
	return newArray(t, arg)
}

func (t defaultArrayType) ReflectType() reflect.Type {
	return reflect.SliceOf(reflectAnyType)
}

func (t defaultArrayType) String() string {
	return TypeString(t)
}

func (t defaultArrayType) Type() dgo.Type {
	return MetaType(t)
}

func (t defaultArrayType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArray
}

func (t defaultArrayType) Unbounded() bool {
	return true
}

func (t *sizedArrayType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *sizedArrayType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	switch ot := other.(type) {
	case defaultArrayType:
		return false // lacks size
	case dgo.Array:
		return t.DeepInstance(guard, ot)
	case dgo.ArrayType:
		return int(t.min) <= ot.Min() && ot.Max() <= int(t.max) && t.elementType.Assignable(ot.ElementType())
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *sizedArrayType) ElementType() dgo.Type {
	return t.elementType
}

func (t *sizedArrayType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *sizedArrayType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*sizedArrayType); ok {
		return t.min == ot.min && t.max == ot.max && equals(seen, t.elementType, ot.elementType)
	}
	return false
}

func (t *sizedArrayType) HashCode() dgo.Hash {
	return deepHashCode(nil, t)
}

func (t *sizedArrayType) deepHashCode(seen []dgo.Value) dgo.Hash {
	h := t.sizeRangeHash(dgo.TiArray)
	if DefaultAnyType != t.elementType {
		h = h*31 + deepHashCode(seen, t.elementType)
	}
	return h
}

func (t *sizedArrayType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *sizedArrayType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	if ov, ok := value.(arraySlice); ok {
		os := ov._slice()
		return t.inRange(len(os)) && allInstance(guard, t.elementType, os)
	}
	return false
}

func (t *sizedArrayType) New(arg dgo.Value) dgo.Value {
	return newArray(t, arg)
}

func (t *sizedArrayType) Resolve(ap dgo.AliasAdder) {
	te := t.elementType
	t.elementType = DefaultAnyType
	t.elementType = ap.Replace(te).(dgo.Type)
}

func (t *sizedArrayType) ReflectType() reflect.Type {
	return reflect.SliceOf(t.elementType.ReflectType())
}

func (t *sizedArrayType) String() string {
	return TypeString(t)
}

func (t *sizedArrayType) Type() dgo.Type {
	return MetaType(t)
}

func (t *sizedArrayType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArray
}

// Array returns a frozen dgo.Array that represents a copy of the given value. The value can be
// a slice or an Iterable
func Array(value interface{}) dgo.Array {
	switch value := value.(type) {
	case dgo.Array:
		return value.FrozenCopy().(dgo.Array)
	case dgo.Iterable:
		return arrayFromIterator(value.Len(), value.Each)
	case []dgo.Value:
		arr := make([]dgo.Value, len(value))
		for i := range value {
			e := value[i]
			if f, ok := e.(dgo.Mutability); ok {
				e = f.FrozenCopy()
			} else if e == nil {
				e = Nil
			}
			arr[i] = e
		}
		return makeFrozenArray(arr)
	case reflect.Value:
		return ValueFromReflected(value).(dgo.Array)
	default:
		return ValueFromReflected(reflect.ValueOf(value)).(dgo.Array)
	}
}

// arrayFromIterator creates an array from a size and an iterator goFunc. The
// iterator goFunc is expected to call its actor exactly size number of times.
func arrayFromIterator(size int, each func(dgo.Consumer)) dgo.Array {
	arr := make([]dgo.Value, size)
	i := 0
	each(func(e dgo.Value) {
		if f, ok := e.(dgo.Mutability); ok {
			e = f.FrozenCopy()
		}
		arr[i] = e
		i++
	})
	return makeFrozenArray(arr)
}

func sliceFromIterable(ir dgo.Iterable) []dgo.Value {
	es := make([]dgo.Value, ir.Len())
	i := 0
	ir.Each(func(e dgo.Value) {
		es[i] = e
		i++
	})
	return es
}

// ArrayFromReflected creates a new array that contains a copy of the given reflected slice
func ArrayFromReflected(vr reflect.Value, frozen bool) dgo.Value {
	if vr.IsNil() {
		return Nil
	}

	var arr []dgo.Value
	if vr.CanInterface() {
		ix := vr.Interface()
		if bs, ok := ix.([]byte); ok {
			return Binary(bs, frozen)
		}

		if vs, ok := ix.([]dgo.Value); ok {
			arr = vs
		}
	}

	if arr == nil {
		top := vr.Len()
		arr = make([]dgo.Value, top)
		for i := 0; i < top; i++ {
			arr[i] = ValueFromReflected(vr.Index(i))
		}
	}

	if frozen {
		arr = util.SliceCopy(arr)
		for i := range arr {
			if f, mutable := arr[i].(dgo.Mutability); mutable {
				arr[i] = f.FrozenCopy()
			}
		}
		return makeFrozenArray(arr)
	}
	return &array{slice: arr}
}

// ArrayWithCapacity creates a new mutable array of the given type and initial capacity.
func ArrayWithCapacity(capacity int) dgo.Array {
	return &array{slice: make([]dgo.Value, 0, capacity)}
}

// WrapSlice wraps the given slice in an array. Unset entries in the slice will be replaced by Nil.
func WrapSlice(values []dgo.Value) dgo.Array {
	ReplaceNil(values)
	return &array{slice: values}
}

// MutableValues returns a frozen dgo.Array that represents the given values
func MutableValues(values []interface{}) dgo.Array {
	cp := make([]dgo.Value, len(values))
	for i := range values {
		cp[i] = Value(values[i])
	}
	return &array{slice: cp}
}

func newArray(t dgo.Type, arg dgo.Value) dgo.Array {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`array`, 1, 1)
		arg = args.Get(0)
	}
	a := Array(arg)
	if !t.Instance(a) {
		panic(catch.Error(IllegalAssignment(t, a)))
	}
	return a
}

func valueSlice(values []interface{}, frozen bool) []dgo.Value {
	cp := make([]dgo.Value, len(values))
	if frozen {
		for i := range values {
			v := Value(values[i])
			if f, ok := v.(dgo.Mutability); ok {
				v = f.FrozenCopy()
			}
			cp[i] = v
		}
	} else {
		for i := range values {
			cp[i] = Value(values[i])
		}
	}
	return cp
}

// Integers returns a dgo.Array that represents the given ints
func Integers(values []int) dgo.Array {
	cp := make([]dgo.Value, len(values))
	for i := range values {
		cp[i] = intVal(values[i])
	}
	return makeFrozenArray(cp)
}

// Strings returns a dgo.Array that represents the given strings
func Strings(values []string) dgo.Array {
	cp := make([]dgo.Value, len(values))
	for i := range values {
		cp[i] = makeHString(values[i])
	}
	return makeFrozenArray(cp)
}

// Values returns a frozen dgo.Array that represents the given values
func Values(values []interface{}) dgo.Array {
	return &arrayFrozen{array{slice: valueSlice(values, true)}}
}

func (v *array) Add(vi interface{}) {
	v.slice = append(v.slice, Value(vi))
}

func (v *arrayFrozen) Add(_ interface{}) {
	panic(frozenArray(`Add`))
}

func (v *array) AddAll(values dgo.Iterable) {
	a := v.slice
	if ar, ok := values.(dgo.Array); ok {
		a = ar.AppendToSlice(a)
	} else {
		values.Each(func(e dgo.Value) { a = append(a, e) })
	}
	v.slice = a
}

func (v *arrayFrozen) AddAll(_ dgo.Iterable) {
	panic(frozenArray(`AddAll`))
}

func (v *array) AddValues(values ...interface{}) {
	v.slice = append(v.slice, valueSlice(values, false)...)
}

func (v *arrayFrozen) AddValues(_ ...interface{}) {
	panic(frozenArray(`AddValues`))
}

func (v *array) All(predicate dgo.Predicate) bool {
	a := v.slice
	for i := range a {
		if !predicate(a[i]) {
			return false
		}
	}
	return true
}

func (v *array) Any(predicate dgo.Predicate) bool {
	a := v.slice
	for i := range a {
		if predicate(a[i]) {
			return true
		}
	}
	return false
}

func (v *array) AppendTo(w dgo.Indenter) {
	w.AppendRune('{')
	ew := w.Indent()
	a := v.slice
	for i := range a {
		if i > 0 {
			w.AppendRune(',')
		}
		ew.NewLine()
		ew.AppendValue(v.slice[i])
	}
	w.NewLine()
	w.AppendRune('}')
}

func (v *array) AppendToSlice(slice []dgo.Value) []dgo.Value {
	return append(slice, v.slice...)
}

func (v *array) CompareTo(other interface{}) (int, bool) {
	return compare(nil, v, Value(other))
}

func (v *array) deepCompare(seen []dgo.Value, other deepCompare) (int, bool) {
	ov, ok := other.(arraySlice)
	if !ok {
		return 0, false
	}
	a := v.slice
	b := ov._slice()
	top := len(a)
	max := len(b)
	r := 0
	if top < max {
		r = -1
		max = top
	} else if top > max {
		r = 1
	}

	for i := 0; i < max; i++ {
		if _, ok = a[i].(dgo.Comparable); !ok {
			r = 0
			break
		}
		var c int
		if c, ok = compare(seen, a[i], b[i]); !ok {
			r = 0
			break
		}
		if c != 0 {
			r = c
			break
		}
	}
	return r, ok
}

func (v *array) Copy(frozen bool) dgo.Array {
	cp := util.SliceCopy(v.slice)
	if frozen {
		for i := range cp {
			if f, ok := cp[i].(dgo.Mutability); ok {
				cp[i] = f.FrozenCopy()
			}
		}
		return makeFrozenArray(cp)
	}
	for i := range cp {
		if f, ok := cp[i].(dgo.Mutability); ok {
			cp[i] = f.ThawedCopy()
		}
	}
	return &array{slice: cp}
}

func (v *arrayFrozen) Copy(frozen bool) dgo.Array {
	if frozen {
		return v
	}
	cp := util.SliceCopy(v.slice)
	for i := range cp {
		if f, ok := cp[i].(dgo.Mutability); ok {
			cp[i] = f.ThawedCopy()
		}
	}
	return &array{slice: cp}
}

func (v *array) ContainsAll(other dgo.Iterable) bool {
	return v._deepContainsAll(nil, other)
}

func (v *array) _deepContainsAll(seen []dgo.Value, other dgo.Iterable) bool {
	a := v.slice
	l := len(a)
	if l < other.Len() {
		return false
	}
	if l == 0 {
		return true
	}

	var vs []dgo.Value
	if oa, ok := other.(arraySlice); ok {
		vs = util.SliceCopy(oa._slice())
	} else {
		vs = sliceFromIterable(other)
	}

	// Keep track of elements that have been found equal using a copy
	// where such elements are set to nil. This avoids excessive calls
	// to Equals
	for i := range vs {
		ea := a[i]
		f := false
		for j := range vs {
			if be := vs[j]; be != nil {
				if equals(seen, be, ea) {
					vs[j] = nil
					f = true
					break
				}
			}
		}
		if !f {
			return false
		}
	}
	return true
}

func (v *array) Each(actor dgo.Consumer) {
	a := v.slice
	for i := range a {
		actor(a[i])
	}
}

func (v *array) EachWithIndex(actor dgo.DoWithIndex) {
	a := v.slice
	for i := range a {
		actor(a[i], i)
	}
}

func (v *array) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *array) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ov, ok := other.(arraySlice); ok {
		return sliceEquals(seen, v.slice, ov._slice())
	}
	return false
}

func (v *array) Find(finder dgo.Mapper) interface{} {
	a := v.slice
	for i := range a {
		if fv := finder(a[i]); fv != nil {
			return fv
		}
	}
	return nil
}

func (v *array) Flatten() dgo.Array {
	if a, changed := flattenSlice(v.slice); changed {
		return &array{slice: a}
	}
	return v
}

func (v *arrayFrozen) Flatten() dgo.Array {
	if a, changed := flattenSlice(v.slice); changed {
		return makeFrozenArray(a)
	}
	return v
}

func flattenSlice(a []dgo.Value) ([]dgo.Value, bool) {
	for i := range a {
		if _, ok := a[i].(arraySlice); ok {
			fs := make([]dgo.Value, i, len(a)*2)
			copy(fs, a)
			return flattenElements(a[i:], fs), true
		}
	}
	return a, false
}

func flattenElements(elements, receiver []dgo.Value) []dgo.Value {
	for i := range elements {
		e := elements[i]
		if a, ok := e.(arraySlice); ok {
			receiver = flattenElements(a._slice(), receiver)
		} else {
			receiver = append(receiver, e)
		}
	}
	return receiver
}

var (
	lb = []byte{'['}
	rb = []byte{']'}
	sp = []byte{' '}
	lc = []byte{'{'}
	rc = []byte{'}'}
	cl = []byte{':'}
	cm = []byte{',', ' '}
)

func (v *array) Format(s fmt.State, ch rune) {
	if ch == 'v' && s.Flag('#') {
		_, _ = s.Write(lb)
		_, _ = s.Write(rb)
		gt := Generic(v.Type()).(dgo.ArrayType)
		if nk, ok := gt.ElementType().(*nativeType); ok {
			_, _ = s.Write([]byte(nk.GoType().String()))
		} else {
			TypeStringOn(gt.ElementType(), s)
		}
		v.formatNV(s, ch)
	} else {
		_, _ = s.Write(lb)
		a := v.slice
		for i := range a {
			if i > 0 {
				_, _ = s.Write(sp)
			}
			formatValue(a[i], s, ch)
		}
		_, _ = s.Write(rb)
	}
}

func (v *array) formatNV(s fmt.State, ch rune) {
	_, _ = s.Write(lc)
	a := v.slice
	for i := range a {
		if i > 0 {
			_, _ = s.Write(cm)
		}
		formatValue(a[i], s, ch)
	}
	_, _ = s.Write(rc)
}

func (v *array) Frozen() bool {
	return false
}

func (v *arrayFrozen) Frozen() bool {
	return true
}

func (v *array) FrozenCopy() dgo.Value {
	return v.Copy(true)
}

func (v *arrayFrozen) FrozenCopy() dgo.Value {
	return v
}

func (v *array) ThawedCopy() dgo.Value {
	return v.Copy(false)
}

func (v *array) GoSlice() []dgo.Value {
	return v.slice
}

func (v *arrayFrozen) GoSlice() []dgo.Value {
	return util.SliceCopy(v.slice)
}

func (v *array) HashCode() dgo.Hash {
	return v.deepHashCode(nil)
}

func (v *array) deepHashCode(seen []dgo.Value) dgo.Hash {
	h := dgo.Hash(1)
	s := v.slice
	for i := range s {
		h = h*31 + deepHashCode(seen, s[i])
	}
	return h
}

func (v *array) Get(index int) dgo.Value {
	return v.slice[index]
}

func (v *array) IndexOf(vi interface{}) int {
	val := Value(vi)
	a := v.slice
	for i := range a {
		if val.Equals(a[i]) {
			return i
		}
	}
	return -1
}

func (v *array) Insert(pos int, vi interface{}) {
	v.slice = append(v.slice[:pos], append([]dgo.Value{Value(vi)}, v.slice[pos:]...)...)
}

func (v *arrayFrozen) Insert(_ int, _ interface{}) {
	panic(frozenArray(`Insert`))
}

// InterfaceSlice returns the values held by the Array as a slice. The slice will
// contain dgo.Value instances.
func (v *array) InterfaceSlice() []interface{} {
	s := v.slice
	is := make([]interface{}, len(s))
	for i := range s {
		is[i] = s[i]
	}
	return is
}

func (v *array) Len() int {
	return len(v.slice)
}

func (v *array) Map(mapper dgo.Mapper) dgo.Array {
	a := v.slice
	vs := make([]dgo.Value, len(a))
	for i := range a {
		vs[i] = Value(mapper(a[i]))
	}
	return &array{slice: vs}
}

func (v *arrayFrozen) Map(mapper dgo.Mapper) dgo.Array {
	a := v.slice
	vs := make([]dgo.Value, len(a))
	for i := range a {
		ev := Value(mapper(a[i]))
		if mv, mutable := ev.(dgo.Mutability); mutable {
			ev = mv.FrozenCopy()
		}
		vs[i] = ev
	}
	return makeFrozenArray(vs)
}

func (v *array) One(predicate dgo.Predicate) bool {
	a := v.slice
	f := false
	for i := range a {
		if predicate(a[i]) {
			if f {
				return false
			}
			f = true
		}
	}
	return f
}

func (v *array) Reduce(mi interface{}, reductor func(memo dgo.Value, elem dgo.Value) interface{}) dgo.Value {
	memo := Value(mi)
	a := v.slice
	for i := range a {
		memo = Value(reductor(memo, a[i]))
	}
	return memo
}

func (v *array) ReflectTo(value reflect.Value) {
	arrayReflectTo(v, value)
}

func (v *arrayFrozen) ReflectTo(value reflect.Value) {
	arrayReflectTo(v, value)
}

func arrayReflectTo(v dgo.Array, value reflect.Value) {
	vt := value.Type()
	ptr := vt.Kind() == reflect.Ptr
	if ptr {
		vt = vt.Elem()
	}
	if vt.Kind() == reflect.Interface && vt.Name() == `` {
		vt = v.Type().ReflectType()
	}
	a := v.(arraySlice)._slice()
	var s reflect.Value
	if !v.Frozen() && vt.Elem() == reflectValueType {
		s = reflect.ValueOf(a)
	} else {
		l := len(a)
		s = reflect.MakeSlice(vt, l, l)
		for i := range a {
			ReflectTo(a[i], s.Index(i))
		}
	}
	if ptr {
		// The created slice cannot be addressed. A pointer to it is necessary
		x := reflect.New(s.Type())
		x.Elem().Set(s)
		s = x
	}
	value.Set(s)
}

func (v *array) removePos(pos int) dgo.Value {
	a := v.slice
	if pos >= 0 && pos < len(a) {
		newLen := len(a) - 1
		val := a[pos]
		copy(a[pos:], a[pos+1:])
		a[newLen] = nil // release to GC
		v.slice = a[:newLen]
		return val
	}
	return nil
}

func (v *array) Remove(pos int) dgo.Value {
	return v.removePos(pos)
}

func (v *arrayFrozen) Remove(_ int) dgo.Value {
	panic(frozenArray(`Remove`))
}

func (v *array) RemoveValue(value interface{}) bool {
	return v.removePos(v.IndexOf(value)) != nil
}

func (v *arrayFrozen) RemoveValue(_ interface{}) bool {
	panic(frozenArray(`RemoveValue`))
}

func (v *array) Resolve(ap dgo.AliasAdder) {
	a := v.slice
	for i := range a {
		a[i] = ap.Replace(a[i])
	}
}

func (v *array) Reject(predicate dgo.Predicate) dgo.Array {
	return &array{slice: sliceReject(v.slice, predicate)}
}

func (v *arrayFrozen) Reject(predicate dgo.Predicate) dgo.Array {
	return makeFrozenArray(sliceReject(v.slice, predicate))
}

func sliceReject(a []dgo.Value, predicate dgo.Predicate) []dgo.Value {
	vs := make([]dgo.Value, 0)
	for i := range a {
		e := a[i]
		if !predicate(e) {
			vs = append(vs, e)
		}
	}
	return vs
}

func (v *array) SameValues(other dgo.Iterable) bool {
	return len(v.slice) == other.Len() && v.ContainsAll(other)
}

func (v *array) Select(predicate dgo.Predicate) dgo.Array {
	return &array{slice: sliceSelect(v.slice, predicate)}
}

func (v *arrayFrozen) Select(predicate dgo.Predicate) dgo.Array {
	return makeFrozenArray(sliceSelect(v.slice, predicate))
}

func sliceSelect(a []dgo.Value, predicate dgo.Predicate) []dgo.Value {
	vs := make([]dgo.Value, 0)
	for i := range a {
		e := a[i]
		if predicate(e) {
			vs = append(vs, e)
		}
	}
	return vs
}

func (v *array) Set(pos int, vi interface{}) dgo.Value {
	old := v.slice[pos]
	v.slice[pos] = Value(vi)
	return old
}

func (v *arrayFrozen) Set(_ int, _ interface{}) dgo.Value {
	panic(frozenArray(`Set`))
}

func (v *array) Slice(i, j int) dgo.Array {
	return &array{slice: v.slice[i:j]}
}

func (v *arrayFrozen) Slice(i, j int) dgo.Array {
	if i == 0 && j == len(v.slice) {
		return v
	}
	return makeFrozenArray(v.slice[i:j])
}

func (v *array) Sort() dgo.Array {
	if len(v.slice) < 2 {
		return v
	}
	return &array{slice: sortSlice(v.slice)}
}

func (v *arrayFrozen) Sort() dgo.Array {
	if len(v.slice) < 2 {
		return v
	}
	return makeFrozenArray(sortSlice(v.slice))
}

func sortSlice(sa []dgo.Value) []dgo.Value {
	sorted := util.SliceCopy(sa)
	sort.SliceStable(sorted, func(i, j int) bool {
		a := sorted[i]
		b := sorted[j]
		if ac, ok := a.(dgo.Comparable); ok {
			var c int
			if c, ok = ac.CompareTo(b); ok {
				return c < 0
			}
		}
		return a.Type().TypeIdentifier() < b.Type().TypeIdentifier()
	})
	return sorted
}

func (v *array) String() string {
	return TypeString(v)
}

func (v *array) ToMap() dgo.Map {
	return sliceToMap(v.slice, false)
}

func (v *arrayFrozen) ToMap() dgo.Map {
	return sliceToMap(v.slice, true)
}

func sliceToMap(ms []dgo.Value, frozen bool) dgo.Map {
	top := len(ms)

	ts := top / 2
	if top%2 != 0 {
		ts++
	}
	tbl := make([]*hashNode, tableSizeFor(ts))
	hl := len(tbl) - 1
	m := &hashMap{table: tbl, len: uint32(ts), frozen: frozen}

	for i := 0; i < top; {
		mk := ms[i]
		i++
		var mv dgo.Value = Nil
		if i < top {
			mv = ms[i]
			i++
		}
		hk := hl & hash(mk.HashCode())
		nd := &hashNode{mapEntry: mapEntry{key: mk, value: mv}, hashNext: tbl[hk], prev: m.last}
		if m.first == nil {
			m.first = nd
		} else {
			m.last.next = nd
		}
		m.last = nd
		tbl[hk] = nd
	}
	return m
}

func (v *array) ToMapFromEntries() (dgo.Map, bool) {
	return sliceToMapFromEntries(v.slice, false)
}

func (v *arrayFrozen) ToMapFromEntries() (dgo.Map, bool) {
	return sliceToMapFromEntries(v.slice, true)
}

func sliceToMapFromEntries(ms []dgo.Value, frozen bool) (dgo.Map, bool) {
	top := len(ms)
	tbl := make([]*hashNode, tableSizeFor(top))
	hl := len(tbl) - 1
	m := &hashMap{table: tbl, len: uint32(top), frozen: frozen}

	for i := range ms {
		nd, ok := ms[i].(*hashNode)
		if !ok {
			var ea arraySlice
			if ea, ok = ms[i].(arraySlice); ok {
				sl := ea._slice()
				if len(sl) != 2 {
					return nil, false
				}
				nd = &hashNode{mapEntry: mapEntry{key: sl[0], value: sl[1]}}
			} else {
				return nil, false
			}
		} else if nd.hashNext != nil {
			// Copy node, it belongs to another map
			c := *nd
			c.next = nil // this one might not get assigned below
			nd = &c
		}

		hk := hl & hash(nd.key.HashCode())
		nd.hashNext = tbl[hk]
		nd.prev = m.last
		if m.first == nil {
			m.first = nd
		} else {
			m.last.next = nd
		}
		m.last = nd
		tbl[hk] = nd
	}
	return m, true
}

func (v *array) Type() dgo.Type {
	return v
}

func (v *arrayFrozen) Type() dgo.Type {
	return v
}

func (v *array) Unique() dgo.Array {
	if a, changed := sliceUnique(v.slice); changed {
		return &array{slice: a}
	}
	return v
}

func (v *arrayFrozen) Unique() dgo.Array {
	if a, changed := sliceUnique(v.slice); changed {
		return makeFrozenArray(a)
	}
	return v
}

func sliceUnique(a []dgo.Value) ([]dgo.Value, bool) {
	top := len(a)
	if top < 2 {
		return a, false
	}
	tbl := make([]*hashNode, tableSizeFor(int(float64(top)/loadFactor)))
	hl := len(tbl) - 1
	u := make([]dgo.Value, top)
	ui := 0

nextVal:
	for i := range a {
		k := a[i]
		hk := hl & hash(k.HashCode())
		for e := tbl[hk]; e != nil; e = e.hashNext {
			if k.Equals(e.key) {
				continue nextVal
			}
		}
		tbl[hk] = &hashNode{mapEntry: mapEntry{key: k}, hashNext: tbl[hk]}
		u[ui] = k
		ui++
	}
	if ui == top {
		return a, false
	}
	return u[:ui], true
}

func (v *array) Pop() (dgo.Value, bool) {
	p := len(v.slice) - 1
	if p >= 0 {
		return v.removePos(p), true
	}
	return nil, false
}

func (v *arrayFrozen) Pop() (dgo.Value, bool) {
	panic(frozenArray(`Pop`))
}

func (v *array) _slice() []dgo.Value {
	return v.slice
}

func (v *array) With(vi interface{}) dgo.Array {
	return &array{slice: append(v.slice, Value(vi))}
}

func (v *arrayFrozen) With(vi interface{}) dgo.Array {
	ev := Value(vi)
	if vm, mutable := ev.(dgo.Mutability); mutable {
		ev = vm.FrozenCopy()
	}
	return makeFrozenArray(append(v.slice, ev))
}

func (v *array) WithAll(values dgo.Iterable) dgo.Array {
	n := values.Len()
	if n == 0 {
		return v
	}
	a := make([]dgo.Value, len(v.slice), len(v.slice)+n)
	copy(a, v.slice)
	values.Each(func(ev dgo.Value) { a = append(a, ev) })
	return &array{slice: a}
}

func (v *arrayFrozen) WithAll(values dgo.Iterable) dgo.Array {
	n := values.Len()
	if n == 0 {
		return v
	}
	values = values.FrozenCopy().(dgo.Iterable)
	a := make([]dgo.Value, len(v.slice), len(v.slice)+n)
	copy(a, v.slice)
	values.Each(func(ev dgo.Value) { a = append(a, ev) })
	return makeFrozenArray(a)
}

func (v *array) WithValues(values ...interface{}) dgo.Array {
	if len(values) == 0 {
		return v
	}
	return &array{slice: append(v.slice, valueSlice(values, false)...)}
}

func (v *arrayFrozen) WithValues(values ...interface{}) dgo.Array {
	if len(values) == 0 {
		return v
	}
	return makeFrozenArray(append(v.slice, valueSlice(values, true)...))
}

// Array as it's own ExactType below

func (v *array) Assignable(other dgo.Type) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v *array) ElementTypeAt(index int) dgo.Type {
	return v.slice[index].Type()
}

func (v *array) ElementType() dgo.Type {
	es := v.slice
	switch len(es) {
	case 0:
		return DefaultAnyType
	case 1:
		return es[0].Type()
	}
	return (*allOfValueType)(v)
}

func (v *arrayFrozen) ElementType() dgo.Type {
	es := v.slice
	switch len(es) {
	case 0:
		return DefaultAnyType
	case 1:
		return es[0].Type()
	}
	return (*allOfValueType)(&v.array)
}

func (v *array) ElementTypes() dgo.Array {
	es := v.slice
	ts := make([]dgo.Value, len(es))
	for i := range es {
		ts[i] = es[i].Type()
	}
	return makeFrozenArray(ts)
}

func (v *array) Generic() dgo.Type {
	return &sizedArrayType{sizeRange: sizeRange{min: 0, max: dgo.UnboundedSize}, elementType: Generic(v.ElementType())}
}

func (v *array) Instance(value interface{}) bool {
	return v.Equals(value)
}

func (v *array) Max() int {
	return v.Len()
}

func (v *array) Min() int {
	return v.Len()
}

func (v *array) New(arg dgo.Value) dgo.Value {
	return newArray(v, arg)
}

func (v *array) ReflectType() reflect.Type {
	return reflect.SliceOf(v.ElementType().ReflectType())
}

func (v *array) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArrayExact
}

func (v *array) Unbounded() bool {
	return false
}

func (v *array) Variadic() bool {
	return false
}

// ReplaceNil performs an in-place replacement of nil interfaces with the NilValue
func ReplaceNil(vs []dgo.Value) {
	for i := range vs {
		if vs[i] == nil {
			vs[i] = Nil
		}
	}
}

// allInstance returns true when all elements of slice vs are assignable to the given type t
func allInstance(guard dgo.RecursionGuard, t dgo.Type, vs []dgo.Value) bool {
	if t == DefaultAnyType {
		return true
	}
	for i := range vs {
		if !Instance(guard, t, vs[i]) {
			return false
		}
	}
	return true
}

func frozenArray(f string) error {
	return catch.Error(`%s called on a frozen Array`, f)
}

func resolveSlice(ts []dgo.Value, ap dgo.AliasAdder) {
	for i := range ts {
		ts[i] = ap.Replace(ts[i])
	}
}
