package vf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// Sensitive creates a new Sensitive that wraps the given value
func Sensitive(v interface{}) dgo.Sensitive {
	return internal.Sensitive(v)
}
