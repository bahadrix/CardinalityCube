package cubeserver
// Basic server commands

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func init() {

	lexicon.Put("PING", &Command{
		ShortDescription: "Simple PING. Server responds with pong",
		Description:      "Server also responds with first argument if defined",
		Executor:         func(server *Server, args ...string) (s string, err error) {
			reply := "PONG"
			if len(args) > 0 {
				reply = fmt.Sprintf("PONG %s", args[0])
			}
			return reply, nil
		},

	})


	lexicon.Put("VERSION", &Command{
		ShortDescription: "Returns server version",
		Description:      "",
		Executor:         func(server *Server, args ...string) (s string, err error) {
			return Version, nil
		},
	})

	lexicon.Put("SHUTDOWN", &Command{
		ShortDescription: "Shutdown server",
		Description:      "Shutdowns server. First argument is seconds before shutdown, default is 1",
		Executor:         func(server *Server, args ...string) (s string, err error) {
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
		},
	})
}