package dgo

type (
	// Finder is a function that is called when a loader is asked to deliver a value that it doesn't
	// have. It's the finders responsibility to obtain the value. The finder must return nil when the
	// value could not be found.
	//
	// A finder can return a Map where the given key is one of multiple keys in case the finder wants
	// to make more elements known to the loader. The given key is required in such a map.
	//
	// The key sent to the finder will never be a multi part name. It will reflect a single name that is
	// relative to the loader that the Finder is configured for.
	Finder func(l Loader, key string) interface{}

	// NsCreator is called when a namespace is requested that does not exist. It is the NsCreator's
	// responsibility to create the Loader that represents the new namespace. The creator must return
	// nil when no such namespace can be created.
	//
	// The name sent to the creator will never be a multi part name. It will reflect a single name that is
	// relative to the loader that the Finder is configured for.
	NsCreator func(l Loader, name string) Loader

	// A Loader loads named values on demand. Loaders for nested namespaces can be obtained using the Namespace
	// method.
	//
	// Implementors of Loader must ensure that all methods are safe to use from concurrent go routines.
	//
	// The Loader implements the Keyed interface. The Get method is different from the Load method in
	// that the name must contain no '/' characters since they separate the name from its namespace.
	// A load involving a nested name must use the Load method.
	Loader interface {
		Value
		Keyed

		// AbsoluteName returns the absolute name of this loader, i.e. the absolute name of the parent
		// namespace + '/' + this loaders name or, if this loader has no parent namespace, just this
		// loaders name prefixed with a '/'.
		AbsoluteName() string

		// Load loads value for the given name and returns it, or nil if the value could not be
		// found. A loaded value is cached.
		//
		// The name may contain the namespace separator character '/'. If it does, the loader will
		// dispatch the load using the Namespace method as many times there are segments in the name
		// - 1 and then let find the element in the lastly obtained namespace.
		Load(name string) Value

		// Name returns this loaders name relative to its parent namespace.
		Name() string

		// Namespace returns the Loader that represents the  given namespace from this loader or nil no such
		// loader exists.
		//
		// Note that this method is called internally by the Load method so there's normally no need to call
		// it explicitly.
		//
		// The name must not contain the separator character '/'. Nested namespaces must be
		// obtained using multiple calls.
		Namespace(name string) Loader

		// NewChild creates a new loader that is parented by this loader.
		NewChild(finder Finder, nsCreator NsCreator) Loader

		// ParentNamespace returns this loaders parent namespace, or nil, if this loader is in the
		// root namespace.
		ParentNamespace() Loader
	}
)
