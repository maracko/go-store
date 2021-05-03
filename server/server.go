package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/maracko/go-store/database"
	h "github.com/maracko/go-store/server/http"
	"github.com/maracko/go-store/server/tcp"
)

// Server is a struct with host info and a database instance
type Server struct {
	Port int
	DB   *database.DB
}

// HTTPStart starts the HTTP server
func (s *Server) HTTP() {
	// Map of all endpoints
	endpoints := map[string]http.HandlerFunc{
		// TODO: is path needed
		"/go-store": h.Redirect,
	}

	// Add middleware from []commonMiddleware to each endpoint
	for endpoint, f := range endpoints {
		http.HandleFunc(endpoint, multipleMiddleware(f, commonMiddleware...))
	}

	// If conn fails
	err := s.DB.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Write and close file on exit
	defer func() {
		// TODO: check error
		_ = s.DB.Disconnect()
	}()

	// TODO: check error
	_ = http.ListenAndServe(fmt.Sprintf(":%v", s.Port), nil)
}

// TCPStart starts a TCP server
func (s *Server) TCP() {
	li, err := net.Listen("tcp", fmt.Sprintf(":%v", s.Port))
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

		go tcp.Handle(conn)
	}
}
