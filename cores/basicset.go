package cores

import (
	"crypto/sha1" //nolint:gosec
	"encoding/base64"
)

type BasicSetCore struct {
	set map[string]bool
}

func BasicSet(opts interface{}) Core {
	return BasicSetCore{
		set: make(map[string]bool),
	}
}

func (b BasicSetCore) Push(item []byte) {
	hash := base64.StdEncoding.EncodeToString(sha1.New().Sum(item)) //nolint:gosec
	b.set[hash] = true
}

func (b BasicSetCore) Count() uint64 {
	return uint64(len(b.set))
}
