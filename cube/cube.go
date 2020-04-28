package cube

import (
	"github.com/bahadrix/cardinalitycube/cores"
	"github.com/bahadrix/cardinalitycube/cube/pb"
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
	core, _ := c.coreGenerator(c.coreOpts, nil)
	return &Cell{core: core}
}

func (c *Cube) deserializeCell(coreBytes []byte) (*Cell, error) {
	core, err := c.coreGenerator(nil, coreBytes)
	if err != nil {
		return nil, err
	}
	return &Cell{core:core}, nil
}

// GetBoardKeys returns board names. Read blocking operation
func (c *Cube) GetBoardKeys() []string {
	c.boardLock.RLock()
	keys := make([]string, 0, len(c.boardMap))
	for key := range c.boardMap {
		keys = append(keys, key)
	}
	c.boardLock.RUnlock()
	return keys
}

// GetBoardCount returns current board count in cube.
func (c *Cube) GetBoardCount() int {
	return len(c.boardMap)
}

func (c *Cube) Dump() (*pb.CubeData, error) {

	c.boardLock.RLock()
	defer c.boardLock.RUnlock()

	dataMap := make(map[string]*pb.BoardData, len(c.boardMap))
	var err error
	for k, b := range c.boardMap{
		dataMap[k], err = b.Dump()

		if err != nil {
			return nil, err
		}

	}

	return &pb.CubeData{ BoardMap:dataMap}, err

}

func (c *Cube) LoadData(data *pb.CubeData) error {

	c.boardLock.Lock()
	defer c.boardLock.Unlock()

	for boardName, boardData := range data.BoardMap {

		board, boardExists := c.boardMap[boardName]
		if !boardExists {
			board = NewBoard(c)
			c.boardMap[boardName] = board
		}

		err := board.LoadData(boardData)

		if err != nil {
			return err
		}
	}

	return nil
}