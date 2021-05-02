package server

import (
	"encoding/json"
	"net/http"
)

//jsonErr writes err to w in JSON format
func jsonErr(w http.ResponseWriter, err errD, code int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

//jsonH sets Content-type header of w to application/json
func jsonH(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
}

//jsonEnc encodes b to JSON and writes to w
func jsonEnc(w http.ResponseWriter, b interface{}) {
	json.NewEncoder(w).Encode(b)
}
