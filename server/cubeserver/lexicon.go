package cubeserver

import (
	"fmt"
	"strings"
)

type Command func(server *Server, args ...string) (string, error)
type Lexicon struct {
	commands map[string]Command
}

const OK = ""




var lexicon = &Lexicon{
	commands: make(map[string]Command),
}

func (l *Lexicon) Put(name string, cmd Command) {
	name = strings.ToUpper(strings.TrimSpace(name))

	_, exists := l.commands[name]

	if exists {
		panic(fmt.Errorf("command already exists with name: %s", name))
	}

	l.commands[name] = cmd

}
// CreateInterpreter creates new interpreter from clone of lexicon
// Note that interpreters not thread safe intentionally
func(l *Lexicon) CreateInterpreter() *Interpreter {
	commands := make(map[string]Command)
	for k, v := range lexicon.commands { // clone the commands
		commands[k] = v
	}
	return &Interpreter{commands:commands}
}

func argParseError(argIndex int, err error) error {
	return fmt.Errorf("argument parse error at argument %d: %s", argIndex, err.Error())
}
