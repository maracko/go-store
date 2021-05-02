package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start client connection over TCP",
	Long:  `Starts a client connected to the server on the provided port`,
	Run: func(cmd *cobra.Command, args []string) {

		conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))

		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()

		fmt.Println("Welcome to go-store server!")
		log.Printf("Connected to %v:%v", host, port)
		scanner := bufio.NewScanner(conn)
		reader := bufio.NewReader(os.Stdin)
		for {

			fmt.Print("$:")
			b, _, _ := reader.ReadLine()
			str := string(b)
			fmt.Fprintln(conn, str)
			ok := scanner.Scan()

			if !ok {
				log.Println("Connection closed by remote host")
				return
			}

			fmt.Println(scanner.Text())
		}

	},
}

func init() {

	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&host, "server", "s", "localhost", "Address of the server")
	clientCmd.Flags().IntVarP(&port, "port", "p", 9999, "Port on which the server is running")
}
