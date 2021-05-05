package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/errors"
	"github.com/maracko/go-store/server"
	"github.com/maracko/go-store/server/http/helpers"
)

type resource struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// New create new server
func New(port int, db *database.DB) server.Server {
	return &httpServer{
		port: port,
		db:   db,
	}
}

// Server is a struct with host info and a database instance
type httpServer struct {
	port int
	db   *database.DB
}

// Clean clean server
func (s *httpServer) Clean() error {
	return s.db.Disconnect()
}

// Serve starts the HTTP server
func (s *httpServer) Serve() {
	// Map of all endpoints
	endpoints := map[string]http.HandlerFunc{
		// TODO: is path needed
		"/go-store": s.handle,
	}

	// Add middleware from []commonMiddleware to each endpoint
	for endpoint, f := range endpoints {
		http.HandleFunc(endpoint, multipleMiddleware(f, commonMiddleware...))
	}

	// If conn fails
	err := s.db.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Write and close file on exit
	defer func() {
		// TODO: check error
		_ = s.db.Disconnect()
	}()

	// TODO: check error
	_ = http.ListenAndServe(fmt.Sprintf(":%v", s.port), nil)
}

// Handle appropriate func based on method and params
func (s *httpServer) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		key, ok := r.URL.Query()["key"]

		if !ok || len(key[0]) < 1 {
			helpers.JSONEncode(w, errors.BadRequest("missing key"))
			return
		}

		keys := strings.Split(key[0], ",")

		if len(keys) == 1 {
			s.read(w, r)
			return
		}

		s.readMany(w, r)
		return

	case "POST":
		s.create(w, r)
	case "PATCH":
		s.update(w, r)
	case "DELETE":
		s.delete(w, r)
	default:
		helpers.JSONEncode(w, errors.MethodNotAllowed("method %s not allowed", r.Method))
	}
}

// Read read database key
func (s *httpServer) read(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()["key"]

	val, err := s.db.Read(key[0])
	if err != nil {
		helpers.JSONEncode(w, errors.NotFoundWrap(err, "not found"))
		return
	}

	// TODO: check error
	_ = json.NewEncoder(w).Encode(resource{key[0], val})
}

// ReadMany read many records
func (s *httpServer) readMany(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()["key"]

	keys := strings.Split(key[0], ",")

	empty := true
	for _, k := range keys {
		if k != "" {
			empty = false
		}
	}

	// If all keys are empty return
	if empty == true {
		helpers.JSONEncode(w, errors.NotFound("all keys are empty"))
		return
	}

	val := s.db.ReadMany(keys...)
	// TODO: check error
	_ = json.NewEncoder(w).Encode(val)
}

// Create create new value
func (s *httpServer) create(w http.ResponseWriter, r *http.Request) {
	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		helpers.JSONEncode(w, errors.InternalWrap(err, "unmarshal error"))
		return
	}

	if err := s.db.Create(res.Key, res.Value); err != nil {
		helpers.JSONEncode(w, errors.BadRequestWrap(err, "duplicate key"))
		return
	}

	helpers.JSONEncode(w, resource{res.Key, res.Value})
}

// Update update key
func (s *httpServer) update(w http.ResponseWriter, r *http.Request) {
	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		helpers.JSONEncode(w, errors.InternalWrap(err, "unmarshal error"))
		return
	}

	if err := s.db.Update(res.Key, res.Value); err != nil {
		helpers.JSONEncode(w, err)
		return
	}

	helpers.JSONEncode(w, resource{res.Key, res.Value})
}

// Delete delete key
func (s *httpServer) delete(w http.ResponseWriter, r *http.Request) {
	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		helpers.JSONEncode(w, errors.InternalWrap(err, "unmarshal error"))
		return
	}

	if err := s.db.Delete(res.Key); err != nil {
		helpers.JSONEncode(w, err)
		return
	}

	val := make(map[string]interface{})
	val["deleted"] = true
	helpers.JSONEncode(w, resource{res.Key, val})
}
