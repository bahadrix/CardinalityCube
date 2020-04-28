package cube

import (
	"github.com/bahadrix/cardinalitycube/cores"
	"github.com/bahadrix/cardinalitycube/cube/pb"
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

// Serialized returns serialized core bytes
func (c *Cell) Serialize() ([]byte, error) {
	return c.core.Serialize()
}

// Dump returns protobuf object
func (c *Cell) Dump() (*pb.CellData, error) {
	c.mux.Lock()
	data, err := c.core.Serialize()
	c.mux.Unlock()

	if err != nil {
		return nil, err
	}

	return &pb.CellData{CoreData:data}, nil
}