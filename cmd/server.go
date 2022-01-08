package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a new server",
	Long: `Starts a new server, 2 types are supported:
	server HTTP -> Serves over HTTP (check serve HTTP help for more info)
	server TCP -> Serves over custom TCP protocol, interaction is done with the 'go-store cli' command  (check serve TCP help for more info)`,
	ValidArgs: []string{"HTTP", "TCP"},
	// Enforce argument constraints
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.OnlyValidArgs(cmd, args); err != nil {
			return err
		}

		exact := cobra.ExactArgs(1)
		if err := exact(cmd, args); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "HTTP":
			_, err := serveHTTPCmd.ExecuteC()
			if err != nil {
				log.Fatalln(err)
			}
		case "TCP":
			_, err := serveTCPCmd.ExecuteC()
			if err != nil {
				log.Fatalln(err)
			}
		}
	},
}

var port int
var host string
var location string
var memory bool
var continousWrite bool

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serveHTTPCmd)
	serverCmd.AddCommand(serveTCPCmd)

	serverCmd.PersistentFlags().IntVarP(&port, "port", "p", 8888, "Port on which to start the server")
	serverCmd.PersistentFlags().StringVarP(&location, "location", "l", "", "Location of the database file. If empty all changes will be lost upon server shutdown")
	serverCmd.PersistentFlags().BoolVarP(&memory, "memory", "m", false, "If present values won't be saved upon exit (Has no effect if location is empty)")
	serverCmd.PersistentFlags().BoolVarP(&continousWrite, "continous-write", "c", false, "Keep writing data to file on every update asynchronously")
}
