package cubeserver

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bahadrix/cardinalitycube/cube"
	"sort"
	"strconv"
	"strings"
)

func cmdPush(server *Server, args ...string) (s string, err error) {
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
}

func cmdGet(server *Server, args ...string) (s string, err error) {
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

}

func cmdSnapshot(server *Server, args ...string) (s string, err error) {
	var board *cube.Board

	if len(args) > 0 {
		board = server.cube.GetBoard(args[0], false)
		if board == nil {
			return "", nil
		}
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
}

func cmdDrop(server *Server, args ...string) (s string, err error) {
	var board *cube.Board

	if len(args) == 2 { // Drop row
		board = server.cube.GetBoard(args[0], false)
		if board != nil {
			board.DropRow(args[1])
		}
	} else if len(args) == 1 { // Drop board
		server.cube.DropBoard(args[0])
	} else {
		err = errors.New("at least board parameter required to drop a board, or both board and row parameters required to drop a row")
	}

	return
}

func cmdExists(server *Server, args ...string) (s string, err error) {
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
}

func cmdList(server *Server, args ...string) (s string, err error) {
	argsLen := len(args)
	var keys []string

	if argsLen == 0 { // Get board keys of cube
		keys = server.cube.GetBoardKeys()
	} else {
		board := server.cube.GetBoard(args[0], false)
		if board == nil {
			return "", nil
		}

		if argsLen == 1 { // Get row keys of board
			keys = board.GetRowKeys()
		} else if argsLen == 2 { // Get cell keys of row
			keys = board.GetCellKeys(args[1])
		} else {
			return "", errors.New("command takes max 3 arguments")
		}
	}

	sort.Strings(keys)
	return strings.Join(keys, "\n"), nil
}

func init() {
	lexicon.Put("PUSH", &Command{
		ShortDescription: "Push given item value to cube.",
		Description:      "Usage: PUSH <board> <row> <cell> <item>",
		Executor:         cmdPush,
	})

	lexicon.Put("GET", &Command{
		ShortDescription: "Returns current count of given cell",
		Description:      "Usage: PUSH <board> <row> <cell>",
		Executor:         cmdGet,
	})

	lexicon.Put("SNAPSHOT", &Command{
		ShortDescription: "Returns current snap shot of given path",
		Description:      "Usage: SNAPSHOT <board> [<row>]",
		Executor:         cmdSnapshot,
	})

	lexicon.Put("DROP", &Command{
		ShortDescription: "Drops board or row",
		Description:      "Usage: DROP board [row]",
		Executor:         cmdDrop,
	})

	lexicon.Put("EXISTS", &Command{
		ShortDescription: "Check the existence of board, row or cell",
		Description:      "Usage: EXISTS <board> [row [cell]]",
		Executor:         cmdExists,
	})

	lexicon.Put("LIST", &Command{
		ShortDescription: "Get list of board, row or cell names",
		Description:      "Usage: LIST [board [row]]",
		Executor:         cmdList,
	})
}
