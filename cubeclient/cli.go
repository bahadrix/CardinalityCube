package main

import (
	"errors"
	"fmt"
	"github.com/bahadrix/cardinalitycube/cubeclient/client"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)

var serverAddress string
var commandToExecute string
var cubeClient *client.Client


func executeCommand(commandString string) (reply string, err error) {
	parts := strings.Fields(commandString)

	if len(parts) < 1 {
		err = errors.New("no command given")
		return
	}

	if len(parts) > 1 {
		return cubeClient.Execute(parts[0], parts[1:]...)
	}
	return cubeClient.Execute(parts[0])

}

func init() {
	flag.StringVarP(&serverAddress, "server", "s", "tcp://127.0.0.1:1994", "Server endpoint address")
	flag.StringVarP(&commandToExecute, "exec", "e", "", "Execute command, get output and exit. If not entered, interactive client will be started.")
}

func startCli() {

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("%s> ", serverAddress),

	}

	for {
		cmd, err := prompt.Run()

		if err != nil {
			log.Errorf("Prompt failed: %s", err.Error())
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
			_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			exitCode = 1
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "%s\n", reply)
			exitCode = 0
		}

		cubeClient.Close()
		os.Exit(exitCode)
	}

	defer cubeClient.Close()
	startCli()
}
