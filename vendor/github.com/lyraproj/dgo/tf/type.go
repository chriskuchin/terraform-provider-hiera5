package tf

import (
	"reflect"
	"sync"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
	"github.com/lyraproj/dgo/parser"

	// ensure that stringer package is initialized prior to using this package
	_ "github.com/lyraproj/dgo/stringer"
)

// FromReflected returns the dgo.Type that represents the given reflected type
func FromReflected(vt reflect.Type) dgo.Type {
	return internal.TypeFromReflected(vt)
}

// ParseType parses the given content into a dgo.Type.
func ParseType(content string) dgo.Type {
	return internal.AsType(parser.Parse(content))
}

// Parse parses the given content into a dgo.Value.
func Parse(content string) dgo.Value {
	return internal.ExactValue(parser.Parse(content))
}

// ParseFile parses the given content into a dgo.Type. The filename is used in error messages.
//
// The alias map is optional. If given, the parser will recognize the type aliases provided in the map
// and also add any new aliases declared within the parsed content to that map.
func ParseFile(aliasMap dgo.AliasAdder, fileName, content string) dgo.Value {
	return parser.ParseFile(aliasMap, fileName, content)
}

// AddDefaultAliases adds the new aliases to the default alias map by passing an AliasAdder to the function
// The function is safe from a concurrency perspective.
func AddDefaultAliases(adderFunc func(aliasAdder dgo.AliasAdder)) {
	internal.AddDefaultAliases(adderFunc)
}

// AddAliases will call the given adder function, and if entries were added, lock the appointed Locker, create
// a copy of the appointed AliasMap, add the entries to that copy, swap the appointed AliasMap for the copy,
// and finally release the lock.
//
// No Locker is locked and no swap will take place if the adder function doesn't add anything.
func AddAliases(mapToReplace *dgo.AliasMap, lock sync.Locker, adder func(adder dgo.AliasAdder)) {
	internal.AddAliases(mapToReplace, lock, adder)
}

// BuiltInAliases returns the frozen built-in dgo.AliasMap
func BuiltInAliases() dgo.AliasMap {
	return internal.BuiltInAliases()
}

// DefaultAliases returns the default dgo.AliasMap
func DefaultAliases() dgo.AliasMap {
	return internal.DefaultAliases()
}

// Meta creates the meta type for the given type
func Meta(t dgo.Type) dgo.Meta {
	return internal.MetaType(t)
}
