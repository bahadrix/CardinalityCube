package cubeserver

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bahadrix/cardinalitycube/cube"
	"strconv"
)

func init() {
	lexicon.Put("PUSH", &Command{
		ShortDescription: "Push given item value to cube.",
		Description:      "Usage: PUSH <board> <row> <cell> <item>",
		Executor: func(server *Server, args ...string) (s string, err error) {
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
		},
	})

	lexicon.Put("GET", &Command{
		ShortDescription: "Returns current count of given cell",
		Description:      "Usage: PUSH <board> <row> <cell>",
		Executor: func(server *Server, args ...string) (s string, err error) {
			if len(args) != 3 {
				err = errors.New("GET method requires exactly 3 arguments")
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

		},
	})


	lexicon.Put("SNAPSHOT", &Command{
		ShortDescription: "Returns current snap shot of given path",
		Description:      "Usage: SNAPSHOT <board> [<row>]",
		Executor: func(server *Server, args ...string) (s string, err error) {
			var board *cube.Board


			if len(args) > 0 {
				board = server.cube.GetBoard(args[0], false)
			} else {
				err = errors.New("at least board parameter required")
				return
			}

			if len(args) > 1 { // Get row snapshot
				rowSnapshot := board.GetRowSnapshot(args[1])

				if rowSnapshot == nil {
					return
				}

				var buffer bytes.Buffer
				for cellName, value := range *rowSnapshot {
					buffer.WriteString(fmt.Sprintf("%s\t%d\n", cellName, value))
				}
				return buffer.String(), nil
			}

			// Get board snapshot

			boardSnapshot := board.GetSnapshot()

			if boardSnapshot == nil {
				return
			}

			var buffer bytes.Buffer
			for rowName, rowSnapshot := range *boardSnapshot {
				for cellName, value := range *rowSnapshot {
					buffer.WriteString(fmt.Sprintf("%s\t%s\t%d\n", rowName, cellName, value))
				}
			}
			return buffer.String(), nil

		},
	})

	lexicon.Put("DROP", &Command{
		ShortDescription: "Drops board or row",
		Description:      "Usage: DROP board [row]",
		Executor: func(server *Server, args ...string) (s string, err error) {
			var board *cube.Board


			if len(args) > 0 {
				board = server.cube.GetBoard(args[0], false)
			} else {
				err = errors.New("at least board parameter required")
				return
			}

			if len(args) > 1 { // Drop row
				board.DropRow(args[1])
				return
			}

			board.Drop()

			return
		},
	})

	lexicon.Put("EXISTS", &Command{
		ShortDescription: "Check the existence of board, row or cell",
		Description:      "Usage: EXISTS <board> [row [cell]]",
		Executor: func(server *Server, args ...string) (s string, err error) {
			argsLen := len(args)
			var board *cube.Board

			if argsLen > 0 {
				board = server.cube.GetBoard(args[0], false)
				if board == nil {
					return "0", nil
				}
			} else {
				err = errors.New("at least board parameter required")
				return
			}

			if argsLen > 1 {
				rowExists := board.CheckRowExists(args[1])
				if !rowExists {
					return "0", nil
				}
			}

			if argsLen > 2 {
				cell := board.GetCell(args[1], args[2], false)
				if cell == nil {
					return "0", nil
				}
			}

			return "1", nil

		},
	})
}
