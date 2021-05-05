package server

import "github.com/maracko/go-store/database"

// Server is a struct with host info and a database instance
type Server struct {
	Port int
	DB   *database.DB
}

// DB is the package wide pointer to a database object used for crud operations, it must be initialized first
var DB = &database.DB{}

type errD struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type resource struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
