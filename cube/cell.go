package cube

import (
	"github.com/bahadrix/cardinalitycube/cores"
	"sync"
)

// Cell is constituent part of Cube. It holds values.
type Cell struct {
	core cores.Core
	mux  sync.Mutex
}

// Push pushes item into core.
func (c *Cell) Push(item []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.core.Push(item)
}

// Count does the accumulation at core and returns the result
func (c *Cell) Count() uint64 {
	return c.core.Count()
}
