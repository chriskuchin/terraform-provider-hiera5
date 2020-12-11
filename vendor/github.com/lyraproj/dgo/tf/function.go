package tf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Function returns a new dgo.FunctionType with the given argument and return value
// types.
func Function(args dgo.TupleType, returns dgo.TupleType) dgo.FunctionType {
	return internal.FunctionType(args, returns)
}
