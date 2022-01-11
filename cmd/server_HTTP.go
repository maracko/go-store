package cmd

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/server/http"
	"github.com/spf13/cobra"
)

// serveHTTPCmd represents the server HTTP command
var serveHTTPCmd = &cobra.Command{
	Use:     "HTTP",
	Aliases: []string{"http"},
	Short:   "Start HTTP server",
	Long: `Starts HTTP server.
	Defaults to port 8888. Empty database will be initialized and kept only in memory if no path is provided.
	If you have json file with data you want to read from but not save provide a location along with -m (memory) flag`,
	Run: func(cmd *cobra.Command, args []string) {
		// create the server
		errChan := make(chan error, 5)
		writeSvcDone := make(chan bool)
		srvDone := &sync.WaitGroup{}
		done := make(chan os.Signal, 1)

		s := http.New(
			port,
			tlsPort,
			token,
			pKey,
			cert,
			database.New(location, memory, continousWrite, errChan, writeSvcDone),
			srvDone,
		)

		// Route shutdown signals to done channel
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		if location == "" {
			log.Println("Starting a blank database")
		}

		s.Serve()
		srvDone.Add(1)
		if pKey != "" && cert != "" {
			srvDone.Add(1)
		}

		for {
			select {
			// Upon receiving a shutdown signal
			case <-done:
				log.Println("Shutting down server")
				if err := s.Clean(); err != nil {
					log.Fatalln("Dirty shutdown:", err)
				}
				return
			case err := (<-errChan):
				log.Println("Error in write service:", err)
			}
		}
	},
}
var cert string
var pKey string
var tlsPort int
var token string

func init() {
	serveHTTPCmd.PersistentFlags().StringVar(&cert, "certificate", "", "Certificate location for your server. If signed by CA, must concatenate thems")
	serveHTTPCmd.PersistentFlags().StringVar(&pKey, "private-key", "", "Your private key location")
	serveHTTPCmd.PersistentFlags().IntVar(&tlsPort, "tls-port", 9999, "Port on which HTTPS will be served")
	serveHTTPCmd.PersistentFlags().StringVarP(&token, "token", "t", "", `Auth. key. It needs to be sent in the "Authorization" header on every request`)
}
