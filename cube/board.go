package cube

import "sync"

type Board struct {
	cube *Cube
	rowMap map[string]*Row
	rowLock sync.RWMutex
	cellLock sync.Mutex
}

type BoardSnapshot map[string]*RowSnapshot

func NewBoard(cube *Cube) *Board {
	return &Board{
		cube: cube,
		rowMap: make(map[string]*Row),
	}
}

func (b *Board) GetCell(rowName string, cellName string, createIfNotExists bool) *Cell {

	var cell *Cell
	b.rowLock.RLock()
	row, _ := b.rowMap[rowName]
	b.rowLock.RUnlock()

	if row == nil {
		if !createIfNotExists {
			return nil
		}
		b.rowLock.Lock() // row sync in ----
		row, _ = b.rowMap[rowName]
		if row == nil {
			row = NewRow()
			b.rowMap[rowName] = row
		}
		b.rowLock.Unlock() // row sync out ---
	}

	cell = row.GetCell(cellName)

	if cell == nil && createIfNotExists {
		b.cellLock.Lock() // cell sync in ----
		cell = row.GetCell(cellName)
		if cell == nil {
			cell = b.cube.generateCell()
			row.SetCell(cellName, cell)
		}
		b.cellLock.Unlock() // cell sync out ---
	}

	return cell
}

func (b* Board) GetSnapshot() *BoardSnapshot {
	ss := make(BoardSnapshot)
	b.rowLock.RLock()
	for key, row := range b.rowMap {
		ss[key] = row.GetSnapshot()
	}
	b.rowLock.RUnlock()
	return &ss
}