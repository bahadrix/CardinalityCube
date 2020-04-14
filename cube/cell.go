package cube

import (
	"github.com/bahadrix/cardinalitycube/cores"
	"sync"
)


type Cell struct {
	core cores.Core
	mux sync.Mutex
}

func (c *Cell) Push(item []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.core.Push(item)
}

func (c *Cell) Count() uint64 {
	return c.core.Count()
}