package vf

import (
	"reflect"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
)

// Map creates an immutable dgo.Map from the given slice which must have 0, 1, or an
// even number of arguments.
//
// Zero arguments: the empty map is returned.
//
// One argument: must be a go map, a go struct, or an Array with an even number of elements.
//
// An even number of arguments: will be considered a flat list of key, value [, key, value, ... ]
func Map(m ...interface{}) dgo.Map {
	return internal.Map(m)
}

// MutableMap creates an empty dgo.Map. The map can be optionally constrained
// by the given type which can be nil, the zero value of a go map, or a dgo.MapType
func MutableMap(m ...interface{}) dgo.Map {
	return internal.MutableMap(m)
}

// MapWithCapacity creates an empty dgo.Map suitable to hold a given number of entries.
func MapWithCapacity(capacity int) dgo.Map {
	return internal.MapWithCapacity(capacity)
}

// FromReflectedMap creates a Map from a reflected map. If frozen is true, the created Map will be
// immutable and the type will reflect exactly that map and nothing else. If frozen is false, the
// created Map will be mutable and its type will be derived from the reflected map.
func FromReflectedMap(rm reflect.Value, frozen bool) dgo.Map {
	return internal.FromReflectedMap(rm, frozen).(dgo.Map)
}
