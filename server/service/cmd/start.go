package cmd

import (
	"github.com/bahadrix/cardinalitycube/cores"
	"github.com/bahadrix/cardinalitycube/cube"
	"github.com/bahadrix/cardinalitycube/server/cubeserver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
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

		kube := cube.NewCube(cores.HLL, nil)

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

	startCmd.Flags().StringP("listen", "t", "tcp://0.0.0.0:1994", "Host and port address for listening TCP connections")
	startCmd.Flags().IntP("command-buffer", "q", 10000, "Command buffer limit")
	startCmd.Flags().IntP("workers", "w", 12, "Number of parallel command processors")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
