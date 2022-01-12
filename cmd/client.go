package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var exec string

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

		if exec != "" {
			execCommands(conn)
			return
		}

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

func execCommands(conn *net.Conn) {
	scanner := bufio.NewScanner(*conn)
	cmds := strings.Split(exec, ";")

	for _, cmd := range cmds {
		fmt.Fprintln(*conn, cmd)
		ok := scanner.Scan()
		start := time.Now().Unix()
		for !ok {
			if time.Now().Unix()-start > 1 {
				fmt.Println("server timeout")
				return
			}
		}
		fmt.Println(scanner.Text())
	}
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringVarP(&host, "server", "s", "localhost", "Address of the server")
	clientCmd.Flags().IntVarP(&port, "port", "p", 9999, "Port on which the server is running")
	clientCmd.PersistentFlags().StringVarP(&exec, "command", "c", "", "directly exec command/s split by \";\"s")

}
