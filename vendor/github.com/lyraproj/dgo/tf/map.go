package tf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Map returns a type that represents an Map value
func Map(args ...interface{}) dgo.MapType {
	return internal.MapType(args)
}

// StructMapEntry returns a new StructMapEntry initiated with the given parameters
func StructMapEntry(key interface{}, value interface{}, required bool) dgo.StructMapEntry {
	return internal.StructMapEntry(key, value, required)
}

// StructMap returns a new StructMapType type built from the given MapEntryTypes. If
// additional is true, the struct will allow additional unconstrained entries
func StructMap(additional bool, entries ...dgo.StructMapEntry) dgo.StructMapType {
	return internal.StructMapType(additional, entries)
}

// StructMapFromMap returns a new type built from a map[string](dgo|type|{type:dgo|type,required?:bool,...})
func StructMapFromMap(additional bool, entries dgo.Map) dgo.StructMapType {
	return internal.StructMapTypeFromMap(additional, entries)
}
