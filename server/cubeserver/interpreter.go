package cubeserver

import (
	"fmt"
	"strings"
)

// Interpreter interprets given string command and executes it on the server.
// Interpreter are not thread safe intentionally, so each routine must own an interpreter.
type Interpreter struct {
	commands map[string]*Command
}

// ParseCommandInput parses command returns command and its argument array separately.
func ParseCommandInput(input string) (command string, args []string) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return "", nil
	}
	command = strings.ToUpper(fields[0])

	if len(fields) > 1 {
		args = fields[1:]
	} else {
		args = []string{}
	}

	return
}

// Interpret parses given command and run it on server then returns the result.
func (interpreter *Interpreter) Interpret(server *Server, input string) (string, error) {

	commandName, args := ParseCommandInput(input)

	command, _ := interpreter.commands[commandName]

	if command == nil {
		if commandName == "" {
			// If you give nothing you get nothing
			return "", nil
		}
		return "", fmt.Errorf("unknown command %s", commandName)
	}

	return command.Executor(server, args...)
}
