package internal

import (
	"sync"

	"github.com/tada/catch"
	"github.com/tada/dgo/dgo"
)

type (
	alias struct {
		dgo.StringType
	}

	aliasMap struct {
		typeNames  hashMap
		namedTypes hashMap
	}

	aliasAdder struct {
		namedTypes hashMap
		backingMap dgo.AliasMap
	}

	dType        = dgo.Type // To avoid collision with method named Type
	deferredCall struct {
		dType
		args dgo.Arguments
	}
)

// BuiltInAliases returns a frozen AliasMap containing the predefined aliases
func BuiltInAliases() dgo.AliasMap {
	return builtinAliases
}

var defaultLock = sync.Mutex{}

// AddDefaultAliases adds the new aliases to the default alias map by passing an AliasAdder to the function
// The function is safe from a concurrency perspective.
func AddDefaultAliases(adder func(adder dgo.AliasAdder)) {
	AddAliases(&defaultAliases, &defaultLock, adder)
}

// AddAliases will call the given adder function, and if entries were added, lock the appointed Locker, create
// a copy of the appointed AliasMap, add the entries to that copy, swap the appointed AliasMap for the copy,
// and finally release the lock.
//
// No Locker is locked and no swap will take place if the adder function doesn't add anything.
func AddAliases(mapToReplace *dgo.AliasMap, lock sync.Locker, adder func(adder dgo.AliasAdder)) {
	am := &aliasAdder{backingMap: *mapToReplace}
	adder(am)
	if am.namedTypes.Len() > 0 {
		lock.Lock()
		defer lock.Unlock()
		*mapToReplace = (*mapToReplace).(*aliasMap).update(am)
	}
}

// DefaultAliases returns the frozen default dgo.AliasMap
func DefaultAliases() dgo.AliasMap {
	return defaultAliases
}

// ResetDefaultAliases will reset the AliasMap returned by the DefaultAliases() method to the BuiltInAliases()
// and thereby throw away any changes made to the DefaultAliases map.
//
// This method is intended for testing purposes only
func ResetDefaultAliases() {
	defaultAliases = builtinAliases
}

// NewCall creates a special interim type that represents a call during parsing, and then nowhere else.
func NewCall(s dgo.Type, args dgo.Arguments) dgo.Type {
	return &deferredCall{s, args}
}

// NewAlias creates a special interim type that represents a type alias used during parsing, and then nowhere else.
func NewAlias(s dgo.String) dgo.Alias {
	return &alias{s.Type().(dgo.StringType)}
}

func (a *alias) Frozen() bool {
	return false
}

func (a *alias) FrozenCopy() dgo.Value {
	panic(catch.Error(`attempt to freeze unresolved alias '%s'`, a.Reference()))
}

func (a *alias) ThawedCopy() dgo.Value {
	return a
}

func (a *alias) Reference() dgo.String {
	return a.StringType.(dgo.String)
}

func (a *alias) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiAlias
}

var builtinAliases dgo.AliasMap
var defaultAliases dgo.AliasMap

func init() {
	m := (&aliasMap{}).Collect(func(b dgo.AliasAdder) {
		dataAlias := NewAlias(String(`data`))
		data := AnyOfType([]interface{}{
			DefaultStringType,
			DefaultIntegerType,
			DefaultFloatType,
			DefaultBooleanType,
			Nil,
			ArrayType([]interface{}{dataAlias}),
			MapType([]interface{}{DefaultStringType, dataAlias})})
		b.Add(data, dataAlias.Reference())

		richDataAlias := NewAlias(String(`richdata`))
		richData := AnyOfType([]interface{}{
			DefaultStringType,
			DefaultIntegerType,
			DefaultFloatType,
			DefaultBooleanType,
			DefaultBinaryType,
			DefaultMetaType,
			DefaultRegexpType,
			DefaultSensitiveType,
			DefaultTimeType,
			Nil,
			ArrayType([]interface{}{richDataAlias}),
			MapType([]interface{}{AnyOfType([]interface{}{DefaultStringType, DefaultIntegerType, DefaultFloatType}), richDataAlias})})
		b.Add(richData, richDataAlias.Reference())
	})
	builtinAliases = m
	defaultAliases = m
}

func (a *aliasMap) Collect(adder func(dgo.AliasAdder)) dgo.AliasMap {
	am := &aliasAdder{backingMap: a}
	adder(am)
	if am.namedTypes.Len() > 0 {
		return a.update(am)
	}
	return a
}

func (a *aliasMap) update(am *aliasAdder) dgo.AliasMap {
	// Create a new alias map with enough room to fit the new entries
	c := &aliasMap{}
	ns := am.namedTypes
	na := ns.Len()
	a.namedTypes.resize(&c.namedTypes, na)
	a.typeNames.resize(&c.typeNames, na)

	// Resolve the added entries
	rs := ns.Map(func(e dgo.MapEntry) interface{} { return am.Replace(e.Value()) })

	// Add entries to the new alias map
	rs.EachEntry(func(e dgo.MapEntry) {
		name := e.Key().(dgo.String)
		t := e.Value().(dgo.Type)
		if mt, mutable := t.(dgo.Mutability); mutable {
			t = mt.FrozenCopy().(dgo.Type)
		}
		c.namedTypes.Put(name, t)
		c.typeNames.Put(t, name)
	})
	return c
}

// GetName returns the name for the given type or nil if the type isn't found
func (a *aliasMap) GetName(t dgo.Type) dgo.String {
	if v := a.typeNames.Get(t); v != nil {
		return v.(dgo.String)
	}
	return nil
}

// GetType returns the type with the given name or nil if the type isn't found
func (a *aliasMap) GetType(n dgo.String) dgo.Type {
	if v := a.namedTypes.Get(n); v != nil {
		return v.(dgo.Type)
	}
	return nil
}

func (a *aliasAdder) Add(t dgo.Type, name dgo.String) {
	a.namedTypes.Put(name, t)
}

func (a *aliasAdder) GetType(n dgo.String) dgo.Type {
	if t := a.namedTypes.Get(n); t != nil {
		return t.(dgo.Type)
	}
	return a.backingMap.GetType(n)
}

func (a *aliasAdder) Replace(t dgo.Value) dgo.Value {
	switch t := t.(type) {
	case *deferredCall:
		return New(t.dType, t.args)
	case dgo.Alias:
		if ra := a.GetType(t.Reference()); ra != nil {
			return ra
		}
		panic(catch.Error(`reference to unresolved type '%s'`, t.Reference()))
	case dgo.AliasContainer:
		t.Resolve(a)
	}
	return t
}
