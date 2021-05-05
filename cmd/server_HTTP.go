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
		s := http.New(
			port,
			database.New(location, memory),
		)

		done := make(chan os.Signal, 1)
		// Route shutdown signals to done channel
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		log.Println("HTTP server started")
		go s.Serve()

		// Upon receiving a shutdown signal
		<-done
		fmt.Println("")
		log.Println("Shutting down server")
		err := s.Clean()
		if err != nil {
			log.Fatal(err)
		}
	},
}
