package internal

import (
	"reflect"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
)

type (
	// tupleType represents an array with an exact number of ordered element types.
	tupleType struct {
		types    []dgo.Value
		variadic bool
	}
)

// DefaultTupleType is a tuple without size and type constraints
var DefaultTupleType = &tupleType{variadic: true, types: []dgo.Value{DefaultAnyType}}

// EmptyTupleType is a tuple that represents an empty array
var EmptyTupleType = &tupleType{variadic: false, types: []dgo.Value{}}

// TupleType creates a new TupleType based on the given types
func TupleType(types []interface{}) dgo.TupleType {
	return newTupleType(types, false)
}

// VariadicTupleType returns a type that represents an Array value with a variadic number of elements. EachEntryType
// given type determines the type of a corresponding element in an array except for the last one which
// determines the remaining elements.
func VariadicTupleType(types []interface{}) dgo.TupleType {
	n := len(types)
	if n == 0 {
		panic(catch.Error(`a variadic tuple must have at least one element`))
	}
	return newTupleType(types, true)
}

func newTupleType(types []interface{}, variadic bool) dgo.TupleType {
	l := len(types)
	if l == 0 {
		return EmptyTupleType
	}
	if variadic && l == 1 && DefaultAnyType.Equals(types[0]) {
		return DefaultTupleType
	}
	exact := !variadic
	es := make([]dgo.Value, l)
	for i := 0; i < l; i++ {
		et := types[i].(dgo.Type)
		if exact && !dgo.IsExact(et) {
			exact = false
		}
		es[i] = et
	}
	if exact {
		return makeFrozenArray(es).(dgo.TupleType)
	}
	return &tupleType{types: es, variadic: variadic}
}

func (t *tupleType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *tupleType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	return tupleAssignable(guard, t, other)
}

func tupleAssignableTuple(guard dgo.RecursionGuard, t, ot dgo.TupleType) bool {
	if t.Min() > ot.Min() || ot.Max() > t.Max() {
		return false
	}

	var tv, ov dgo.Type
	tn := t.Len()
	if t.Variadic() {
		tn--
		tv = t.ElementTypeAt(tn)
	}
	on := ot.Len()
	if ot.Variadic() {
		on--
		ov = ot.ElementTypeAt(on)
	}

	// n := max(tn, on)
	n := tn
	if n < on {
		n = on
	}

	for i := 0; i < n; i++ {
		te := tv
		if i < tn {
			te = t.ElementTypeAt(i)
		}

		oe := ov
		if i < on {
			oe = ot.ElementTypeAt(i)
		}
		if te == nil || oe == nil || !Assignable(guard, te, oe) {
			return false
		}
	}
	return true
}

func tupleAssignableArray(guard dgo.RecursionGuard, t dgo.TupleType, ot *sizedArrayType) bool {
	if t.Min() <= ot.Min() && ot.Max() <= t.Max() {
		et := ot.ElementType()
		n := t.Len()
		if t.Variadic() {
			n--
		}
		for i := 0; i < n; i++ {
			if !Assignable(guard, t.ElementTypeAt(i), et) {
				return false
			}
		}
		return !t.Variadic() || Assignable(guard, t.ElementTypeAt(n), et)
	}
	return false
}

func tupleAssignable(guard dgo.RecursionGuard, t dgo.TupleType, other dgo.Type) bool {
	switch ot := other.(type) {
	case defaultArrayType:
		return false
	case *array, *arrayFrozen:
		return Instance(guard, t, ot)
	case dgo.TupleType:
		return tupleAssignableTuple(guard, t, ot)
	case *sizedArrayType:
		return tupleAssignableArray(guard, t, ot)
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *tupleType) ElementTypeAt(index int) dgo.Type {
	return t.types[index].(dgo.Type)
}

func (t *tupleType) ElementType() dgo.Type {
	return tupleElementType(t)
}

func tupleElementType(t dgo.TupleType) (et dgo.Type) {
	switch t.Len() {
	case 0:
		et = DefaultAnyType
	case 1:
		et = t.ElementTypeAt(0)
	default:
		ea := t.ElementTypes().Unique()
		if ea.Len() == 1 {
			return ea.Get(0).(dgo.Type)
		}
		return &allOfType{slice: ea.(arraySlice)._slice()}
	}
	return
}

func (t *tupleType) ElementTypes() dgo.Array {
	return makeFrozenArray(t.types)
}

func (t *tupleType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *tupleType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*tupleType); ok {
		return t.variadic == ot.variadic && sliceEquals(seen, t.types, ot.types)
	}
	return tupleEquals(seen, t, other)
}

func tupleEquals(seen []dgo.Value, t dgo.TupleType, other interface{}) bool {
	if ot, ok := other.(dgo.TupleType); ok {
		n := t.Len()
		if t.Variadic() == ot.Variadic() && n == ot.Len() {
			for i := 0; i < n; i++ {
				if !equals(seen, t.ElementTypeAt(i), ot.ElementTypeAt(i)) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (t *tupleType) Generic() dgo.Type {
	return newArrayType(Generic(t.ElementType()), 0, dgo.UnboundedSize)
}

func (t *tupleType) HashCode() dgo.Hash {
	return deepHashCode(nil, t)
}

func (t *tupleType) deepHashCode(seen []dgo.Value) dgo.Hash {
	return tupleHashCode(t, seen)
}

func tupleHashCode(t dgo.TupleType, seen []dgo.Value) dgo.Hash {
	h := dgo.Hash(1)
	if t.Variadic() {
		h = 7
	}
	l := t.Len()
	for i := 0; i < l; i++ {
		h = h*31 + deepHashCode(seen, t.ElementTypeAt(i))
	}
	return h
}

func (t *tupleType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *tupleType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	return tupleInstance(guard, t, value)
}

func tupleInstance(guard dgo.RecursionGuard, t dgo.TupleType, value interface{}) bool {
	ov, ok := value.(arraySlice)
	if !ok {
		return false
	}

	s := ov._slice()
	n := len(s)
	if t.Variadic() {
		if t.Min() > n {
			return false
		}
		tn := t.Len() - 1
		for i := 0; i < tn; i++ {
			if !Instance(guard, t.ElementTypeAt(i), s[i]) {
				return false
			}
		}
		vt := t.ElementTypeAt(tn)
		for ; tn < n; tn++ {
			if !Instance(guard, vt, s[tn]) {
				return false
			}
		}
		return true
	}

	if n != t.Len() {
		return false
	}

	for i := range s {
		if !Instance(guard, t.ElementTypeAt(i), s[i]) {
			return false
		}
	}
	return true
}

func (t *tupleType) Len() int {
	return len(t.types)
}

func (t *tupleType) Max() int {
	return tupleMax(t)
}

func tupleMax(t dgo.TupleType) int {
	n := t.Len()
	if t.Variadic() {
		n = dgo.UnboundedSize
	}
	return n
}

func (t *tupleType) Min() int {
	return tupleMin(t)
}

func (t *tupleType) New(arg dgo.Value) dgo.Value {
	return newArray(t, arg)
}

func tupleMin(t dgo.TupleType) int {
	n := t.Len()
	if t.Variadic() {
		n--
	}
	return n
}

func (t *tupleType) ReflectType() reflect.Type {
	return reflect.SliceOf(t.ElementType().ReflectType())
}

func (t *tupleType) Resolve(ap dgo.AliasAdder) {
	s := t.types
	t.types = nil
	resolveSlice(s, ap)
	t.types = s
}

func (t *tupleType) String() string {
	return TypeString(t)
}

func (t *tupleType) Type() dgo.Type {
	return MetaType(t)
}

func (t *tupleType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiTuple
}

func (t *tupleType) Unbounded() bool {
	return t.variadic
}

func (t *tupleType) Variadic() bool {
	return t.variadic
}
