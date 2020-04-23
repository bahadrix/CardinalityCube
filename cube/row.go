package cube

import "sync"

// Row is a build block of board. It holds Cells.
type Row struct {
	cellMap map[string]*Cell
	mux     sync.RWMutex
}

// RowSnapshot contains cell names and their accumulated values.
type RowSnapshot map[string]uint64

// NewRow creates empty row and returns it.
func NewRow() *Row {
	return &Row{
		cellMap: map[string]*Cell{},
	}
}

// GetCell returns cell at given name.
// Returns nil if it not found.
func (r *Row) GetCell(cellName string) *Cell {
	r.mux.RLock()
	cell, _ := r.cellMap[cellName]
	r.mux.RUnlock()
	return cell
}

// SetCell sets the given cell to given name
func (r *Row) SetCell(cellName string, cell *Cell) {
	r.mux.Lock()
	r.cellMap[cellName] = cell
	r.mux.Unlock()
}

// GetSnapshot returns snapshot of row. Blocks row while getting snapshot.
func (r *Row) GetSnapshot() *RowSnapshot {
	ss := make(RowSnapshot)
	r.mux.RLock()
	for key, cell := range r.cellMap {
		ss[key] = cell.Count()
	}
	r.mux.RUnlock()
	return &ss
}

// GetCellKeys returns keys of cells. Read blocking operation.
func (r *Row) GetCellKeys() []string {
	r.mux.RLock()
	keys := make([]string, 0, len(r.cellMap))
	for key := range r.cellMap {
		keys = append(keys, key)
	}
	r.mux.RUnlock()
	return keys
}

// GetCellCount returns cell count of row.
func (r *Row) GetCellCount() int {
	return len(r.cellMap)
}
