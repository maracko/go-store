package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/maracko/go-store/errors"
)

// JSONEncode encodes b to JSON and writes to w
func JSONEncode(w http.ResponseWriter, b interface{}) {
	w.Header().Set("Content-type", "application/json")

	e, ok := b.(errors.Error)
	if ok {
		b, _ := json.Marshal(e)
		w.WriteHeader(e.Status)
		w.Write(b)
		return
	}

	// TODO: check error
	_ = json.NewEncoder(w).Encode(b)
}
