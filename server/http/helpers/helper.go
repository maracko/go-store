package helpers

import (
	"encoding/json"
	"log"
	"net/http"

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
		w.Write(b)
		return
	}

	err := json.NewEncoder(w).Encode(b)
	if err != nil {
		log.Println(err)
	}
}
