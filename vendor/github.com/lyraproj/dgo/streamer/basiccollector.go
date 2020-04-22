package streamer

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/vf"
)

// A BasicCollector is an extendable basic implementation of the Consumer interface
type BasicCollector struct {
	// Values is an array of all values that are added to the BasicCollector. When adding
	// a reference, the reference is considered to be an index in this array.
	Values dgo.Array

	// The Stack of values that is used when adding nested constructs (arrays and maps)
	Stack []dgo.Array
}

// NewCollector returns a new BasicCollector instance
func NewCollector() Collector {
	c := &BasicCollector{}
	c.Init()
	return c
}

// Init initializes the internal stack and reference storage
func (c *BasicCollector) Init() {
	c.Values = vf.MutableValues()
	c.Stack = make([]dgo.Array, 1, 8)
	c.Stack[0] = vf.MutableValues()
}

// AddArray initializes and adds a new array and then calls the function with is supposed to
// add the elements.
func (c *BasicCollector) AddArray(cap int, doer dgo.Doer) {
	a := vf.ArrayWithCapacity(cap)
	c.Add(a)
	top := len(c.Stack)
	c.Stack = append(c.Stack, a)
	doer()
	c.Stack = c.Stack[0:top]
}

// AddMap initializes and adds a new map and then calls the function with is supposed to
// add an even number of elements as a sequence of key, value, [key, value, ...]
func (c *BasicCollector) AddMap(cap int, doer dgo.Doer) {
	h := vf.MapWithCapacity(cap)
	c.Add(h)
	a := vf.ArrayWithCapacity(cap * 2)
	top := len(c.Stack)
	c.Stack = append(c.Stack, a)
	doer()
	c.Stack = c.Stack[0:top]
	h.PutAll(a.ToMap())
}

// Add adds a new value
func (c *BasicCollector) Add(element dgo.Value) {
	c.StackTop().Add(element)
	c.Values.Add(element)
}

// AddRef adds the nth value of the values that has been added once again.
func (c *BasicCollector) AddRef(ref int) {
	c.StackTop().Add(c.Values.Get(ref))
}

// CanDoBinary returns true
func (c *BasicCollector) CanDoBinary() bool {
	return true
}

// CanDoTime returns true
func (c *BasicCollector) CanDoTime() bool {
	return true
}

// CanDoComplexKeys returns true
func (c *BasicCollector) CanDoComplexKeys() bool {
	return true
}

// StringDedupThreshold returns 0
func (c *BasicCollector) StringDedupThreshold() int {
	return 0
}

// Value returns the last value added to this collector
func (c *BasicCollector) Value() dgo.Value {
	return c.Stack[0].Get(0)
}

// StackTop returns the Array at the top of the collector stack.
func (c *BasicCollector) StackTop() dgo.Array {
	return c.Stack[len(c.Stack)-1]
}

// PeekLast returns the last added value from
func (c *BasicCollector) PeekLast() dgo.Value {
	a := c.StackTop()
	return a.Get(a.Len() - 1)
}

// ReplaceLast replaces the last added value with the given value
func (c *BasicCollector) ReplaceLast(v dgo.Value) {
	a := c.StackTop()
	a.Set(a.Len()-1, v)
}
