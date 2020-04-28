package cube

import (
	"fmt"
	"github.com/bahadrix/cardinalitycube/cores"
	"github.com/bahadrix/cardinalitycube/cube/pb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"reflect"

	"sync"

	"testing"
)

func TestCubeThreadSafety(t *testing.T) {

	setCube := NewCube(cores.BasicSet, nil)

	numBoards := 8
	numRows := 100
	numCells := 3
	valuePerCell := 10

	totalItems := numBoards * numRows * numCells * valuePerCell

	numPushThreads := numBoards

	// Prepare test set
	type Item struct {
		Board string
		Row   string
		Cell  string
		Value []byte
	}

	items := make(chan *Item, totalItems)

	i := -1
	for b := 0; b < numBoards; b++ {
		for r := 0; r < numRows; r++ {
			for c := 0; c < numCells; c++ {
				for v := 0; v < valuePerCell; v++ {
					i++
					items <- &Item{
						Board: fmt.Sprintf("board_%d", b),
						Row:   fmt.Sprintf("row_%d", r),
						Cell:  fmt.Sprintf("cell_%d", c),
						Value: []byte(fmt.Sprintf("value_%d", i)),
					}
				}

			}
		}
	}

	close(items)

	// Test push
	t.Run(fmt.Sprintf("Concurrent push %d", totalItems), func(t *testing.T) {
		var wg sync.WaitGroup

		for i := 0; i < numPushThreads; i++ {
			wg.Add(1)
			go func() {
				for {
					item, more := <-items
					if !more { // Queue exhausted
						wg.Done()
						break
					}

					board := setCube.GetBoard(item.Board, true)
					cell := board.GetCell(item.Row, item.Cell, true)
					cell.Push(item.Value)
				}

			}()

		}

		wg.Wait()
	})

	// Test read
	t.Run("Concurrent read", func(t *testing.T) {

		var wg sync.WaitGroup
		for b := 0; b < numBoards; b++ { // Check each board concurrently
			wg.Add(1)
			go func(boardIndex int) {
				// also test snapshot concurrency
				_ = setCube.GetSnapshot()

				for r := 0; r < numRows; r++ {
					for c := 0; c < numCells; c++ {
						boardName := fmt.Sprintf("board_%d", boardIndex)
						rowName := fmt.Sprintf("row_%d", r)
						cellName := fmt.Sprintf("cell_%d", c)

						board := setCube.GetBoard(boardName, false)
						assert.NotNil(t, board, fmt.Sprintf("Board %s not exists", boardName))

						cell := board.GetCell(rowName, cellName, false)
						assert.NotNil(t, board, fmt.Sprintf("Cell %s not exists at %s, %s", cellName, boardName, rowName))

						assert.Equal(t, uint64(valuePerCell), cell.Count(), "Count mismatch")
					}
				}
				wg.Done()
			}(b)

		}

		wg.Wait()
	})

}

func TestCube_Serialization(t *testing.T) {


	subjects := []struct {
		name string
		cube *Cube
	}{
		{
			name: "Basic set cube",
			cube: NewCube(cores.BasicSet, nil),
		},
		{
			name: "HLL Cube",
			cube: NewCube(cores.HLL, nil),
		},
	}

	ops := []func(c *Cube){
		func(c *Cube) {
			// do nothing
		},
		func(c *Cube) {
			c.GetBoard("b1", true).GetCell("r1", "c1", true).Push([]byte("test"))
		},
		func(c *Cube) {
			c.GetBoard("b2", true).GetCell("r1", "c1", true).Push([]byte("test"))
			c.GetBoard("b2", true).GetCell("r1", "c1", true).Push([]byte("test"))
		},
		func(c *Cube) {
			c.GetBoard("b3", true).GetCell("r1", "c1", true).Push([]byte("test"))
			c.GetBoard("b3", true).GetCell("r1", "c2", true).Push([]byte("test"))
		},
	}

	for _, subject := range subjects {
		for i, op := range ops {
			t.Run(fmt.Sprintf("Test_%d_for_%s", i, subject.name), func(t *testing.T) {

				op(subject.cube)

				dataObj, err := subject.cube.Dump()

				if err != nil {
					t.Error(err)
				}

				data, err := proto.Marshal(dataObj)

				if err != nil {
					t.Error(err)
				}

				var readObj pb.CubeData

				err = proto.Unmarshal(data, &readObj)

				if err != nil {
					t.Error(err)
				}

				kube := NewCube(subject.cube.coreGenerator, subject.cube.coreOpts)
				err = kube.LoadData(&readObj)

				if err != nil {
					t.Error(err)
				}

				// Compare
				origSnap := subject.cube.GetSnapshot()
				loadSnap := kube.GetSnapshot()
				if !reflect.DeepEqual(origSnap, loadSnap) {
					t.Error("Snapshots not equal")
				}

				// Add after test

				kube.GetBoard("bX", true).GetCell("r1", "c1", true).Push([]byte("test"))
				subject.cube.GetBoard("bX", true).GetCell("r1", "c1", true).Push([]byte("test"))

				if !reflect.DeepEqual(subject.cube.GetSnapshot(), kube.GetSnapshot()) {
					t.Error("Add after Snapshots not equal")
				}
			})
		}
	}

}