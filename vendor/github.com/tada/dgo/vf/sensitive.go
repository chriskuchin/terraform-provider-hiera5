package vf

import (
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/internal"
)

// Sensitive creates a new Sensitive that wraps the given value
func Sensitive(v interface{}) dgo.Sensitive {
	return internal.Sensitive(v)
}
