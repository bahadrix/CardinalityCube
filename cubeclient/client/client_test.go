package client

import (
	"github.com/bahadrix/cardinalitycube/cores"
	"github.com/bahadrix/cardinalitycube/cube"
	"github.com/bahadrix/cardinalitycube/server/cubeserver"
	"github.com/bmizerany/assert"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var testServer *cubeserver.Server

const (
	TestSocket = "ipc://test.sock"
)

func TestMain(m *testing.M) {
	testServer = cubeserver.NewServer(cube.NewCube(cores.HLL, nil), TestSocket, 10, 10, 4)

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
	time.Sleep(2 * time.Second) // grace period
	os.Exit(r)
}


func TestClient_Close(t *testing.T) {

	cli, err := NewClient(TestSocket)

	if err != nil {
		log.Fatalf("Can't create client: %s", err.Error())
	}

	log.Info("Executing ping command")
	reply, err := cli.Execute("PING")

	if err != nil {
		log.Fatalf("Error on getting reply from server %s", err.Error())
	}

	assert.Equal(t, "PONG", reply)


	cli.Close()
}
