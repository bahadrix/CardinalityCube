package cube

import (
	"github.com/bahadrix/cardinalitycube/cores"
	"sync"
)

type Cube struct {
	boardMap map[string]*Board
	coreGenerator cores.CoreInitiator
	coreOpts interface{}
	boardLock sync.RWMutex
}

type Snapshot map[string]*BoardSnapshot

func NewCube(coreGenerator cores.CoreInitiator, coreOpts interface{}) *Cube {
	return &Cube{
		boardMap: map[string]*Board{},
		coreGenerator: coreGenerator,
		coreOpts:      coreOpts,
	}
}

func (c *Cube) GetBoard(name string, createIfNotExists bool) *Board {
	c.boardLock.RLock() // Concurrent map read is not allowed
	board, _ := c.boardMap[name]
	c.boardLock.RUnlock()

	if board == nil && createIfNotExists {
		c.boardLock.Lock() // board sync in
		board, _ = c.boardMap[name]
		if board == nil { // board still not exists
			board = NewBoard(c)
			c.boardMap[name] = board
		}
		c.boardLock.Unlock() // board sync out
	}
	return board
}

func (c *Cube) GetSnapshot() *Snapshot {
	ss := make(Snapshot)
	c.boardLock.RLock()
	for key, board := range c.boardMap {
		ss[key] = board.GetSnapshot()
	}
	c.boardLock.RUnlock()
	return &ss
}

func (c *Cube) generateCell() *Cell {
	core := c.coreGenerator(c.coreOpts)
	return &Cell{core:core}
}

func d() {
	cube := NewCube(cores.HLL, nil)
	board := cube.GetBoard("s", true)
	cell := board.GetCell("e", "d", true)
	_ = cell


}