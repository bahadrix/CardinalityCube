package cores

// Core is an abstraction of basic functionality of a Cube cell.
// Different implementations of a core result in different cubes.
type Core interface {
	// Push adds item to store to count. Thread safe supplied by caller.
	Push(item []byte)
	// Count returns unique item count in the store. Not thread safe.
	Count() uint64
	// Serialize core into bytes
	Serialize() ([]byte, error)
}

// CoreInitiator defines a function that initializes the cell core.
// The function will be called lazily whenever new cell is needed or
// at the deserialization phase.
// Use nil value for each argument to create new core with default options.
// If serializedBytes is not nil deserialization will be initiated.
type CoreInitiator func(opts interface{}, serializedBytes []byte) (Core, error)
