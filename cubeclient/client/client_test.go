package client

import (
	"encoding/json"
	"github.com/bahadrix/cardinalitycube/cores"
	"github.com/bahadrix/cardinalitycube/cube"
	"github.com/bahadrix/cardinalitycube/server/cubeserver"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var testServer *cubeserver.Server

const (
	TestSocket = "ipc://test.sock"
)

func TestMain(m *testing.M) {
	testServer = cubeserver.NewServer(cube.NewCube(cores.BasicSet, nil), TestSocket, 10, 10, 4)

	go func() {
		err := testServer.Start()
		if err != nil {
			log.Fatalf("Server routine error %s", err.Error())
		}
		defer os.Remove("test.sock")
	}()
	time.Sleep(time.Second) // wait for server warmup
	r := m.Run()

	testServer.Shutdown()

	defer os.Exit(r)
}

func TestClient_Commands(t *testing.T) {

	cli, err := NewClient(TestSocket)
	if err != nil {
		log.Fatalf("Can't create client: %s", err.Error())
	}
	defer cli.Close()

	tests := []struct {
		Command          string
		ExpectedResponse string
	}{
		{
			Command:          "PING test",
			ExpectedResponse: "PONG test",
		},
		{
			Command:          "PUSH board_1 row_1 cell_1 value_1",
			ExpectedResponse: "",
		},
		{
			Command:          "PUSH board_1 row_1 cell_1 value_1",
			ExpectedResponse: "",
		},
		{
			Command:          "PUSH board_1 row_1 cell_1 value_2",
			ExpectedResponse: "",
		},
		{
			Command:          "GET board_1 row_1 cell_1",
			ExpectedResponse: "2",
		},
		{
			Command:          "EXISTS board_1 row_1 cell_1",
			ExpectedResponse: "1",
		},
		{
			Command:          "EXISTS board_non_exist row_1 cell_1",
			ExpectedResponse: "0",
		},
		{
			Command:          "SNAPSHOT board_1 row_1",
			ExpectedResponse: "cell_1\t2\n",
		},
		{
			Command:          "SNAPSHOT board_1",
			ExpectedResponse: "row_1\tcell_1\t2\n",
		},
		{
			Command:          "DROP board_1 row_1",
			ExpectedResponse: "",
		},
		{
			Command:          "EXISTS board_1 row_1",
			ExpectedResponse: "0",
		},
		{
			Command:          "DROP board_1",
			ExpectedResponse: "",
		},
		{
			Command:          "EXISTS board_1",
			ExpectedResponse: "0",
		},
	}

	for _, test := range tests {
		t.Run(test.Command, func(t *testing.T) {
			cmd, args := cubeserver.ParseCommandInput(test.Command)
			reply, err := cli.Execute(cmd, args...)

			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, test.ExpectedResponse, reply)
		})
	}

}

func TestClient_Lexicon(t *testing.T) {
	cli, err := NewClient(TestSocket)
	if err != nil {
		t.Error(err)
	}
	defer cli.Close()

	reply, err := cli.Execute("LEXICON")

	var lx interface{}
	err = json.Unmarshal([]byte(reply), &lx)
	if err != nil {
		t.Error(err)
	}

}
