package cmd

import (
	"fmt"
	"github.com/bahadrix/cardinalitycube/cores"
	"github.com/bahadrix/cardinalitycube/cube"
	"github.com/bahadrix/cardinalitycube/server/cubeserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// startCmd represents the serve command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Cardinality Cube server",
	Long:  "Start a Cardinality Cube server",
	Run: func(cmd *cobra.Command, args []string) {
		tcpAddr, _ := cmd.Flags().GetString("listen")
		queueSize, _ := cmd.Flags().GetInt("command-buffer")
		workers, _ := cmd.Flags().GetInt("workers")
		coreType, _ := cmd.Flags().GetString("core")

		coreType = strings.ToLower(coreType)
		var kube *cube.Cube
		switch coreType {
		case "hll":
			kube = cube.NewCube(cores.HLL, nil)
		case "basicset":
			kube = cube.NewCube(cores.BasicSet, nil)
		default:
			fmt.Println("Only hll and basicset type supported.")
			os.Exit(128)
			return
		}

		cubeServer := cubeserver.NewServer(kube, tcpAddr, queueSize, queueSize, workers)

		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-s
			cubeServer.Shutdown()
		}()

		err := cubeServer.Start()

		if err != nil {
			log.Error(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP("listen", "l", "tcp://0.0.0.0:1994", "Host and port address for listening TCP connections")
	startCmd.Flags().IntP("command-buffer", "b", 10000, "Executor buffer limit")
	startCmd.Flags().IntP("workers", "w", 12, "Number of parallel command processors")
	startCmd.Flags().StringP("core", "c", "hll", "Core type of cube. Currently hll and basicset supported")

}
