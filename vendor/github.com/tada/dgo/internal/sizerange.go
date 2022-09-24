package internal

import "github.com/tada/dgo/dgo"

// sizeRange is an inclusive positive integer range from zero to dgo.UnboundedSize
type sizeRange struct {
	min uint32
	max uint32
}

func (t *sizeRange) sizeRangeHash(base dgo.TypeIdentifier) dgo.Hash {
	h := dgo.Hash(base)
	if t.min > 0 {
		h = h*31 + dgo.Hash(t.min)
	}
	if t.max < dgo.UnboundedSize {
		h = h*31 + dgo.Hash(t.max)
	}
	return h
}

func (t *sizeRange) inRange(sz int) bool {
	return int(t.min) <= sz && sz <= int(t.max)
}

func (t *sizeRange) Max() int {
	return int(t.max)
}

func (t *sizeRange) Min() int {
	return int(t.min)
}

func (t *sizeRange) Unbounded() bool {
	return t.min == 0 && t.max == dgo.UnboundedSize
}
