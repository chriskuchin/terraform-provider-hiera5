package loader

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/vf"
)

type (
	multipleEntries struct {
		dgo.Map
	}

	mapLoader struct {
		name     string
		parentNs dgo.Loader
		entries  dgo.Map
	}

	loader struct {
		mapLoader
		lock       sync.RWMutex
		namespaces dgo.Map
		finder     dgo.Finder
		nsCreator  dgo.NsCreator
	}

	childLoader struct {
		dgo.Loader
		parent dgo.Loader
	}
)

// Multiple creates an value that holds multiple loader entries. The function is used by a finder that
// wishes to return more than the one entry that was requested. The requested element must be one of the
// entries.
func Multiple(m dgo.Map) dgo.Value {
	return multipleEntries{m}
}

func load(l dgo.Loader, name string) dgo.Value {
	parts := strings.Split(name, `/`)
	last := len(parts) - 1
	for i := 0; i < last; i++ {
		l = l.Namespace(parts[i])
		if l == nil {
			return nil
		}
	}
	return l.Get(parts[last])
}

// New returns a new Loader instance
func New(parentNs dgo.Loader, name string, entries dgo.Map, finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	if entries == nil {
		entries = vf.Map()
	}
	if finder == nil && nsCreator == nil {
		// Immutable map based loader
		return &mapLoader{parentNs: parentNs, name: name, entries: entries.FrozenCopy().(dgo.Map)}
	}

	var namespaces dgo.Map
	if nsCreator == nil {
		namespaces = vf.Map()
	} else {
		namespaces = vf.MutableMap()
	}
	return &loader{
		mapLoader:  mapLoader{parentNs: parentNs, name: name, entries: entries.Copy(finder == nil)},
		namespaces: namespaces,
		finder:     finder,
		nsCreator:  nsCreator}
}

// Type is the basic immutable loader dgo.Type
var Type = tf.NewNamed(`mapLoader`,
	func(arg dgo.Value) dgo.Value {
		l := &mapLoader{}
		l.init(arg.(dgo.Map))
		return l
	},
	func(v dgo.Value) dgo.Value {
		return v.(*mapLoader).initMap()
	},
	reflect.TypeOf(&mapLoader{}),
	reflect.TypeOf((*dgo.Loader)(nil)).Elem(),
	nil)

func (l *mapLoader) init(im dgo.Map) {
	l.name = im.Get(`name`).(dgo.String).GoString()
	l.entries = im.Get(`entries`).(dgo.Map)
}

func (l *mapLoader) initMap() dgo.Map {
	m := vf.MapWithCapacity(5)
	m.Put(`name`, l.name)
	m.Put(`entries`, l.entries)
	return m
}

func (l *mapLoader) String() string {
	return Type.ValueString(l)
}

func (l *mapLoader) Type() dgo.Type {
	return tf.ExactNamed(Type, l)
}

func (l *mapLoader) Equals(other interface{}) bool {
	if ov, ok := other.(*mapLoader); ok {
		return l.entries.Equals(ov.entries)
	}
	return false
}

func (l *mapLoader) HashCode() int {
	return l.entries.HashCode()
}

func (l *mapLoader) Get(key interface{}) dgo.Value {
	return l.entries.Get(key)
}

func (l *mapLoader) Load(name string) dgo.Value {
	return load(l, name)
}

func (l *mapLoader) AbsoluteName() string {
	an := `/` + l.name
	if l.parentNs != nil {
		if pn := l.parentNs.AbsoluteName(); pn != `/` {
			an = pn + an
		}
	}
	return an
}

func (l *mapLoader) Name() string {
	return l.name
}

func (l *mapLoader) Namespace(name string) dgo.Loader {
	if name == `` {
		return l
	}
	return nil
}

func (l *mapLoader) NewChild(finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func (l *mapLoader) ParentNamespace() dgo.Loader {
	return l.parentNs
}

// MutableType is the mutable loader dgo.Type
var MutableType = tf.NewNamed(`loader`,
	func(args dgo.Value) dgo.Value {
		l := &loader{}
		l.init(args.(dgo.Map))
		return l
	},
	func(v dgo.Value) dgo.Value {
		return v.(*loader).initMap()
	},
	reflect.TypeOf(&loader{}),
	reflect.TypeOf((*dgo.Loader)(nil)).Elem(),
	nil)

func (l *loader) init(im dgo.Map) {
	l.mapLoader.init(im)
	l.namespaces = im.Get(`namespaces`).(dgo.Map)
}

func (l *loader) initMap() dgo.Map {
	m := l.mapLoader.initMap()
	m.Put(`namespaces`, l.namespaces)
	return m
}

func (l *loader) add(key, value dgo.Value) dgo.Value {
	l.lock.Lock()
	defer l.lock.Unlock()

	addEntry := func(key, value dgo.Value) {
		if old := l.entries.Get(key); old == nil {
			l.entries.Put(key, value)
		} else if !old.Equals(value) {
			panic(fmt.Errorf(`attempt to override entry %q`, key))
		}
	}

	if m, ok := value.(multipleEntries); ok {
		value = m.Get(key)
		if value == nil {
			panic(fmt.Errorf(`map returned from finder doesn't contain original key %q`, key))
		}
		m.EachEntry(func(e dgo.MapEntry) { addEntry(e.Key(), e.Value()) })
	} else {
		addEntry(key, value)
	}
	return value
}

func (l *loader) Equals(other interface{}) bool {
	if ov, ok := other.(*loader); ok {
		return l.entries.Equals(ov.entries) && l.namespaces.Equals(ov.namespaces)
	}
	return false
}

func (l *loader) HashCode() int {
	return l.entries.HashCode()*31 + l.namespaces.HashCode()
}

func (l *loader) Get(ki interface{}) dgo.Value {
	key, ok := vf.Value(ki).(dgo.String)
	if !ok {
		return nil
	}
	l.lock.RLock()
	v := l.entries.Get(key)
	l.lock.RUnlock()
	if v == nil && l.finder != nil {
		v = vf.Value(l.finder(l, key.GoString()))
		v = l.add(key, v)
	}
	if vf.Nil == v {
		v = nil
	}
	return v
}

func (l *loader) Load(name string) dgo.Value {
	return load(l, name)
}

func (l *loader) Namespace(name string) dgo.Loader {
	if name == `` {
		return l
	}

	l.lock.RLock()
	ns, ok := l.namespaces.Get(name).(dgo.Loader)
	l.lock.RUnlock()

	if ok || l.nsCreator == nil {
		return ns
	}

	if ns = l.nsCreator(l, name); ns != nil {
		var old dgo.Value

		l.lock.Lock()
		if old = l.namespaces.Get(name); old == nil {
			l.namespaces.Put(name, ns)
		}
		l.lock.Unlock()

		if nil != old {
			// Either the nsCreator did something wrong that resulted in the creation of this
			// namespace or a another one has been created from another go routine.
			if !old.Equals(ns) {
				panic(fmt.Errorf(`namespace %q is already defined`, name))
			}

			// Get rid of the duplicate
			ns = old.(dgo.Loader)
		}
	}
	return ns
}

func (l *loader) NewChild(finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func (l *loader) String() string {
	return MutableType.ValueString(l)
}

func (l *loader) Type() dgo.Type {
	return tf.ExactNamed(MutableType, l)
}

// ChildType is the parented loader dgo.Type
var ChildType = tf.NewNamed(`childLoader`,
	func(args dgo.Value) dgo.Value {
		l := &childLoader{}
		l.init(args.(dgo.Map))
		return l
	},
	func(v dgo.Value) dgo.Value {
		return v.(*childLoader).initMap()
	},
	reflect.TypeOf(&loader{}),
	reflect.TypeOf((*dgo.Loader)(nil)).Elem(),
	nil)

func (l *childLoader) init(im dgo.Map) {
	l.Loader = im.Get(`loader`).(dgo.Loader)
	l.parent = im.Get(`parent`).(dgo.Loader)
}

func (l *childLoader) initMap() dgo.Map {
	m := vf.MapWithCapacity(2)
	m.Put(`loader`, l.Loader)
	m.Put(`parent`, l.parent)
	return m
}

func (l *childLoader) Equals(other interface{}) bool {
	if ov, ok := other.(*childLoader); ok {
		return l.Loader.Equals(ov.Loader) && l.parent.Equals(ov.parent)
	}
	return false
}

func (l *childLoader) Get(key interface{}) dgo.Value {
	v := l.parent.Get(key)
	if v == nil {
		v = l.Loader.Get(key)
	}
	return v
}

func (l *childLoader) HashCode() int {
	return l.Loader.HashCode()*31 + l.parent.HashCode()
}

func (l *childLoader) Load(name string) dgo.Value {
	return load(l, name)
}

func (l *childLoader) Namespace(name string) dgo.Loader {
	if name == `` {
		return l
	}
	pv := l.parent.Namespace(name)
	v := l.Loader.Namespace(name)
	switch {
	case v == nil:
		v = pv
	case pv == nil:
	default:
		v = &childLoader{Loader: v, parent: pv}
	}
	return v
}

func (l *childLoader) NewChild(finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	return loaderWithParent(l, finder, nsCreator)
}

func (l *childLoader) String() string {
	return ChildType.ValueString(l)
}

func (l *childLoader) Type() dgo.Type {
	return tf.ExactNamed(ChildType, l)
}

func loaderWithParent(parent dgo.Loader, finder dgo.Finder, nsCreator dgo.NsCreator) dgo.Loader {
	return &childLoader{Loader: New(parent.ParentNamespace(), parent.Name(), vf.Map(), finder, nsCreator), parent: parent}
}
