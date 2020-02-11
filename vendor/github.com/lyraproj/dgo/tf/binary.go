package tf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Binary returns a new dgo.BinaryType. It can be called with two optional integer arguments denoting
// the min and max length of the binary. If only one integer is given, it represents the min length.
func Binary(args ...interface{}) dgo.BinaryType {
	return internal.BinaryType(args...)
}
