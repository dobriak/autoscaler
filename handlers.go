package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

type Message struct {
	Message string `json:"message"`
}

func JsonResponse(w http.ResponseWriter, message string) {
	response := Message{message}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func RemoveApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app_id := vars["appid"]
	if err := RepoRemoveApp(app_id); err != nil {
		w.WriteHeader(400)
		panic(err)
	}
	JsonResponse(w, "OK")
}

func AddApp(w http.ResponseWriter, r *http.Request) {
	var app App
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &app); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	RepoAddApp(app)

	w.WriteHeader(200)
	JsonResponse(w, "OK")
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	app_id := vars["appid"]
	app := RepoFindApp(app_id)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(app); err != nil {
		panic(err)
	}
}

func IndexApps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(apps); err != nil {
		panic(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	JsonResponse(w, "Autoscaler, v0.0.1")
}
