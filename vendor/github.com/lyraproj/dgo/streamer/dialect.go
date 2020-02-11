package streamer

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

// A Dialect determines how dgo values are serialized
type Dialect interface {
	// TypeKey returns the string that is used as a special hash key to denote a type. The default string is "__type"
	TypeKey() dgo.String

	// ValueKey returns the string that is used as a special hash key to denote a type. The default string is "__value"
	ValueKey() dgo.String

	// RefKey is the key used to signify the ordinal number of a previously serialized value. The
	// value is always an integer
	RefKey() dgo.String

	// AliasTypeName returns the string that denotes an alias. The default string is "alias"
	AliasTypeName() dgo.String

	// BinaryTypeName returns the string that denotes an alias. The default string is "binary"
	BinaryTypeName() dgo.String

	// MapTypeName returns the string that denotes an map that contains non-string keys. The default string is "map"
	MapTypeName() dgo.String

	// SensitiveTypeName returns the string that denotes a sensitive value. The default string is "sensitive"
	SensitiveTypeName() dgo.String

	// TimeTypeName returns the string that denotes a time. The default string is "time"
	TimeTypeName() dgo.String

	// ParseType parses the given type string and returns the resulting Type. The default parser will parse dgo syntax
	ParseType(aliasMap dgo.AliasMap, typeString dgo.String) dgo.Type
}

// DgoDialect returns the default dialect which is dgo
func DgoDialect() Dialect {
	return dgoDialectSingleton
}

type dgoDialect int

const dgoDialectSingleton = dgoDialect(0)

var typeKey = vf.String(`__type`)
var valueKey = vf.String(`__value`)
var refKey = vf.String(`__ref`)
var aliasType = vf.String(`alias`)
var binaryType = vf.String(`binary`)
var sensitiveType = vf.String(`sensitive`)
var mapType = vf.String(`map`)
var timeType = vf.String(`time`)

func (d dgoDialect) TypeKey() dgo.String {
	return typeKey
}

func (d dgoDialect) ValueKey() dgo.String {
	return valueKey
}

func (d dgoDialect) RefKey() dgo.String {
	return refKey
}

func (d dgoDialect) AliasTypeName() dgo.String {
	return aliasType
}

func (d dgoDialect) BinaryTypeName() dgo.String {
	return binaryType
}

func (d dgoDialect) MapTypeName() dgo.String {
	return mapType
}

func (d dgoDialect) SensitiveTypeName() dgo.String {
	return sensitiveType
}

func (d dgoDialect) TimeTypeName() dgo.String {
	return timeType
}

func (d dgoDialect) ParseType(aliasMap dgo.AliasMap, typeString dgo.String) dgo.Type {
	return typ.AsType(tf.ParseFile(aliasMap, ``, typeString.GoString()))
}
