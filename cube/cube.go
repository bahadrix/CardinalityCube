package cube

import (
	"github.com/bahadrix/cardinalitycube/cores"
	"sync"
)

// Cube is a data structure consists of cells that allows thread safe parallel operations over them.
// A cube organizes it cells by holding them in maps called rows. Rows are also combined into a another map
// called board.
type Cube struct {
	boardMap      map[string]*Board
	coreGenerator cores.CoreInitiator
	coreOpts      interface{}
	boardLock     sync.RWMutex
}

// A Snapshot of cube is a map of board snapshots
type Snapshot map[string]*BoardSnapshot

// NewCube creates new cube in the type of give CoreInitiator
func NewCube(coreGenerator cores.CoreInitiator, coreOpts interface{}) *Cube {
	return &Cube{
		boardMap:      map[string]*Board{},
		coreGenerator: coreGenerator,
		coreOpts:      coreOpts,
	}
}

// GetBoard returns board at given key name. Returns nil if it not found
// or returns newly created one if createIfNotExists is set to true.
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

// GetSnapshot Returns snapshot of whole cube. Blocking operation.
func (c *Cube) GetSnapshot() *Snapshot {
	ss := make(Snapshot)
	c.boardLock.RLock()
	for key, board := range c.boardMap {
		ss[key] = board.GetSnapshot()
	}
	c.boardLock.RUnlock()
	return &ss
}

// DropBoard deletes board at given name.
func (c *Cube) DropBoard(boardName string) {
	c.boardLock.Lock()
	_, exists := c.boardMap[boardName]
	if exists {
		delete(c.boardMap, boardName)
	}
	c.boardLock.Unlock()
}

func (c *Cube) generateCell() *Cell {
	core := c.coreGenerator(c.coreOpts)
	return &Cell{core: core}
}
