package cores

import (
	"crypto/sha1" //nolint:gosec
	"encoding/base64"
	"encoding/json"
)

// BasicSetCore is a very basic core that counts exact distinct items.
// There is no approximation in this approach but memory usage is inefficient.
type BasicSetCore struct {
	set map[string]bool
}


// BasicSet is a CoreInitiator
func BasicSet(opts interface{}, serializedBytes []byte) (Core, error) {
	var set map[string]bool

	if serializedBytes != nil {
		err := json.Unmarshal(serializedBytes, &set)
		if err != nil {
			return nil, err
		}
	} else {
		set = make(map[string]bool)
	}

	return BasicSetCore{
		set: set,
	}, nil
}

// Push pushes item into set
func (b BasicSetCore) Push(item []byte) {
	hash := base64.StdEncoding.EncodeToString(sha1.New().Sum(item)) //nolint:gosec
	b.set[hash] = true
}

// Count returns item count in the set
func (b BasicSetCore) Count() uint64 {
	return uint64(len(b.set))
}

func (b BasicSetCore) Serialize() ([]byte, error) {
	return json.Marshal(b.set)
}
