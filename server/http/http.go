package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/maracko/go-store/database"
	"github.com/maracko/go-store/errors"
	"github.com/maracko/go-store/server/http/helpers"
)

type resource struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// DB is the package wide pointer to a database object used for crud operations, it must be initialized first
var DB = &database.DB{}

// Redirect appropriate func based on method and params
func Redirect(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		key, ok := r.URL.Query()["key"]

		if !ok || len(key[0]) < 1 {
			helpers.JSONEncode(w, errors.BadRequest("missing key"))
			return
		}

		keys := strings.Split(key[0], ",")

		if len(keys) == 1 {
			Read(w, r)
			return
		}

		ReadMany(w, r)
		return

	case "POST":
		Create(w, r)
	case "PATCH":
		Update(w, r)
	case "DELETE":
		Delete(w, r)
	default:
		helpers.JSONEncode(w, errors.MethodNotAllowed("method %s not allowed", r.Method))
	}
}

// Read read database key
func Read(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query()["key"]

	val, err := DB.Read(key[0])
	if err != nil {
		helpers.JSONEncode(w, errors.NotFoundWrap(err, "not found"))
		return
	}

	// TODO: check error
	_ = json.NewEncoder(w).Encode(resource{key[0], val})
}

// ReadMany read many records
func ReadMany(w http.ResponseWriter, r *http.Request) {
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

	val := DB.ReadMany(keys...)
	// TODO: check error
	_ = json.NewEncoder(w).Encode(val)
}

// Create create new value
func Create(w http.ResponseWriter, r *http.Request) {
	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		helpers.JSONEncode(w, errors.InternalWrap(err, "unmarshal error"))
		return
	}

	if err := DB.Create(res.Key, res.Value); err != nil {
		helpers.JSONEncode(w, errors.BadRequestWrap(err, "duplicate key"))
		return
	}

	helpers.JSONEncode(w, resource{res.Key, res.Value})
}

// Update update key
func Update(w http.ResponseWriter, r *http.Request) {
	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		helpers.JSONEncode(w, errors.InternalWrap(err, "unmarshal error"))
		return
	}

	if err := DB.Update(res.Key, res.Value); err != nil {
		helpers.JSONEncode(w, err)
		return
	}

	helpers.JSONEncode(w, resource{res.Key, res.Value})
}

// Delete delete key
func Delete(w http.ResponseWriter, r *http.Request) {
	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		helpers.JSONEncode(w, errors.InternalWrap(err, "unmarshal error"))
		return
	}

	if err := DB.Delete(res.Key); err != nil {
		helpers.JSONEncode(w, err)
		return
	}

	val := make(map[string]interface{})
	val["deleted"] = true
	helpers.JSONEncode(w, resource{res.Key, val})
}
