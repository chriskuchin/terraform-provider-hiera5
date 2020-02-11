package hiera

import (
	"github.com/lyraproj/dgo/dgo"
)

type (
	// DataDig is a Hiera 'data_dig' function looks up a value by a key consisting of several segments.
	// The segments are either strings or ints. No other types of segments are allowed.
	DataDig func(ic ProviderContext, key dgo.Array) dgo.Value

	// DataHash is a Hiera 'data_hash' function returns a Map that Hiera can use as the source for
	// lookups.
	DataHash func(ic ProviderContext) dgo.Map

	// LookupKey is a Hiera 'lookup_key' function returns the value that corresponds to the given key.
	LookupKey func(ic ProviderContext, key string) dgo.Value
)
