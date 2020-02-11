package streamer

import (
	"github.com/lyraproj/dgo/dgo"
)

// A Consumer is used by a Data streaming mechanism that maintains a reference index
// which is increased by one for each value that it streams. The reference index
// originates from zero.
type Consumer interface {
	// CanDoBinary returns true if the value can handle binary efficiently. This tells
	// the Serializer to pass dgo.Binary verbatim to Add
	CanDoBinary() bool

	// CanDoTime returns true if the value can handle timestamp efficiently. This tells
	// the Serializer to pass dgo.Time verbatim to Add
	CanDoTime() bool

	// CanDoComplexKeys returns true if complex values can be used as keys. If this
	// method returns false, all keys must be strings
	CanDoComplexKeys() bool

	// StringDedupThreshold returns the preferred threshold for dedup of strings. Strings
	// shorter than this threshold will not be subjected to de-duplication.
	StringDedupThreshold() int

	// AddArray starts a new array, calls the doer function, and then ends the Array.
	//
	// The callers reference index is increased by one.
	AddArray(len int, doer dgo.Doer)

	// AddMap starts a new map, calls the doer function, and then ends the map.
	//
	// The callers reference index is increased by one.
	AddMap(len int, doer dgo.Doer)

	// Add adds the next value.
	//
	// Calls following a StartArray will add elements to the Array
	//
	// Calls following a StartHash will first add a key, then a value. This
	// repeats until End or StartArray is called.
	//
	// The callers reference index is increased by one.
	Add(element dgo.Value)

	// AddRef adds a reference to a previously added afterElement, hash, or array.
	AddRef(ref int)
}

// A Collector receives streaming events and produces an Value
type Collector interface {
	Consumer

	Value() dgo.Value
}
