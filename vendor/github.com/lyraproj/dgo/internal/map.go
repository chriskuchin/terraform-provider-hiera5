package internal

import (
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

const initialCapacity = 1 << 3
const maximumCapacity = 1 << 30
const loadFactor = 0.75

type (
	// defaultMapType is the unconstrained map type
	defaultMapType int

	// exactMapType represents a map exactly
	exactMapType struct {
		deepExactType
		value dgo.Map
	}

	// sizedMapType represents a map with constraints on key type, value type, and size
	sizedMapType struct {
		keyType   dgo.Type
		valueType dgo.Type
		min       int
		max       int
	}

	mapEntry struct {
		key   dgo.Value
		value dgo.Value
	}

	exactEntryType struct {
		deepExactType
		value *mapEntry
	}

	hashNode struct {
		mapEntry
		hashNext *hashNode
		next     *hashNode
		prev     *hashNode
	}

	// hashMap is an unsorted Map that uses a hash table
	hashMap struct {
		table  []*hashNode
		len    int
		first  *hashNode
		last   *hashNode
		frozen bool
	}
)

func (t *exactEntryType) Generic() dgo.Type {
	return DefaultAnyType
}

func (t *exactEntryType) ReflectType() reflect.Type {
	return reflect.TypeOf((*dgo.MapEntry)(nil)).Elem()
}

func (t *exactEntryType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMapEntryExact
}

// NewMapEntry returns a new MapEntry instance with the given key and value
func NewMapEntry(key, value interface{}) dgo.MapEntry {
	return &mapEntry{Value(key), Value(value)}
}

func (t *exactEntryType) ExactValue() dgo.Value {
	return t.value
}

func (v *mapEntry) AppendTo(w dgo.Indenter) {
	w.AppendValue(v.key)
	w.Append(`:`)
	if w.Indenting() {
		w.Append(` `)
	}
	w.AppendValue(v.value)
}

func (v *mapEntry) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *mapEntry) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ov, ok := other.(dgo.MapEntry); ok {
		return equals(seen, v.key, ov.Key()) && equals(seen, v.value, ov.Value())
	}
	return false
}

func (v *mapEntry) HashCode() int {
	return deepHashCode(nil, v)
}

func (v *mapEntry) deepHashCode(seen []dgo.Value) int {
	return deepHashCode(seen, v.key) ^ deepHashCode(seen, v.value)
}

func (v *mapEntry) Freeze() {
	if f, ok := v.value.(dgo.Freezable); ok {
		f.Freeze()
	}
}

func (v *mapEntry) Frozen() bool {
	if f, ok := v.value.(dgo.Freezable); ok && !f.Frozen() {
		return false
	}
	return true
}

func (v *mapEntry) FrozenCopy() dgo.Value {
	if v.Frozen() {
		return v
	}
	c := &mapEntry{key: v.key, value: v.value}
	c.copyFreeze()
	return c
}

func (v *mapEntry) ThawedCopy() dgo.Value {
	c := &mapEntry{key: v.key, value: v.value}
	if f, ok := v.value.(dgo.Freezable); ok {
		v.value = f.ThawedCopy()
	}
	return c
}

func (v *hashNode) FrozenCopy() dgo.Value {
	if v.Frozen() {
		return v
	}
	c := &hashNode{mapEntry: mapEntry{key: v.key, value: v.value}}
	c.copyFreeze()
	return c
}

func (v *hashNode) ThawedCopy() dgo.Value {
	c := &hashNode{mapEntry: mapEntry{key: v.key, value: v.value}}
	if f, ok := v.value.(dgo.Freezable); ok {
		v.value = f.ThawedCopy()
	}
	return c
}

func (v *mapEntry) Key() dgo.Value {
	return v.key
}

func (v *mapEntry) String() string {
	return util.ToStringERP(v)
}

func (v *mapEntry) Type() dgo.Type {
	ea := &exactEntryType{value: v}
	ea.ExactType = ea
	return ea
}

func (v *mapEntry) Value() dgo.Value {
	return v.value
}

func (v *mapEntry) copyFreeze() {
	if f, ok := v.value.(dgo.Freezable); ok {
		v.value = f.FrozenCopy()
	}
}

func (v *mapEntry) copyThaw() {
	if f, ok := v.value.(dgo.Freezable); ok {
		v.value = f.ThawedCopy()
	}
}

var emptyMap = &hashMap{frozen: true}

// Map creates an immutable dgo.Map from the given slice which must have 0, 1, or an
// even number of elements.
//
// Zero elements: the empty map is returned.
//
// One element: must be a go map, a go struct, or an Array with an even number of elements.
//
// An even number of elements: will be considered a flat list of key, value [, key, value, ... ]
func Map(args []interface{}) dgo.Map {
	return mapFromArgs(args, true)
}

// MutableMap creates a mutable dgo.Map from the given slice which must have 0, 1, or an
// even number of elements.
//
// Zero elements: the empty map is returned.
//
// One element: must be a go map, a go struct, or an Array with an even number of elements.
//
// An even number of elements: will be considered a flat list of key, value [, key, value, ... ]
func MutableMap(args []interface{}) dgo.Map {
	return mapFromArgs(args, false)
}

func newMap(t dgo.Type, arg dgo.Value) dgo.Map {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`map`, 1, 1)
		arg = args.Get(0)
	}
	m, ok := arg.(dgo.Map)
	if ok {
		m = m.FrozenCopy().(dgo.Map)
	} else {
		m = mapFromArgs([]interface{}{arg}, true)
	}
	if !t.Instance(m) {
		panic(IllegalAssignment(t, m))
	}
	return m
}

func mapFromArgs(args []interface{}, frozen bool) dgo.Map {
	l := len(args)
	switch {
	case l == 0:
		if frozen {
			return emptyMap
		}
		return MapWithCapacity(0)
	case l == 1:
		a0 := args[0]
		if ar, ok := a0.(dgo.Array); ok && ar.Len()%2 == 0 {
			if frozen {
				return ar.FrozenCopy().(dgo.Array).ToMap()
			}
			m := ar.ToMap().(*hashMap)
			m.frozen = false
			return m
		}
		rm := reflect.ValueOf(a0)
		switch rm.Kind() {
		case reflect.Map:
			return FromReflectedMap(rm, frozen).(dgo.Map)
		case reflect.Ptr:
			re := rm.Elem()
			if re.Kind() == reflect.Struct {
				return FromReflectedStruct(re)
			}
		case reflect.Struct:
			return FromReflectedStruct(rm)
		}
		panic(fmt.Errorf(`illegal argument: %t is not a map, a struct, or an array with even number of elements`, a0))
	case l%2 == 0:
		if frozen {
			return Values(args).ToMap()
		}
		return MutableValues(args).ToMap()
	default:
		panic(fmt.Errorf(`the number of arguments to Map must be 1 or an even number, got: %d`, l))
	}
}

// FromReflectedMap creates a Map from a reflected map. It panics if rm's kind is not reflect.Map. If frozen is true,
// the created Map will be immutable and the type will reflect exactly that map and nothing else. If frozen is false,
// the created Map will be mutable and its type will be derived from the reflected map.
func FromReflectedMap(rm reflect.Value, frozen bool) dgo.Value {
	if rm.IsNil() {
		return Nil
	}
	keys := rm.MapKeys()
	top := len(keys)
	ic := top
	if top == 0 {
		if frozen {
			return emptyMap
		}
		ic = initialCapacity
	}
	tbl := make([]*hashNode, tableSizeFor(int(float64(ic)/loadFactor)))
	hl := len(tbl) - 1
	se := make([][2]dgo.Value, len(keys))
	for i := range keys {
		key := keys[i]
		se[i] = [2]dgo.Value{ValueFromReflected(key), ValueFromReflected(rm.MapIndex(key))}
	}

	// Sort by key to always get predictable order
	sort.Slice(se, func(i, j int) bool {
		less := false
		if cm, ok := se[i][0].(dgo.Comparable); ok {
			var c int
			if c, ok = cm.CompareTo(se[j][0]); ok {
				less = c < 0
			}
		}
		return less
	})

	m := &hashMap{table: tbl, len: top, frozen: frozen}
	for i := range se {
		e := se[i]
		k := e[0]
		hk := hl & hash(k.HashCode())
		hn := &hashNode{mapEntry: mapEntry{key: k, value: e[1]}, hashNext: tbl[hk], prev: m.last}
		if m.last == nil {
			m.first = hn
		} else {
			m.last.next = hn
		}
		m.last = hn
		tbl[hk] = hn
	}
	return m
}

// FromReflectedStruct creates a frozen Map from the exported fields of a go struct. It panics if rm's kind is not
// reflect.Struct.
func FromReflectedStruct(rv reflect.Value) dgo.Struct {
	return &structVal{rs: rv, frozen: false}
}

// MapWithCapacity creates an empty dgo.Map suitable to hold a given number of entries.
func MapWithCapacity(capacity int) dgo.Map {
	if capacity <= 0 {
		capacity = initialCapacity
	}
	capacity = int(float64(capacity) / loadFactor)
	return &hashMap{table: make([]*hashNode, tableSizeFor(capacity)), len: 0, frozen: false}
}

func (g *hashMap) All(predicate dgo.EntryPredicate) bool {
	for e := g.first; e != nil; e = e.next {
		if !predicate(e) {
			return false
		}
	}
	return true
}

func (g *hashMap) AllKeys(predicate dgo.Predicate) bool {
	for e := g.first; e != nil; e = e.next {
		if !predicate(e.key) {
			return false
		}
	}
	return true
}

func (g *hashMap) AllValues(predicate dgo.Predicate) bool {
	for e := g.first; e != nil; e = e.next {
		if !predicate(e.value) {
			return false
		}
	}
	return true
}

func (g *hashMap) Any(predicate dgo.EntryPredicate) bool {
	for e := g.first; e != nil; e = e.next {
		if predicate(e) {
			return true
		}
	}
	return false
}

func (g *hashMap) AnyKey(predicate dgo.Predicate) bool {
	for e := g.first; e != nil; e = e.next {
		if predicate(e.key) {
			return true
		}
	}
	return false
}

func (g *hashMap) AnyValue(predicate dgo.Predicate) bool {
	for e := g.first; e != nil; e = e.next {
		if predicate(e.value) {
			return true
		}
	}
	return false
}

func (g *hashMap) AppendTo(w dgo.Indenter) {
	appendMapTo(g, w)
}

func appendMapTo(m dgo.Map, w dgo.Indenter) {
	w.AppendRune('{')
	ew := w.Indent()
	first := true
	m.EachEntry(func(e dgo.MapEntry) {
		if first {
			first = false
		} else {
			ew.AppendRune(',')
		}
		ew.NewLine()
		ew.AppendValue(e)
	})
	w.NewLine()
	w.AppendRune('}')
}

func (g *hashMap) ContainsKey(key interface{}) bool {
	return g.Get(key) != nil
}

func (g *hashMap) Copy(frozen bool) dgo.Map {
	if frozen && g.frozen {
		return g
	}

	c := &hashMap{len: g.len, frozen: frozen}
	g.resize(c, 0)
	if frozen {
		for e := c.first; e != nil; e = e.next {
			e.copyFreeze()
		}
	} else {
		for e := c.first; e != nil; e = e.next {
			e.copyThaw()
		}
	}
	return c
}

func (g *hashMap) Each(actor dgo.Consumer) {
	for e := g.first; e != nil; e = e.next {
		actor(e)
	}
}

func (g *hashMap) EachEntry(actor dgo.EntryActor) {
	for e := g.first; e != nil; e = e.next {
		actor(e)
	}
}

func (g *hashMap) EachKey(actor dgo.Consumer) {
	for e := g.first; e != nil; e = e.next {
		actor(e.key)
	}
}

func (g *hashMap) EachValue(actor dgo.Consumer) {
	for e := g.first; e != nil; e = e.next {
		actor(e.value)
	}
}

func (g *hashMap) Equals(other interface{}) bool {
	return equals(nil, g, other)
}

func (g *hashMap) deepEqual(seen []dgo.Value, other deepEqual) bool {
	return mapEqual(seen, g, other)
}

func mapEqual(seen []dgo.Value, g dgo.Map, other deepEqual) bool {
	if om, ok := other.(dgo.Map); ok && g.Len() == om.Len() {
		return g.All(func(e dgo.MapEntry) bool { return equals(seen, e.Value(), om.Get(e.Key())) })
	}
	return false
}

func (g *hashMap) Find(predicate dgo.EntryPredicate) dgo.MapEntry {
	for e := g.first; e != nil; e = e.next {
		if predicate(e) {
			return e
		}
	}
	return nil
}

func (g *hashMap) Freeze() {
	if !g.frozen {
		g.frozen = true
		for e := g.first; e != nil; e = e.next {
			e.Freeze()
		}
	}
}

func (g *hashMap) Frozen() bool {
	return g.frozen
}

func (g *hashMap) FrozenCopy() dgo.Value {
	return g.Copy(true)
}

func (g *hashMap) ThawedCopy() dgo.Value {
	return g.Copy(false)
}

func (g *hashMap) Get(key interface{}) dgo.Value {
	tbl := g.table
	tl := len(tbl) - 1
	if tl >= 0 {
		// This switch increases performance a great deal because using the direct implementation
		// instead of the dgo.Value enables inlining of the HashCode() method
		switch k := key.(type) {
		case *hstring:
			for e := tbl[tl&hash(k.HashCode())]; e != nil; e = e.hashNext {
				if k.Equals(e.key) {
					return e.value
				}
			}
		case intVal:
			for e := tbl[tl&hash(k.HashCode())]; e != nil; e = e.hashNext {
				if k == e.key {
					return e.value
				}
			}
		case string:
			gk := makeHString(k)
			for e := tbl[tl&hash(gk.HashCode())]; e != nil; e = e.hashNext {
				if gk.Equals(e.key) {
					return e.value
				}
			}
		default:
			gk := Value(k)
			for e := tbl[tl&hash(gk.HashCode())]; e != nil; e = e.hashNext {
				if gk.Equals(e.key) {
					return e.value
				}
			}
		}
	}
	return nil
}

func (g *hashMap) HashCode() int {
	return deepHashCode(nil, g)
}

func (g *hashMap) deepHashCode(seen []dgo.Value) int {
	// compute order independent hash code. This is necessary to withhold the
	// contract that when two maps are equal, their hashes are equal.
	hs := make([]int, g.len)
	i := 0
	for e := g.first; e != nil; e = e.next {
		hs[i] = deepHashCode(seen, e)
		i++
	}
	sort.Ints(hs)
	h := 1
	for i = range hs {
		h = h*31 + hs[i]
	}
	return h
}

func (g *hashMap) Keys() dgo.Array {
	return arrayFromIterator(g.len, g.EachKey)
}

func (g *hashMap) Len() int {
	return g.len
}

func (g *hashMap) Map(mapper dgo.EntryMapper) dgo.Map {
	c := &hashMap{len: g.len, frozen: g.frozen}
	g.resize(c, 0)
	for e := c.first; e != nil; e = e.next {
		e.value = Value(mapper(e))
	}
	return c
}

func (g *hashMap) Merge(associations dgo.Map) dgo.Map {
	if associations.Len() == 0 || g == associations {
		return g
	}
	l := g.len
	if l == 0 {
		return associations
	}
	c := &hashMap{len: l}
	g.resize(c, l+associations.Len())
	c.PutAll(associations)
	c.frozen = g.frozen
	return c
}

func (g *hashMap) Put(ki, vi interface{}) dgo.Value {
	if g.frozen {
		panic(frozenMap(`Put`))
	}
	k := Value(ki)
	v := Value(vi)
	hs := hash(k.HashCode())
	var hk int

	tbl := g.table
	if tbl == nil {
		tbl = make([]*hashNode, tableSizeFor(1))
		g.table = tbl
	}
	hk = (len(tbl) - 1) & hs
	for e := tbl[hk]; e != nil; e = e.hashNext {
		if k.Equals(e.key) {
			old := e.value
			e.value = v
			return old
		}
	}

	if float64(g.len+1) > float64(len(g.table))*loadFactor {
		g.resize(g, 1)
		tbl = g.table
		hk = (len(tbl) - 1) & hs
	}

	nd := &hashNode{mapEntry: mapEntry{key: frozenCopy(k), value: v}, hashNext: tbl[hk], prev: g.last}
	if g.first == nil {
		g.first = nd
	} else {
		g.last.next = nd
	}
	g.last = nd
	tbl[hk] = nd
	g.len++
	return nil
}

func (g *hashMap) PutAll(associations dgo.Map) {
	al := associations.Len()
	if al == 0 {
		return
	}
	if g.frozen {
		panic(frozenMap(`PutAll`))
	}

	l := g.len
	if float64(l+al) > float64(l)*loadFactor {
		g.resize(g, al)
	}
	tbl := g.table

	associations.EachEntry(func(entry dgo.MapEntry) {
		key := entry.Key()
		val := entry.Value()
		hk := (len(tbl) - 1) & hash(key.HashCode())
		for e := tbl[hk]; e != nil; e = e.hashNext {
			if key.Equals(e.key) {
				e.value = val
				return
			}
		}
		nd := &hashNode{mapEntry: mapEntry{key: frozenCopy(key), value: val}, hashNext: tbl[hk], prev: g.last}
		if g.first == nil {
			g.first = nd
		} else {
			g.last.next = nd
		}
		g.last = nd
		tbl[hk] = nd
		l++
	})
	g.len = l
}

func (g *hashMap) ReflectTo(value reflect.Value) {
	ht := value.Type()
	ptr := ht.Kind() == reflect.Ptr
	if ptr {
		ht = ht.Elem()
	}
	if ht.Kind() == reflect.Interface && ht.Name() == `` {
		ht = g.Type().ReflectType()
	}
	keyType := ht.Key()
	valueType := ht.Elem()
	m := reflect.MakeMapWithSize(ht, g.Len())
	g.EachEntry(func(e dgo.MapEntry) {
		rk := reflect.New(keyType).Elem()
		ReflectTo(e.Key(), rk)
		rv := reflect.New(valueType).Elem()
		ReflectTo(e.Value(), rv)
		m.SetMapIndex(rk, rv)
	})
	if ptr {
		// The created map cannot be addressed. A pointer to it is necessary
		x := reflect.New(m.Type())
		x.Elem().Set(m)
		m = x
	}
	value.Set(m)
}

func (g *hashMap) Remove(ki interface{}) dgo.Value {
	if g.frozen {
		panic(frozenMap(`Remove`))
	}
	key := Value(ki)
	hk := (len(g.table) - 1) & hash(key.HashCode())

	var p *hashNode
	for e := g.table[hk]; e != nil; e = e.hashNext {
		if key.Equals(e.key) {
			old := e.value
			if p == nil {
				g.table[hk] = e.hashNext
			} else {
				p.hashNext = e.hashNext
			}
			if e.prev == nil {
				g.first = e.next
			} else {
				e.prev.next = e.next
			}
			if e.next == nil {
				g.last = e.prev
			} else {
				e.next.prev = e.prev
			}
			g.len--
			return old
		}
		p = e
	}
	return nil
}

func (g *hashMap) RemoveAll(keys dgo.Array) {
	if g.frozen {
		panic(frozenMap(`RemoveAll`))
	}
	if g.len == 0 || keys.Len() == 0 {
		return
	}

	tbl := g.table
	kl := len(tbl) - 1
	keys.Each(func(k dgo.Value) {
		hk := kl & hash(k.HashCode())
		var p *hashNode
		for e := tbl[hk]; e != nil; e = e.hashNext {
			if k.Equals(e.key) {
				if p == nil {
					tbl[hk] = e.hashNext
				} else {
					p.hashNext = e.hashNext
				}
				if e.prev == nil {
					g.first = e.next
				} else {
					e.prev.next = e.next
				}
				if e.next == nil {
					g.last = e.prev
				} else {
					e.next.prev = e.prev
				}
				g.len--
				break
			}
			p = e
		}
	})
}

func (g *hashMap) Resolve(ap dgo.AliasAdder) {
	for e := g.first; e != nil; e = e.next {
		e.value = ap.Replace(e.value)
	}
}

func (g *hashMap) String() string {
	return util.ToStringERP(g)
}

func (g *hashMap) StringKeys() bool {
	for e := g.first; e != nil; e = e.next {
		if _, str := e.key.(*hstring); !str {
			return false
		}
	}
	return true
}

func (g *hashMap) With(ki, vi interface{}) dgo.Map {
	var c *hashMap
	key := Value(ki)
	val := Value(vi)
	if g.table == nil {
		c = &hashMap{table: make([]*hashNode, tableSizeFor(initialCapacity)), len: g.len}
	} else {
		if val.Equals(g.Get(key)) {
			return g
		}
		c = &hashMap{len: g.len}
		g.resize(c, 1)
	}
	c.Put(key, val)
	c.frozen = g.frozen
	return c
}

func (g *hashMap) Without(ki interface{}) dgo.Map {
	key := Value(ki)
	if g.Get(key) == nil {
		return g
	}
	c := &hashMap{len: g.len}
	g.resize(c, 0)
	c.Remove(key)
	c.frozen = g.frozen
	return c
}

func (g *hashMap) WithoutAll(keys dgo.Array) dgo.Map {
	if g.len == 0 || keys.Len() == 0 {
		return g
	}
	c := &hashMap{len: g.len}
	g.resize(c, 0)
	c.RemoveAll(keys)
	if g.len == c.len {
		return g
	}
	c.frozen = g.frozen
	return c
}

func (g *hashMap) Type() dgo.Type {
	et := &exactMapType{value: g}
	et.ExactType = et
	return et
}

func (g *hashMap) Values() dgo.Array {
	return &array{slice: g.values(), frozen: g.frozen}
}

func (g *hashMap) values() []dgo.Value {
	ks := make([]dgo.Value, g.len)
	p := 0
	for e := g.first; e != nil; e = e.next {
		ks[p] = e.value
		p++
	}
	return ks
}

func (g *hashMap) resize(c *hashMap, capInc int) {
	tbl := g.table
	tl := tableSizeFor(len(tbl) + capInc)
	nt := make([]*hashNode, tl)
	c.table = nt
	tl--
	var prev *hashNode
	for oe := g.first; oe != nil; oe = oe.next {
		hk := tl & hash(oe.key.HashCode())
		ne := &hashNode{mapEntry: mapEntry{key: oe.key, value: oe.value}, hashNext: nt[hk]}
		if prev == nil {
			c.first = ne
		} else {
			prev.next = ne
			ne.prev = prev
		}
		nt[hk] = ne
		prev = ne
	}
	c.last = prev
}

func tableSizeFor(cap int) int {
	if cap < 1 {
		return 1
	}
	n := (uint)(cap - 1)
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	if n > maximumCapacity {
		return maximumCapacity
	}
	return int(n)
}

func frozenCopy(v dgo.Value) dgo.Value {
	if f, ok := v.(dgo.Freezable); ok {
		v = f.FrozenCopy()
	}
	return v
}

func hash(h int) int {
	return h ^ (h >> 16)
}

func mapTypeOne(args []interface{}) dgo.MapType {
	// min integer
	a0, ok := Value(args[0]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Map`, `Integer`, args, 0))
	}
	return newMapType(nil, nil, int(a0.GoInt()), math.MaxInt64)
}

func mapTypeTwo(args []interface{}) dgo.MapType {
	// key and value types or min and max integers
	switch a0 := Value(args[0]).(type) {
	case dgo.Type:
		a1, ok := Value(args[1]).(dgo.Type)
		if !ok {
			panic(illegalArgument(`Map`, `Type`, args, 1))
		}
		return newMapType(a0, a1, 0, math.MaxInt64)
	case dgo.Integer:
		a1, ok := Value(args[1]).(dgo.Integer)
		if !ok {
			panic(illegalArgument(`Map`, `Integer`, args, 1))
		}
		return newMapType(nil, nil, int(a0.GoInt()), int(a1.GoInt()))
	default:
		panic(illegalArgument(`Map`, `Type or Integer`, args, 0))
	}
}

func mapTypeThree(args []interface{}) dgo.MapType {
	// key and value types, and min integer
	a0, ok := Value(args[0]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Map`, `Type`, args, 0))
	}
	a1, ok := Value(args[1]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Map`, `Type`, args, 1))
	}
	a2, ok := Value(args[2]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Map`, `Integer`, args, 2))
	}
	return newMapType(a0, a1, int(a2.GoInt()), math.MaxInt64)
}

func mapTypeFour(args []interface{}) dgo.MapType {
	// key and value types, and min and max integers
	a0, ok := Value(args[0]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Map`, `Type`, args, 0))
	}
	a1, ok := Value(args[1]).(dgo.Type)
	if !ok {
		panic(illegalArgument(`Map`, `Type`, args, 1))
	}
	a2, ok := Value(args[2]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Map`, `Integer`, args, 2))
	}
	a3, ok := Value(args[3]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Map`, `Integer`, args, 3))
	}
	return newMapType(a0, a1, int(a2.GoInt()), int(a3.GoInt()))
}

// MapType returns a type that represents an Map value
func MapType(args []interface{}) dgo.MapType {
	switch len(args) {
	case 0:
		return DefaultMapType
	case 1:
		return mapTypeOne(args)
	case 2:
		return mapTypeTwo(args)
	case 3:
		return mapTypeThree(args)
	case 4:
		return mapTypeFour(args)
	default:
		panic(illegalArgumentCount(`MapType`, 0, 4, len(args)))
	}
}

func newMapType(kt, vt dgo.Type, min, max int) dgo.MapType {
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
	if kt == nil {
		kt = DefaultAnyType
	}
	if vt == nil {
		vt = DefaultAnyType
	}
	if min == 0 && max == math.MaxInt64 {
		// Unbounded
		if kt == DefaultAnyType && vt == DefaultAnyType {
			return DefaultMapType
		}
	}
	return &sizedMapType{keyType: kt, valueType: vt, min: min, max: max}
}

func (t *sizedMapType) Assignable(other dgo.Type) bool {
	return Assignable(nil, t, other)
}

func (t *sizedMapType) DeepAssignable(guard dgo.RecursionGuard, other dgo.Type) bool {
	if ot, ok := other.(dgo.MapType); ok {
		return t.min <= ot.Min() && ot.Max() <= t.max &&
			Assignable(guard, t.keyType, ot.KeyType()) && Assignable(guard, t.valueType, ot.ValueType())
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *sizedMapType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *sizedMapType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*sizedMapType); ok {
		return t.min == ot.min && t.max == ot.max && equals(seen, t.keyType, ot.keyType) && equals(seen, t.valueType, ot.valueType)
	}
	return false
}

func (t *sizedMapType) HashCode() int {
	return deepHashCode(nil, t)
}

func (t *sizedMapType) deepHashCode(seen []dgo.Value) int {
	h := int(dgo.TiMap)
	if t.min > 0 {
		h = h*31 + t.min
	}
	if t.max < math.MaxInt64 {
		h = h*31 + t.max
	}
	if DefaultAnyType != t.keyType {
		h = h*31 + deepHashCode(seen, t.keyType)
	}
	if DefaultAnyType != t.valueType {
		h = h*31 + deepHashCode(seen, t.keyType)
	}
	return h
}

func (t *sizedMapType) Instance(value interface{}) bool {
	return Instance(nil, t, value)
}

func (t *sizedMapType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	if ov, ok := value.(dgo.Map); ok {
		l := ov.Len()
		if t.min <= l && l <= t.max {
			kt := t.keyType
			vt := t.valueType
			if DefaultAnyType == kt {
				if DefaultAnyType == vt {
					return true
				}
				return ov.AllValues(func(v dgo.Value) bool { return Instance(guard, vt, v) })
			}
			if DefaultAnyType == vt {
				return ov.AllKeys(func(k dgo.Value) bool { return kt.Instance(k) })
			}
			return ov.All(func(e dgo.MapEntry) bool { return Instance(guard, kt, e.Key()) && Instance(guard, vt, e.Value()) })
		}
	}
	return false
}

func (t *sizedMapType) KeyType() dgo.Type {
	return t.keyType
}

func (t *sizedMapType) Max() int {
	return t.max
}

func (t *sizedMapType) Min() int {
	return t.min
}

func (t *sizedMapType) New(arg dgo.Value) dgo.Value {
	return newMap(t, arg)
}

func (t *sizedMapType) ReflectType() reflect.Type {
	return reflect.MapOf(t.KeyType().ReflectType(), t.ValueType().ReflectType())
}

func (t *sizedMapType) Resolve(ap dgo.AliasAdder) {
	kt := t.keyType
	vt := t.valueType
	t.keyType = DefaultAnyType
	t.valueType = DefaultAnyType
	kt = ap.Replace(kt).(dgo.Type)
	vt = ap.Replace(vt).(dgo.Type)
	t.keyType = kt
	t.valueType = vt
}

func (t *sizedMapType) String() string {
	return TypeString(t)
}

func (t *sizedMapType) Type() dgo.Type {
	return &metaType{t}
}

func (t *sizedMapType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMap
}

func (t *sizedMapType) ValueType() dgo.Type {
	return t.valueType
}

func (t *sizedMapType) Unbounded() bool {
	return t.min == 0 && t.max == math.MaxInt64
}

// DefaultMapType is the unconstrained Map type
const DefaultMapType = defaultMapType(0)

func (t defaultMapType) Assignable(other dgo.Type) bool {
	switch other.(type) {
	case defaultMapType, *sizedMapType, *exactMapType:
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t defaultMapType) Equals(other interface{}) bool {
	return t == other
}

func (t defaultMapType) HashCode() int {
	return int(dgo.TiMap)
}

func (t defaultMapType) Instance(value interface{}) bool {
	_, ok := value.(dgo.Map)
	return ok
}

func (t defaultMapType) KeyType() dgo.Type {
	return DefaultAnyType
}

func (t defaultMapType) Max() int {
	return math.MaxInt64
}

func (t defaultMapType) Min() int {
	return 0
}

func (t defaultMapType) New(arg dgo.Value) dgo.Value {
	return newMap(t, arg)
}

func (t defaultMapType) ReflectType() reflect.Type {
	return reflect.MapOf(reflectAnyType, reflectAnyType)
}

func (t defaultMapType) String() string {
	return TypeString(t)
}

func (t defaultMapType) Type() dgo.Type {
	return &metaType{t}
}

func (t defaultMapType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMap
}

func (t defaultMapType) ValueType() dgo.Type {
	return DefaultAnyType
}

func (t defaultMapType) Unbounded() bool {
	return true
}

func (t *exactMapType) Additional() bool {
	return false
}

func (t *exactMapType) Each(actor func(dgo.StructMapEntry)) {
	t.value.EachEntry(func(e dgo.MapEntry) {
		actor(&structEntry{mapEntry{e.Key().Type(), e.Value().Type()}, true})
	})
}

func (t *exactMapType) Generic() dgo.Type {
	kt := Generic(t.KeyType())
	vt := Generic(t.ValueType())
	if kt == DefaultAnyType && vt == DefaultAnyType {
		return DefaultMapType
	}
	return &sizedMapType{
		keyType:   kt,
		valueType: vt,
		min:       0,
		max:       math.MaxInt64}
}

func (t *exactMapType) Get(key interface{}) dgo.StructMapEntry {
	k := Value(key)
	if et, ok := k.(dgo.ExactType); ok {
		k = et.ExactValue()
	}
	if v := t.value.Get(k); v != nil {
		return &structEntry{mapEntry{k.Type(), v.Type()}, true}
	}
	return nil
}

func (t *exactMapType) KeyType() dgo.Type {
	l := t.value.Len()
	if l == 0 {
		return DefaultAnyType
	}
	return (*allOfValueType)(arrayFromIterator(l, t.value.EachKey))
}

func (t *exactMapType) Len() int {
	return t.value.Len()
}

func (t *exactMapType) Max() int {
	return t.value.Len()
}

func (t *exactMapType) Min() int {
	return t.value.Len()
}

func (t *exactMapType) New(arg dgo.Value) dgo.Value {
	return newMap(t, arg)
}

func (t *exactMapType) ReflectType() reflect.Type {
	return reflect.MapOf(t.KeyType().ReflectType(), t.ValueType().ReflectType())
}

func (t *exactMapType) Resolve(ap dgo.AliasAdder) {
	if ac, ok := t.value.(dgo.AliasContainer); ok {
		ac.Resolve(ap)
	}
}

func (t *exactMapType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiMapExact
}

func (t *exactMapType) Validate(keyLabel func(key dgo.Value) string, value interface{}) []error {
	return validate(t, keyLabel, value)
}

func (t *exactMapType) ValidateVerbose(value interface{}, out dgo.Indenter) bool {
	return validateVerbose(t, value, out)
}

func (t *exactMapType) ExactValue() dgo.Value {
	return t.value
}

func (t *exactMapType) ValueType() dgo.Type {
	l := t.value.Len()
	if l == 0 {
		return DefaultAnyType
	}
	return (*allOfValueType)(arrayFromIterator(l, t.value.EachValue))
}

func frozenMap(f string) error {
	return fmt.Errorf(`%s called on a frozen Map`, f)
}

func (t *exactMapType) Unbounded() bool {
	return false
}
