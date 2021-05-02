package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//HTTPStart starts the HTTP server
func (s *Server) HTTPStart() {

	//Map of all endpoints
	endpoints := map[string]http.HandlerFunc{
		"/go-store": hRedirect,
	}

	//Add middleware from []commonMiddleware to each endpoint
	for endpoint, f := range endpoints {
		http.HandleFunc(endpoint, multipleMiddleware(f, commonMiddleware...))
	}

	//If conn fails
	err := DB.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	//Write and close file on exit
	defer s.DB.Disconnect()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", s.Port), nil))

}

//Choose appropriate func based on method and params
func hRedirect(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		key, ok := r.URL.Query()["key"]

		if !ok || len(key[0]) < 1 {
			jsonEnc(w, errD{"missing key", "no key param found"})
			return
		}

		keys := strings.Split(key[0], ",")

		if len(keys) == 1 {
			hRead(w, r)
			return
		}

		hReadMany(w, r)
		return

	case "POST":
		hCreate(w, r)
	case "PATCH":
		hUpdate(w, r)
	case "DELETE":
		hDel(w, r)
	default:
		jsonErr(w, errD{"Not allowed", "Method not allowed"}, http.StatusMethodNotAllowed)
	}
	return
}

func hRead(w http.ResponseWriter, r *http.Request) {

	key, _ := r.URL.Query()["key"]

	val, err := DB.Read(key[0])

	if err != nil {
		jsonErr(w, errD{"not found", err.Error()}, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resource{key[0], val})
	return
}

func hReadMany(w http.ResponseWriter, r *http.Request) {

	key, _ := r.URL.Query()["key"]

	keys := strings.Split(key[0], ",")

	empty := true
	for _, k := range keys {
		if k != "" {
			empty = false
		}
	}

	//If all keys are empty return
	if empty == true {
		jsonErr(w, errD{"no keys provided", "all keys are empty"}, 404)
		return
	}

	val := DB.ReadMany(keys...)

	json.NewEncoder(w).Encode(val)
	return
}

func hCreate(w http.ResponseWriter, r *http.Request) {

	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		jsonErr(w, errD{"Unknown", err.Error()}, 500)
		return
	}

	if err := DB.Create(res.Key, res.Value); err != nil {
		jsonErr(w, errD{"duplicate key", err.Error()}, http.StatusConflict)
		return
	}

	jsonEnc(w, resource{res.Key, res.Value})
}

func hUpdate(w http.ResponseWriter, r *http.Request) {

	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		jsonErr(w, errD{"Unknown", err.Error()}, 500)
		return
	}

	if err := DB.Update(res.Key, res.Value); err != nil {
		jsonErr(w, errD{"not found", err.Error()}, 404)
		return
	}

	jsonEnc(w, resource{res.Key, res.Value})
}

func hDel(w http.ResponseWriter, r *http.Request) {

	var res resource
	b, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		jsonErr(w, errD{"Unknown", err.Error()}, 500)
		return
	}

	if err := DB.Delete(res.Key); err != nil {
		jsonErr(w, errD{"not found", err.Error()}, 404)
		return
	}

	val := make(map[string]interface{})
	val["deleted"] = true
	jsonEnc(w, resource{res.Key, val})
}
