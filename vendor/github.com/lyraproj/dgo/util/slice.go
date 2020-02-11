package util

import "github.com/lyraproj/dgo/dgo"

// SliceCopy copies the given value slice to a new slice
func SliceCopy(s []dgo.Value) []dgo.Value {
	c := make([]dgo.Value, len(s))
	copy(c, s)
	return c
}

// RecursionHit returns true if the given this dgo.Value exists within the given seen slice. The comparison
// uses == and hence, compares for identity rather than equality. The function is intended to be used when
// detecting endless recursion in self recursion structures. Hence its name.
func RecursionHit(seen []dgo.Value, this dgo.Value) bool {
	for i := range seen {
		if this == seen[i] {
			return true
		}
	}
	return false
}
