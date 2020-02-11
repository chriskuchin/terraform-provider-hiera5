package internal

import (
	"fmt"

	"github.com/lyraproj/dgo/dgo"
)

type (
	alias struct {
		dgo.StringType
	}

	aliasMap struct {
		typeNames  hashMap
		namedTypes hashMap
	}

	dType        = dgo.Type // To avoid collision with method named Type
	deferredCall struct {
		dType
		args dgo.Arguments
	}
)

// NewCall creates a special interim type that represents a call during parsing, and then nowhere else.
func NewCall(s dgo.Type, args dgo.Arguments) dgo.Type {
	return &deferredCall{s, args}
}

// NewAlias creates a special interim type that represents a type alias used during parsing, and then nowhere else.
func NewAlias(s dgo.String) dgo.Alias {
	return &alias{s.Type().(dgo.StringType)}
}

func (a *alias) Reference() dgo.String {
	return a.StringType.(dgo.ExactType).ExactValue().(dgo.String)
}

func (a *alias) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiAlias
}

var builtinAliases *aliasMap

func init() {
	m := &aliasMap{}
	dataAlias := NewAlias(String(`data`))
	data := AnyOfType([]interface{}{
		DefaultStringType,
		DefaultIntegerType,
		DefaultFloatType,
		DefaultBooleanType,
		DefaultNilType,
		ArrayType([]interface{}{dataAlias}),
		MapType([]interface{}{DefaultStringType, dataAlias})})
	m.Add(data, dataAlias.Reference())

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
		DefaultNilType,
		ArrayType([]interface{}{richDataAlias}),
		MapType([]interface{}{AnyOfType([]interface{}{DefaultStringType, DefaultIntegerType, DefaultFloatType}), richDataAlias})})
	m.Add(richData, richDataAlias.Reference())

	data.(dgo.AliasContainer).Resolve(m)
	richData.(dgo.AliasContainer).Resolve(m)
	builtinAliases = m
}

// NewAliasMap creates a new dgo.AliasMap to be used as a scope when parsing types
func NewAliasMap() dgo.AliasMap {
	m := &aliasMap{}
	builtinAliases.namedTypes.resize(&m.namedTypes, 0)
	builtinAliases.typeNames.resize(&m.typeNames, 0)
	return m
}

func (a *aliasMap) GetName(t dgo.Type) dgo.String {
	if v := a.typeNames.Get(t); v != nil {
		return v.(dgo.String)
	}
	return nil
}

func (a *aliasMap) GetType(n dgo.String) dgo.Type {
	if v := a.namedTypes.Get(n); v != nil {
		return v.(dgo.Type)
	}
	return nil
}

func (a *aliasMap) Add(t dgo.Type, name dgo.String) {
	a.typeNames.Put(t, name)
	a.namedTypes.Put(name, t)
}

func (a *aliasMap) Replace(t dgo.Value) dgo.Value {
	switch t := t.(type) {
	case *deferredCall:
		return New(t.dType, t.args)
	case dgo.Alias:
		if ra := a.GetType(t.Reference()); ra != nil {
			return ra
		}
		panic(fmt.Errorf(`reference to unresolved type '%s'`, t.Reference()))
	case dgo.AliasContainer:
		t.Resolve(a)
	}
	return t
}
