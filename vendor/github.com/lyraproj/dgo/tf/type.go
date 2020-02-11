package tf

import (
	"reflect"

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
func ParseFile(aliasMap dgo.AliasMap, fileName, content string) dgo.Value {
	return parser.ParseFile(aliasMap, fileName, content)
}

// NewAliasMap creates a new dgo.Alias map to be used as a scope when parsing types
func NewAliasMap() dgo.AliasMap {
	return internal.NewAliasMap()
}

// Meta creates the meta type for the given type
func Meta(t dgo.Type) dgo.Meta {
	return internal.MetaType(t)
}
