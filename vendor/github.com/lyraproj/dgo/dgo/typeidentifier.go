package dgo

import "fmt"

// TypeIdentifier is a unique identifier for each type known to the system. The order of the TypeIdentifier
// determines the sort order for elements that are not comparable
type TypeIdentifier int

const (
	// TiAlias is the type identifier for the alias type reference
	TiAlias = TypeIdentifier(iota)

	// TiAllOf is the type identifier for the AllOf type
	TiAllOf

	// TiAllOfValue is the type identifier for the AllOf type that uses the type of its contained values
	TiAllOfValue

	// TiAny is the type identifier for the Any type
	TiAny

	// TiAnyOf is the type identifier for the AnyOf type
	TiAnyOf

	// TiArray is the type identifier for the Array type
	TiArray

	// TiBinary is the type identifier for the Binary type
	TiBinary

	// TiBoolean is the type identifier for the Boolean type
	TiBoolean

	// TiCiString is the type identifier for the case insensitive String type
	TiCiString

	// TiDgoString is the type identifier for for the DgoString type
	TiDgoString

	// TiError is the type identifier for for the Error type
	TiError

	// TiFloat is the type identifier for the Float type
	TiFloat

	// TiFloatRange is the type identifier for the Float range type
	TiFloatRange

	// TiFunction is the type identifier for for the Function type
	TiFunction

	// TiInteger is the type identifier for the Integer type
	TiInteger

	// TiIntegerRange is the type identifier for the Integer range type
	TiIntegerRange

	// TiMap is the type identifier for the Map type
	TiMap

	// TiMeta is the type identifier for the Meta type
	TiMeta

	// TiNamed is the type identifier for for named types
	TiNamed

	// TiNative is the type identifier for the Native type
	TiNative

	// TiNot is the type identifier for the Not type
	TiNot

	// TiOneOf is the type identifier for the OneOf type
	TiOneOf

	// TiRegexp is the type identifier for the Regexp type
	TiRegexp

	// TiSensitive is the type identifier for for the Sensitive type
	TiSensitive

	// TiString is the type identifier for the String type
	TiString

	// TiStringPattern is the type identifier for the String pattern type
	TiStringPattern

	// TiStringSized is the type identifier for the size constrained String type
	TiStringSized

	// TiStruct is the type identifier for the Struct type
	TiStruct

	// TiTime is the type identifier for for the Time type
	TiTime

	// TiTuple is the type identifier for the Tuple type
	TiTuple

	// exactStart denotes the index of where the range of exact types start. All
	// exact types must be added below this entry
	exactStart

	// TiArrayExact is the type identifier for the exact Array type
	TiArrayExact

	// TiBinaryExact is the type identifier for the exact Binary type
	TiBinaryExact

	// TiBooleanExact is the type identifier for the exact Boolean type
	TiBooleanExact

	// TiErrorExact is the type identifier for for the exact Error type
	TiErrorExact

	// TiFloatExact is the type identifier for the exact Float type
	TiFloatExact

	// TiFunctionExact is the type identifier for for the exact Function type
	TiFunctionExact

	// TiIntegerExact is the type identifier for the exact Integer type
	TiIntegerExact

	// TiMapExact is the type identifier for exact Map type
	TiMapExact

	// TiMapEntryExact is the type identifier the map entry type of the exact Map type
	TiMapEntryExact

	// TiNamedExact is the type identifier for for exact Named types
	TiNamedExact

	// TiNil is the type identifier for the Nil type
	TiNil

	// TiRegexpExact is the type identifier for the exact Regexp type
	TiRegexpExact

	// TiStringExact is the type identifier for the exact String type
	TiStringExact

	// TiTimeExact is the type identifier for the exact Time type
	TiTimeExact
)

var tiLabels = map[TypeIdentifier]string{
	TiNil:           `nil`,
	TiAny:           `any`,
	TiMeta:          `type`,
	TiBoolean:       `bool`,
	TiBooleanExact:  `bool`,
	TiInteger:       `int`,
	TiIntegerExact:  `int`,
	TiIntegerRange:  `int range`,
	TiFloat:         `float`,
	TiFloatExact:    `float`,
	TiFloatRange:    `float range`,
	TiBinary:        `binary`,
	TiBinaryExact:   `binary`,
	TiString:        `string`,
	TiStringExact:   `string`,
	TiStringSized:   `string`,
	TiStringPattern: `pattern`,
	TiCiString:      `string`,
	TiRegexp:        `regexp`,
	TiRegexpExact:   `regexp`,
	TiTime:          `time`,
	TiTimeExact:     `time`,
	TiNative:        `native`,
	TiArray:         `slice`,
	TiArrayExact:    `slice`,
	TiTuple:         `tuple`,
	TiMap:           `map`,
	TiMapExact:      `map`,
	TiMapEntryExact: `map entry`,
	TiStruct:        `struct`,
	TiNot:           `not`,
	TiAllOf:         `all of`,
	TiAllOfValue:    `all of`,
	TiAnyOf:         `any of`,
	TiOneOf:         `one of`,
	TiError:         `error`,
	TiErrorExact:    `error`,
	TiDgoString:     `dgo`,
	TiSensitive:     `sensitive`,
	TiFunction:      `function`,
	TiFunctionExact: `function`,
	TiNamed:         `named`,
}

func (ti TypeIdentifier) String() string {
	if s, ok := tiLabels[ti]; ok {
		return s
	}
	panic(fmt.Errorf("unhandled TypeIdentifier %d", ti))
}

// IsExact returns true if the given type represents an exact value.
func IsExact(value Type) bool {
	return value.TypeIdentifier() > exactStart
}
