package internal

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

// deepEqual is implemented by values that need DeepEqual comparisons.
type deepEqual interface {
	deepEqual(seen []dgo.Value, other deepEqual) bool

	deepHashCode(seen []dgo.Value) int
}

type doubleSeen struct {
	aSeen   []dgo.Value
	bSeen   []dgo.Value
	seenInA bool
	seenInB bool
}

func (s *doubleSeen) Hit() bool {
	return s.seenInA && s.seenInB
}

func (s *doubleSeen) Append(a, b dgo.Value) dgo.RecursionGuard {
	if !s.seenInA {
		s.seenInA = util.RecursionHit(s.aSeen, a)
	}
	if !s.seenInB {
		s.seenInB = util.RecursionHit(s.bSeen, b)
	}
	c := *s
	if !c.seenInA {
		c.aSeen = append(c.aSeen, a)
	}
	if !c.seenInB {
		c.bSeen = append(c.bSeen, b)
	}
	return &c
}

func (s *doubleSeen) Swap() dgo.RecursionGuard {
	return &doubleSeen{seenInA: s.seenInB, seenInB: s.seenInA, aSeen: s.bSeen, bSeen: s.aSeen}
}

type deepCompare interface {
	deepCompare(seen []dgo.Value, other deepCompare) (int, bool)
}

// Assignable checks if b is assignable to a while guarding for endless recursion
func Assignable(guard dgo.RecursionGuard, a dgo.Type, b dgo.Type) bool {
	if a == b {
		return true
	}

	da, ok := a.(dgo.DeepAssignable)
	if !ok {
		return a.Assignable(b)
	}

	_, ok = b.(dgo.DeepAssignable)
	if ok {
		if guard == nil {
			guard = &doubleSeen{aSeen: []dgo.Value{a}, bSeen: []dgo.Value{b}}
		} else {
			guard = guard.Append(a, b)
			if guard.Hit() {
				return true
			}
		}
	}
	return da.DeepAssignable(guard, b)
}

// Instance checks if b is an instance of a to a while guarding for endless recursion
func Instance(guard dgo.RecursionGuard, a dgo.Type, b interface{}) bool {
	da, ok := a.(dgo.DeepInstance)
	if !ok {
		return a.Instance(b)
	}

	_, ok = b.(deepEqual) // only deepEqual implementations may be recursive
	if ok {
		bv := b.(dgo.Value)
		if guard == nil {
			guard = &doubleSeen{aSeen: []dgo.Value{a}, bSeen: []dgo.Value{bv}}
		} else {
			guard = guard.Append(a, bv)
			if guard.Hit() {
				return true
			}
		}
	}
	return da.DeepInstance(guard, b)
}

func deepHashCode(seen []dgo.Value, e dgo.Value) int {
	if de, ok := e.(deepEqual); ok {
		if util.RecursionHit(seen, e) {
			return 0
		}
		return de.deepHashCode(append(seen, e))
	}
	return e.HashCode()
}

// equals performs a deep equality comparison of a and b using the Value.Equals method. The given seen slice
// is used to prevent endless recursion. The rationale using a slice rather than a map for this is that the
// depth is typically very limited. The seen slice should be nil at the point where the comparison starts.
func equals(seen []dgo.Value, a dgo.Value, b interface{}) bool {
	if a == b {
		return true
	}
	if nil == a || nil == b {
		return false
	}
	da, ok := a.(deepEqual)
	if !ok {
		return a.Equals(b)
	}
	if util.RecursionHit(seen, a) {
		// Recursion, so assume true
		return true
	}
	db, ok := Value(b).(deepEqual)
	if !ok {
		// Must be false since only one implements deepEqual
		return false
	}
	return da.deepEqual(append(seen, a), db)
}

func sliceEquals(seen []dgo.Value, a, b []dgo.Value) bool {
	l := len(a)
	if l != len(b) {
		return false
	}
	for i := 0; i < l; i++ {
		if !equals(seen, a[i], b[i]) {
			return false
		}
	}
	return true
}

func compare(seen []dgo.Value, a dgo.Value, b dgo.Value) (int, bool) {
	if a == b {
		return 0, true
	}
	if a == Nil {
		return -1, true
	}
	if b == Nil {
		return 1, true
	}

	da, ok := a.(deepCompare)
	if !ok {
		return a.(dgo.Comparable).CompareTo(b)
	}

	db, ok := b.(deepCompare)
	if !ok {
		// Calling back to a.CompareTo at this point would cause endless recursion but it should be safe to
		// assume that a deepCompare cannot be compared to a non deepCompare.
		return 0, false
	}

	if util.RecursionHit(seen, a) {
		// Recursion, so assume equal
		return 0, true
	}
	return da.deepCompare(append(seen, a), db)
}
