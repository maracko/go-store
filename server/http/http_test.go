package http

import (
	"testing"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/server"
)

var port int
var db *database.DB
var s server.Server

func init() {
	port = 8888
	errChan := make(chan error, 10)
	db = database.New(".test.file", false, true, errChan)
	s = New(port, db, errChan)

}

func TestNew(t *testing.T) {
	if _, ok := s.(*httpServer); !ok {
		t.Error("Failed to init server on port {} with empty db", port)
	}
}

func TestClean(t *testing.T) {
	if e := s.Clean(); e != nil {
		t.Error("{}", e)
	}
}
