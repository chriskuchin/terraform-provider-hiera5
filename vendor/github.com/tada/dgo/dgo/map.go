package dgo

type (
	// MapEntry is a key-value association in a Map
	MapEntry interface {
		Value
		Mutability

		Key() Value

		Value() Value
	}

	// StructMapEntry describes a MapEntry
	StructMapEntry interface {
		MapEntry

		Required() bool
	}

	// EntryActor performs some task on behalf of a caller
	EntryActor func(entry MapEntry)

	// EntryMapper maps produces the value of an entry to a new value
	EntryMapper func(entry MapEntry) interface{}

	// EntryPredicate returns true of false based on the given entry
	EntryPredicate func(entry MapEntry) bool

	// Keyed is the simples possible interface for a key store.
	Keyed interface {
		// Get returns the value for the given key. The method will return nil when the key is not found. A
		// vf.Nil is returned if the key is found but associated with nil.
		Get(key interface{}) Value
	}

	// Map represents an ordered set of key-value associations. The Map preserves the order by which the entries
	// were added. Associations retain their order even if their value change. When creating a Map from a go map
	// the associations will be sorted based on the natural order of the keys.
	Map interface {
		Iterable
		Keyed
		ReflectedValue
		Indentable

		// All returns true if the predicate returns true for all entries of this Map.
		All(predicate EntryPredicate) bool

		// AllKeys returns true if the predicate returns true for all keys of this Map.
		AllKeys(predicate Predicate) bool

		// AllValues returns true if the predicate returns true for all values of this Map.
		AllValues(predicate Predicate) bool

		// Any returns true if the predicate returns true for any entry of this Map.
		Any(actor EntryPredicate) bool

		// AnyKey returns true if the predicate returns true for any key of this Map.
		AnyKey(actor Predicate) bool

		// AnyValue returns true if the predicate returns true for any value of this Map.
		AnyValue(actor Predicate) bool

		// ContainsKey returns true if the map contains the give key
		ContainsKey(key interface{}) bool

		// Copy returns a copy of the Map. The copy is frozen or mutable depending on
		// the given argument. A request to create a frozen copy of an already frozen Map
		// is a no-op that returns the receiver
		//
		// If a frozen copy is requested from a non-frozen Map, then all non-frozen keys and
		// values will be copied and frozen recursively.
		//
		// A Copy of a map that contains back references to itself will result in a stack
		// overflow panic.
		Copy(frozen bool) Map

		// EachEntry calls the given actor with each entry of this Map
		EachEntry(actor EntryActor)

		// EachKey calls the given actor with each key of this Map
		EachKey(actor Consumer)

		// EachValue calls the given actor with each value of this Map
		EachValue(actor Consumer)

		// Find returns the first entry for which the entry predicate returns true
		Find(predicate EntryPredicate) MapEntry

		// Keys returns frozen snapshot of all the keys of this map
		Keys() Array

		// Map returns a new map with the same keys where each value has been replaced using the
		// given mapper function.
		Map(mapper EntryMapper) Map

		// Merge returns a Map where all associations from this and the given Map are merged. The associations of the
		// given map have priority.
		Merge(associations Map) Map

		// Put adds an association between the given key and value. The old value for the key or nil is returned. The
		// method will panic if the map is immutable
		Put(key, value interface{}) Value

		// PutAll adds all associations from the given Map, overwriting any that has the same key. It will panic if the
		// map is immutable.
		PutAll(associations Map)

		// Remove returns a Map that is guaranteed to have no value associated with the given key. The previous value
		// associated with the key or nil is returned. The method will panic if the map is immutable.
		Remove(key interface{}) Value

		// RemoveAll returns a Map that is guaranteed to have no values associated with any of the given keys. It will
		// panic if the map is immutable.
		RemoveAll(keys Array)

		// StringKeys returns true if this map's key type is assignable to String (i.e. if all keys are strings)
		StringKeys() bool

		// Values returns snapshot of all the values of this map.
		Values() Array

		// With creates a copy of this Map containing an association between the given key and value.
		With(key, value interface{}) Map

		// Without returns a Map that is guaranteed to have no value associated with the given key.
		Without(key interface{}) Map

		// WithoutAll returns a Map that is guaranteed to have no values associated with any of the given keys.
		WithoutAll(keys Array) Map
	}

	// A Struct represents a go struct as a Value.
	Struct interface {
		Map

		// GoStruct returns a pointer to the wrapped struct value.
		GoStruct() interface{}
	}

	// MapType is implemented by types representing implementations of the Map value
	MapType interface {
		sizedType

		// KeyType returns the type of the keys for instances of this type
		KeyType() Type

		// ValueType returns the type of the values for instances of this type
		ValueType() Type
	}

	// StructMapType represent a Map with explicitly defined typed entries.
	StructMapType interface {
		MapType

		// Additional returns true if the maps that is described by this type are allowed to
		// have additional entries.
		Additional() bool

		// EachEntryType iterates over each entry of the StructMapType
		EachEntryType(actor func(StructMapEntry))

		// GetEntryType returns the StructMapEntry that is identified with the given key
		GetEntryType(key interface{}) StructMapEntry

		// Len returns the number of StructEntrys in this StructMapType
		Len() int
	}

	// MapValidation provides methods for validate a Map against a MapType as a set of named parameters
	MapValidation interface {
		// Validate checks that the given value represents a Map which is an instance of this struct and returns a
		// possibly empty slice of errors explaining why that's not the case. Errors are generated if a required key
		// is missing, not recognized, or if it is of incorrect type.
		//
		// The keyLabel argument is an optional function that produces a suitable label for a key. If it is nil,
		// then a default function that produces the string "parameter '<key>'" will be used. The function
		// is called when errors are produced.
		//
		// An empty slice indicates a successful validation
		Validate(keyLabel func(key Value) string, value interface{}) []error

		// ValidateVerbose checks that the given value represents a Map which is an instance of this struct and returns
		// a boolean result. During validation, both successful and failing errors are verbosely explained on the given
		// Indenter.
		ValidateVerbose(value interface{}, out Indenter) bool
	}
)
