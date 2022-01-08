package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/maracko/go-store/errors"
)

// JSONEncode encodes b to JSON and writes to w
func JSONEncode(w http.ResponseWriter, b interface{}) {
	w.Header().Set("Content-type", "application/json")

	e, ok := b.(errors.Error)
	if ok {
		b, err := json.Marshal(e)
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(e.Status)
		_, _ = w.Write(b)
		return
	}

	err := json.NewEncoder(w).Encode(b)
	if err != nil {
		log.Println(err)
	}
}

//Return an array of keys(everything after / in url path),
//they are strings split by a comma
func ExtractKeys(r *http.Request) []string {
	key := strings.TrimPrefix(r.URL.Path, "/")
	keys := strings.Split(key, ",")
	return keys
}
