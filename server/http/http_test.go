package http

import (
	"os"
	"sync"
	"testing"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/database/helpers"
)

var port int
var db *database.DB
var srv *httpServer
var path string

func init() {
	port = 8888
	tlsPort := 9999
	errChan := make(chan error, 10)
	dc := make(chan bool)
	path = ".test.file"

	db = database.New(path, false, true, errChan, dc, 1)
	s := New(port, tlsPort, "", "", "", db, &sync.WaitGroup{})
	srv = s.(*httpServer)
}

func TestConnect(t *testing.T) {
	if err := srv.db.Connect(); err != nil {
		t.Errorf("db connection failed: %s", err)
	}
	if err := srv.db.Disconnect(); err != nil {
		t.Errorf("db disconnect failed: %s", err)
	}
}

func TestClean(t *testing.T) {
	if e := srv.Clean(); e != nil {
		t.Error("{}", e)
	}

	if helpers.FileExists(path) {
		err := os.Remove(path)
		if err != nil {
			t.Error("delete failed:", err)
		}
	}

}
