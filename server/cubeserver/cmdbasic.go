package cubeserver
// Basic server commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func init() {
	lexicon.Put("PING", func(server *Server, args ...string) (s string, err error) {
		reply := "PONG"
		if len(args) > 0 {
			reply = fmt.Sprintf("PONG %s", args[0])
		}
		return reply, nil
	})

	lexicon.Put("VERSION", func(server *Server, args ...string) (s string, err error) {
		return Version, nil
	})

	lexicon.Put("SHUTDOWN", func(server *Server, args ...string) (s string, err error) {
		var delay uint64 = 1

		if len(args) > 0 {
			delay, err = strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				err = argParseError(0, err)
				return
			}
		}

		defer func() {
			go func() {
				log.Infof("Shutting down in %d seconds", delay)
				time.Sleep(time.Duration(delay) * time.Second)
				server.Shutdown()
			}()
		}()

		return OK, nil
	})
}