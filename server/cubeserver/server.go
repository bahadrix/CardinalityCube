package cubeserver

import (
	"github.com/bahadrix/cardinalitycube/cube"
	log "github.com/sirupsen/logrus"
	"github.com/zeromq/goczmq"
)

const (
	// Version for server
	Version = "0.1.0"
)

const ( // Message State Prefix Bytes
	// MsgPrefixFail for failed operation replies
	MsgPrefixFail = byte(0)
	// MsgPrefixSuccess for successful operation replies
	MsgPrefixSuccess = byte(1)
)

// Server is Cardinality Cube Server
type Server struct {
	processQueue      chan *Message
	responseQueue     chan *Message
	endpoints         string
	numProcessWorkers int
	router            *goczmq.Sock
	cube              *cube.Cube
	isShuttingDown    bool
}

// Message between client and server
type Message struct {
	Route   []byte
	Message []byte
}

// NewServer creates new server
func NewServer(cubeToServe *cube.Cube, endPoints string, processQueueSize int, responseQueueSize int, numProcessWorkers int) *Server {
	return &Server{
		processQueue:      make(chan *Message, processQueueSize),
		responseQueue:     make(chan *Message, responseQueueSize),
		endpoints:         endPoints,
		numProcessWorkers: numProcessWorkers,
		cube:              cubeToServe,
	}
}

func (s *Server) startRouting() error {
	var err error
	s.router, err = goczmq.NewRouter(s.endpoints)
	if err != nil {
		return err
	}

	go func() {
		for {
			request, err := s.router.RecvMessage()
			if err == goczmq.ErrRecvFrameAfterDestroy || s.isShuttingDown {
				break
			} else if err != nil {
				log.Errorf("Error while receiving message from client: %s", err.Error())
			} else {
				s.processQueue <- &Message{
					Route:   request[0],
					Message: request[1],
				}
			}
		}
	}()

	return nil
}

func (s *Server) process() {
	interpreter := lexicon.CreateInterpreter()

	for {
		message, ok := <-s.processQueue

		if !ok {
			break
		}
		state := MsgPrefixSuccess
		reply, err := interpreter.Interpret(s, string(message.Message))
		if err != nil {
			reply = err.Error()
			state = MsgPrefixFail
		}
		s.responseQueue <- &Message{
			Route:   message.Route,
			Message: append([]byte{state}, []byte(reply)...),
		}
	}
}

// Shutdown gracefully stops server
func (s *Server) Shutdown() {
	log.Info("Shutting down server.")
	s.isShuttingDown = true
	log.Info("Destroying router.")
	s.router.Destroy()
	log.Info("Stopping processors")
	close(s.processQueue)
	log.Info("Stopping responder")
	close(s.responseQueue)

	log.Info("Goodbye!")
}

// Start initiates the server routines and waits for outputs.
func (s *Server) Start() error {

	log.Infof("Starting router for %s", s.endpoints)
	err := s.startRouting()

	if err != nil {
		return err
	}

	log.Infof("Starting command processors")
	for i := 0; i < s.numProcessWorkers; i++ {
		go s.process()
	}

	log.Infof("Connecting exhaust")

	for { // Sequential exhaust
		message, ok := <-s.responseQueue
		if !ok { // Channel is closed
			break
		}
		err = s.router.SendFrame(message.Route, goczmq.FlagMore)
		if err != nil {
			log.Errorf("Error on sending reply: %s", err.Error())
			continue
		}

		err = s.router.SendFrame(message.Message, goczmq.FlagNone)
		if err != nil {
			log.Errorf("Error on sending reply: %s", err.Error())
			continue
		}

	}

	return nil
}
