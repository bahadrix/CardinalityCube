package cores

import (
	"crypto/sha1"
	"encoding/base64"

)

type BasicSetCore struct {
	set map[string]bool
}


var BasicSet CoreInitiator = func(opts interface{}) Core {
	return BasicSetCore{
		set: make(map[string]bool),
	}
}

func (b BasicSetCore) Push(item []byte) {
	hash := base64.StdEncoding.EncodeToString(sha1.New().Sum(item))
	b.set[hash] = true
}

func (b BasicSetCore) Count() uint64 {
	return uint64(len(b.set))
}
