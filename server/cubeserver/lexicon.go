package cubeserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Executor function does all the work for command.
type Executor func(server *Server, args ...string) (string, error)

// Command is a constituent of server language
type Command struct {
	ShortDescription string   `json:"short"`
	Description      string   `json:"description"`
	Executor         Executor `json:"-"`
}

// A Lexicon consists of all commands of server
type Lexicon struct {
	commands map[string]*Command
}

// OK is a reply for successful but idempotent operations
const OK = ""

var lexicon = &Lexicon{
	commands: make(map[string]*Command),
}

// Put adds a new command to lexicon
func (l *Lexicon) Put(name string, cmd *Command) {
	name = strings.ToUpper(strings.TrimSpace(name))

	_, exists := l.commands[name]

	if exists {
		panic(fmt.Errorf("command already exists with name: %s", name))
	}

	l.commands[name] = cmd

}

// CreateInterpreter creates new interpreter from clone of lexicon
// Note that interpreters not thread safe intentionally
func (l *Lexicon) CreateInterpreter() *Interpreter {
	commands := make(map[string]*Command)
	for k, v := range lexicon.commands { // Shallow clone is enough for evading parallel read access situation
		commands[k] = v
	}
	return &Interpreter{commands: commands}
}

// AsJSON Converts whole lexicon to a JSON string
func (l *Lexicon) AsJSON() (string, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(l.commands)

	if err != nil {
		return "", err
	}

	return buffer.String(), err
}

func argParseError(argIndex int, err error) error {
	return fmt.Errorf("argument parse error at argument %d: %s", argIndex, err.Error())
}
