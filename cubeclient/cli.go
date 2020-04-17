package main

import (
	"encoding/json"
	"fmt"
	"github.com/bahadrix/cardinalitycube/cubeclient/client"
	"github.com/bahadrix/cardinalitycube/server/cubeserver"
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)

var serverAddress string
var commandToExecute string
var cubeClient *client.Client

var suggestions []prompt.Suggest

func executeCommand(commandString string) (reply string, err error) {
	cmd, args := cubeserver.ParseCommandInput(commandString)
	return cubeClient.Execute(cmd, args...)

}

func init() {
	flag.StringVarP(&serverAddress, "server", "s", "tcp://127.0.0.1:1994", "Server endpoint address")
	flag.StringVarP(&commandToExecute, "exec", "e", "", "Execute command, get output and exit. If not entered, interactive client will be started.")
}

func autoComplete(in prompt.Document) []prompt.Suggest {

	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	if len(strings.Fields(in.Text)) >= 2 {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func loadSuggestions() {
	var lexicon map[string]map[string]string

	lexiconJSON, err := executeCommand("LEXICON")

	if err != nil {
		fmt.Printf("Unable to receive lexicon %s\n", err.Error())
	}

	err = json.Unmarshal([]byte(lexiconJSON), &lexicon)
	if err != nil {
		fmt.Printf("Unable to parse lexicon %s\n", err.Error())
	}

	cmdCount := len(lexicon)
	suggestions = make([]prompt.Suggest, cmdCount)

	i := -1
	for cmd, info := range lexicon {
		i++
		suggestions[i] = prompt.Suggest{
			Text:        cmd,
			Description: info["short"],
		}
	}

	suggestions = append(suggestions, prompt.Suggest{
		Text:        "EXIT",
		Description: "Terminate console session",
	})

}

func startCli() {

	loadSuggestions()

	for {

		cmd := prompt.Input(fmt.Sprintf("%s> ", serverAddress), autoComplete,
			prompt.OptionTitle("Cardinality Cube Server"),
			prompt.OptionPrefixTextColor(prompt.DarkBlue),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSelectedSuggestionTextColor(prompt.DarkGray),
			prompt.OptionDescriptionTextColor(prompt.DarkGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))

		if strings.ToLower(cmd) == "exit" {
			break
		}
		reply, err := executeCommand(cmd)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		} else {
			fmt.Printf("%s\n", reply)
		}

	}

}

func main() {
	flag.Parse()

	var err error

	cubeClient, err = client.NewClient(serverAddress)

	if err != nil {
		log.Errorf("Error on connecting client: %s", err.Error())
		os.Exit(1)
	}

	if commandToExecute != "" {
		reply, err := executeCommand(commandToExecute)
		exitCode := 0
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Fail: %s\n", err.Error())
			exitCode = 1
		} else {
			if reply == "" {
				reply = "OK"
			}
			_, _ = fmt.Fprintf(os.Stdout, "%s\n", reply)
			exitCode = 0
		}

		cubeClient.Close()
		os.Exit(exitCode)
	}

	defer cubeClient.Close()
	startCli()
}
