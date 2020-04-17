package cubeserver

// Basic server commands

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func cmdPing(server *Server, args ...string) (s string, err error) {
	reply := "PONG"
	if len(args) > 0 {
		reply = fmt.Sprintf("PONG %s", args[0])
	}
	return reply, nil
}

func cmdVersion(server *Server, args ...string) (s string, err error) {
	return Version, nil
}

func cmdShutdown(server *Server, args ...string) (s string, err error) {
	var delay uint64 = 1

	if len(args) > 0 {
		delay, err = strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			err = argParseError(0, err)
			return
		}
	} else {
		err = errors.New("seconds_before_shutdown argument is required")
		return
	}

	defer func() {
		go func() {
			log.Infof("Shutting down in %d seconds", delay)
			time.Sleep(time.Duration(delay) * time.Second)
			server.Shutdown()
		}()
	}()

	return OK, nil
}

func cmdLexicon(server *Server, args ...string) (s string, err error) {
	return lexicon.AsJson()
}

func init() {

	lexicon.Put("PING", &Command{
		ShortDescription: "Simple PING. Server responds with pong",
		Description:      "Server also responds with first argument if defined",
		Executor:         cmdPing,
	})

	lexicon.Put("VERSION", &Command{
		ShortDescription: "Returns server version",
		Description:      "",
		Executor:         cmdVersion,
	})

	lexicon.Put("SHUTDOWN", &Command{
		ShortDescription: "Shutdown server",
		Description:      "Usage: SHUTDOWN <seconds_before_shutdown>",
		Executor:         cmdShutdown,
	})

	lexicon.Put("LEXICON", &Command{
		ShortDescription: "Returns current command lexicon as json",
		Description:      "",
		Executor:         cmdLexicon,
	})
}
