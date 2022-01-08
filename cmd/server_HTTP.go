package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/server/http"
	"github.com/spf13/cobra"
)

// serveHTTPCmd represents the server HTTP command
var serveHTTPCmd = &cobra.Command{
	Use:   "HTTP",
	Short: "Start HTTP server",
	Long: `Starts HTTP server.
	Defaults to port 8888. Empty database will be initialized and kept only in memory if no path is provided.
	If you have json file with data you want to read from but not save provide a location along with -m (memory) flag`,
	Run: func(cmd *cobra.Command, args []string) {
		// create the server
		errChan := make(chan error, 5)
		s := http.New(
			port,
			database.New(location, memory, continousWrite, errChan),
			errChan,
		)

		done := make(chan os.Signal, 1)
		// Route shutdown signals to done channel
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		log.Printf("HTTP server started on port %d", port)
		go s.Serve()

		for {
			select {
			// Upon receiving a shutdown signal
			case <-done:
				fmt.Println("")
				log.Println("Shutting down server")
				err := s.Clean()
				if err != nil {
					log.Fatal(err)
				}
				return
			case err, ok := (<-errChan):
				if !ok {
					log.Println("Exiting")
					return
				}
				log.Println("Error in write service:", err)
			}
		}
	},
}
