package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/maracko/go-store/database"
)

// New create new server
func New(port int, db *database.DB) *tcpServer {
	s := &tcpServer{
		port: port,
		db:   db,
	}

	if err := s.db.Connect(); err != nil {
		log.Fatalln(err)
	}

	return s
}

// Server is a struct with host info and a database instance
type tcpServer struct {
	port int
	db   *database.DB
}

// Clean cleans a server
func (s *tcpServer) Clean() error {
	return s.db.Disconnect()
}

// Serve starts a TCP server
func (s *tcpServer) Serve() {
	li, err := net.Listen("tcp", fmt.Sprintf(":%v", s.port))
	if err != nil {
		log.Fatalln(err)
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.handle(conn)
	}
}
func (s *tcpServer) handle(conn net.Conn) {
	log.Printf("Accepted connection from %v", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	defer conn.Close()

	for scanner.Scan() {
		ln := scanner.Text()
		resp := s.command(ln)
		log.Printf("Host: %v Command: %v Response: %v", conn.RemoteAddr(), ln, resp)
		fmt.Fprintln(conn, resp)
	}

	log.Printf("Connection from %v closed\n", conn.RemoteAddr())
}

func (s *tcpServer) command(input string) interface{} {
	e := "Invalid command"
	data := strings.Split(input, " ")
	l := len(data)

	if l < 2 {
		return e
	}

	switch strings.ToLower(data[0]) {
	case "get":
		res, _ := s.db.Read(data[1])
		return res
	case "set":
		if l == 3 {
			err := s.db.Create(data[1], data[2])
			if err != nil {
				return err
			}
			return fmt.Sprintf("created %v", data[1])
		}
		return "usage: [set] [key] [value]"
	case "upd":
		if l == 3 {
			if err := s.db.Update(data[1], data[2]); err != nil {
				return err
			}
			return fmt.Sprintf("updated %v", data[1])
		}
		return "usage: [update] [key] [value]"
	case "del":
		if err := s.db.Delete(data[1]); err != nil {
			return err
		}
		return fmt.Sprintf("deleted %v", data[1])
	}

	return nil
}
