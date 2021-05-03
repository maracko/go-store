package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/maracko/go-store/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the server TCP command
var serveTCPCmd = &cobra.Command{
	Use:   "TCP",
	Short: "Start TCP server",
	Long: `Starts a server over TCP.
	Current limitation of TCP server is that it stores/reads values as strings while HTTP wil encode/decode to json.
	Defaults to port 8888. Empty database will be initialized and kept only in memory if no path is provided.
	If you have json file with data you want to read from but not save provide a location along with -m (memory) flag`,
	Run: func(cmd *cobra.Command, args []string) {
		// init DB
		server.DB.Init(location, memory)

		// create the server
		S := &server.Server{
			Port: port,
			DB:   server.DB,
		}

		// TODO: check error
		_ = S.DB.Connect()

		done := make(chan os.Signal, 1)
		// Route shutdown signals to done channel
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		log.Println("TCP server started")
		go S.TCPStart()

		// Upon receiving a shutdown signal
		<-done
		fmt.Println("")
		log.Println("Shutting down server")
		err := S.DB.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {

}
