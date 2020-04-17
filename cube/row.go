package cube

import "sync"

type Row struct {
	cellMap map[string]*Cell
	mux     sync.RWMutex
}

type RowSnapshot map[string]uint64

func NewRow() *Row {
	return &Row{
		cellMap: map[string]*Cell{},
	}
}

func (r *Row) GetCell(cellName string) *Cell {
	r.mux.RLock()
	cell, _ := r.cellMap[cellName]
	r.mux.RUnlock()
	return cell
}

func (r *Row) SetCell(cellName string, cell *Cell) {
	r.mux.Lock()
	r.cellMap[cellName] = cell
	r.mux.Unlock()
}

func (r *Row) GetSnapshot() *RowSnapshot {
	ss := make(RowSnapshot)
	r.mux.RLock()
	for key, cell := range r.cellMap {
		ss[key] = cell.Count()
	}
	r.mux.RUnlock()
	return &ss
}
