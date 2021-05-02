package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

//TCPStart starts a TCP server
func (s *Server) TCPStart() {
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

		go tHandle(conn)
	}
}

func tHandle(conn net.Conn) {

	log.Printf("Accepted connection from %v", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	defer conn.Close()

	for scanner.Scan() {
		ln := scanner.Text()
		resp := tCommand(ln)
		log.Printf("Host: %v Command: %v Response: %v", conn.RemoteAddr(), ln, resp)
		fmt.Fprintln(conn, resp)
	}

	log.Printf("Connection from %v closed\n", conn.RemoteAddr())
	return
}

func tCommand(s string) interface{} {
	e := "Invalid command"
	data := strings.Split(s, " ")
	l := len(data)

	if l < 2 {
		return e
	}

	switch strings.ToLower(data[0]) {

	case "get":
		res, _ := DB.Read(data[1])
		return res
	case "set":
		if l == 3 {
			err := DB.Create(data[1], data[2])
			if err != nil {
				return err
			}
			return fmt.Sprintf("Created new key %v", data[1])
		}
		return "Usage: [set] [key] [value]"
	case "upd":
		if l == 3 {
			if err := DB.Update(data[1], data[2]); err != nil {
				return err
			}
			return fmt.Sprintf("Created key: %v to value: %v", data[1], data[2])
		}
		return "Usage: [update] [key] [value]"
	case "del":
		if err := DB.Delete(data[1]); err != nil {
			return err
		}
		return fmt.Sprintf("Deleted key: %v", data[1])

	}

	return nil
}
