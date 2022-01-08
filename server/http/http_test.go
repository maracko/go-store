package http

import (
	"os"
	"testing"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/database/helpers"
)

var port int
var db *database.DB
var s httpServer
var path string

func init() {
	port = 8888
	errChan := make(chan error, 10)
	path = ".test.file"

	db = database.New(path, false, true, errChan)
	s = *New(port, db, errChan)
}

func TestConnect(t *testing.T) {
	if err := s.db.Connect(); err != nil {
		t.Errorf("db connection failed: %s", err)
	}
}

func TestClean(t *testing.T) {
	if e := s.Clean(); e != nil {
		t.Error("{}", e)
	}

	if helpers.FileExists(path) {
		err := os.Remove(path)
		if err != nil {
			t.Error("delete failed:", err)
		}
	}

}
