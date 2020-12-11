package tf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Sensitive returns a Sensitive dgo.Type that wraps the given dgo.Type
func Sensitive(args ...interface{}) dgo.Type {
	return internal.SensitiveType(args)
}
