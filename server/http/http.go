package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/errors"
	"github.com/maracko/go-store/server/http/helpers"
)

type resource struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var key string

func init() {
	rand.Seed(time.Now().UnixNano())
}

// New create new server
func New(port, tlsPort int, token, pKey, cert string, db *database.DB, wg *sync.WaitGroup) *httpServer {
	srv := &http.Server{
		Addr: ":" + fmt.Sprint(port),
	}
	return &httpServer{
		port:  port,
		token: token,
		pKey:  pKey,
		cert:  cert,
		db:    db,
		srv:   srv,
		wg:    wg,
	}
}

// Server is a struct with host info and a database instance
type httpServer struct {
	port              int
	tlsPort           int
	token, pKey, cert string
	wg                *sync.WaitGroup
	db                *database.DB
	srv               *http.Server
}

// Clean clean server
func (s *httpServer) Clean() error {

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(time.Millisecond * 1500)
		cancel()
	}()

	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Println("server error:", err)
	}
	log.Println("HTTP/S shut down")

	return s.db.Disconnect()
}

// Serve starts the HTTP server
func (s *httpServer) Serve() {
	key = s.token

	// Map of all endpoints
	endpoints := map[string]http.HandlerFunc{
		"/": s.handle,
	}

	// Add middleware from []commonMiddleware to each endpoint
	for endpoint, f := range endpoints {
		http.HandleFunc(endpoint, multipleMiddleware(f, commonMiddleware...))
	}

	err := s.db.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	if s.pKey != "" && s.cert != "" {
		go func() {
			// let main know we are done cleaning up
			defer s.wg.Done()
			s.srv.Addr = ":" + fmt.Sprint(s.tlsPort)
			if err := s.srv.ListenAndServeTLS(s.cert, s.pKey); err != http.ErrServerClosed {
				log.Println("TLS error:", err)
			}
		}()
		log.Println("HTTPS server started")
	}

	go func() {
		defer s.wg.Done()
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(err)
		}
	}()
}

// Handle appropriate func based on method and params
func (s *httpServer) handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		keys := helpers.ExtractKeys(r)

		switch {
		case len(keys) == 0:
			helpers.JSONEncode(w, errors.BadRequest("missing key"))
			return
		case len(keys) == 1:
			s.read(w, r)
			return
		case len(keys) > 1:
			s.readMany(w, r)
			return
		}

	case "POST":
		s.create(w, r)
	case "PATCH":
		s.update(w, r)
	case "DELETE":
		keys := helpers.ExtractKeys(r)

		switch {
		case len(keys) == 0:
			helpers.JSONEncode(w, errors.BadRequest("missing key"))
			return
		case len(keys) == 1:
			s.delete(w, r)
			return
		case len(keys) > 1:
			s.deleteMany(w, r)
		}

	default:
		helpers.JSONEncode(w, errors.MethodNotAllowed("method %s not allowed", r.Method))
	}
}

// Read read database key
func (s *httpServer) read(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/")

	val, err := s.db.Read(key)
	if err != nil {
		helpers.JSONEncode(w, errors.NotFoundWrap(err, "not found"))
		return
	}

	helpers.JSONEncode(w, resource{key, val})
}

// ReadMany read many records
func (s *httpServer) readMany(w http.ResponseWriter, r *http.Request) {

	keys := helpers.ExtractKeys(r)

	empty := true
	for _, k := range keys {
		if k != "" {
			empty = false
			break
		}
	}

	if empty {
		helpers.JSONEncode(w, errors.NotFound("all keys are empty"))
		return
	}

	res := s.db.ReadMany(keys...)
	resp := []resource{}

	for k, v := range res {
		resp = append(resp, resource{k, v})
	}
	helpers.JSONEncode(w, resp)
}

// Create create new value
func (s *httpServer) create(w http.ResponseWriter, r *http.Request) {
	var res resource
	var multiRes []resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		err = json.Unmarshal(b, &multiRes)
		if err != nil {
			helpers.JSONEncode(w, errors.InternalWrap(err, "unmarshal error"))
			return
		}
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
		helpers.JSONEncode(w, errors.BadRequestWrap(err, "update error"))
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
		helpers.JSONEncode(w, errors.NotFoundWrap(err, "delete error"))
		return
	}

	del := make(map[string]bool)
	del["deleted"] = true
	helpers.JSONEncode(w, resource{res.Key, del})
}

func (s *httpServer) deleteMany(w http.ResponseWriter, r *http.Request) {

	keys := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), ",")

	res := s.db.DeleteMany(keys...)

	delFlag := false
	errFlag := false
	// Set flags so we can show appropriate status code
	for _, v := range res {
		err, ok := v.(map[string]string)
		if ok {
			if _, ok := err["error"]; ok {
				errFlag = true
			}
		}

		del, ok := v.(map[string]bool)
		if ok {
			if _, ok := del["deleted"]; ok {
				delFlag = true
			}
		}
	}

	if errFlag && delFlag {
		w.WriteHeader(http.StatusMultiStatus)
	} else if errFlag && !delFlag {
		w.WriteHeader(http.StatusNotFound)
	}

	resp := []resource{}
	for k, v := range res {
		resp = append(resp, resource{k, v})
	}
	helpers.JSONEncode(w, resp)
}
