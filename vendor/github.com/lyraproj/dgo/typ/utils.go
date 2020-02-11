package typ

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// ExactValue returns the "exact value" that a value represents. If the given value is a dgo.ExactType, then the value
// that it represents is the exact value. For all other cases, the exact value is the value itself.
func ExactValue(value dgo.Value) dgo.Value {
	return internal.ExactValue(value)
}

// AsType returns the value as a type. If the value already is a type, it is returned. Otherwise the
// exact type of the value is returned.
func AsType(value dgo.Value) dgo.Type {
	if tp, ok := value.(dgo.Type); ok {
		return tp
	}
	return value.Type()
}

// Generic returns the generic form of the given type. All non exact types are considered generic
// and will be returned directly. Exact types will loose information about what instance they represent
// and also range and size information. Nested types will return a generic version of the contained
// types as well.
func Generic(t dgo.Type) dgo.Type {
	return internal.Generic(t)
}
