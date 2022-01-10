package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/server/tcp"
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
		// create the server
		errChan := make(chan error, 10)
		writeDone := make(chan bool)
		done := make(chan os.Signal, 1)

		s := tcp.New(
			port,
			database.New(location, memory, continousWrite, errChan, writeDone),
		)

		// Route shutdown signals to done channel
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		log.Println("TCP server started")
		go s.Serve()

		for {
			select {
			// Upon receiving a shutdown signal
			case <-done:
				fmt.Println("")
				log.Println("Shutting down server")
				go func() {
					err := s.Clean()
					if err != nil {
						log.Fatal(err)
					}
				}()
				<-writeDone
				return
			case err, ok := (<-errChan):
				log.Println("Error in write service:", err)
				if !ok {
					log.Println("Write service stopped")
				}
			default:
				time.Sleep(time.Millisecond * 500)
			}
		}
	},
}

func init() {
}
