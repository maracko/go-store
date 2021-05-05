package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-store",
	Short: "Go store is a key/value data store",
	Long: `Go store creates an in memory database similar to redis.
	You can start TCP or HTTP server and access them with a CLI or HTTP requests.
	After server shutdown a json encoded file with your data will be written.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

}
