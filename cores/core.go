package cores

type Core interface {
	// Push adds item to store to count. Thread safe supplied by caller.
	Push(item []byte)
	// Count returns unique item count in the store. Not thread safe.
	Count() uint64
}

// CoreInitiator defines a function that initializes the cell core.
// The function will be called lazily whenever new cell is needed.
// Use nil value for opts argument to default options
type CoreInitiator func(opts interface{}) Core
