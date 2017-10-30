package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Message struct to be displayed as json
type Message struct {
	Message string `json:"message"`
}

//JSONResponse write generic message as json response
func JSONResponse(w http.ResponseWriter, message string) {
	response := Message{message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Panicln(err)
	}
}

//RemoveApp removes app by its ID from the pool of apps being monitored
func RemoveApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID := vars["appid"]
	if err := RepoRemoveApp(appID); err != nil {
		w.WriteHeader(400)
		log.Panicln(err)
	}
	JSONResponse(w, "OK")
}

//AddApp adds an app to the pool of monitored apps
func AddApp(w http.ResponseWriter, r *http.Request) {
	var app App
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Panicln(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Panicln(err)
	}
	if err := json.Unmarshal(body, &app); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Panicln(err)
		}
	}
	RepoAddApp(app)

	w.WriteHeader(200)
	JSONResponse(w, "OK")
}

//GetApp finds and displays a monitored app by its ID
func GetApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appID := vars["appid"]
	app := RepoFindApp(appID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(app); err != nil {
		log.Panicln(err)
	}
}

//IndexApps displays all monitored apps
func IndexApps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(apps); err != nil {
		log.Panicln(err)
	}
}

//Index for slash, returns version
func Index(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, "Autoscaler, v0.0.1")
}
