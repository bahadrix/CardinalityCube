package cubeserver

import (
	"errors"
	"strconv"
)

func init() {
	lexicon.Put("PUSH", func(server *Server, args ...string) (s string, err error) {
		if len(args) != 4 {
			err = errors.New("PUSH method requires exactly 4 arguments")
			return
		}
		boardName := args[0]
		rowName := args[1]
		cellName := args[2]
		item := args[3]

		board := server.cube.GetBoard(boardName, true)
		cell := board.GetCell(rowName, cellName, true)
		cell.Push([]byte(item))

		return
	})

	lexicon.Put("GET", func(server *Server, args ...string) (s string, err error) {
		if len(args) != 3 {
			err = errors.New("PUSH method requires exactly 4 arguments")
			return
		}
		boardName := args[0]
		rowName := args[1]
		cellName := args[2]

		board := server.cube.GetBoard(boardName, false)
		if board == nil {
			return
		}

		cell := board.GetCell(rowName, cellName, false)
		if cell == nil {
			return
		}

		return strconv.FormatUint(cell.Count(), 10), nil

	})
}
