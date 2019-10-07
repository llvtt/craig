package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getListSearchesHandler(conf *CraigConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		searches, _ := json.Marshal(conf.Searches)
		w.Write(searches)
	}
}

func getRunSearchHandler(conf *CraigConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			methodNotFoundResponse(w)
			return
		}

		// todo respond with info about the search
		search(conf)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"status\":\"healthy\"}"))
}

func methodNotFoundResponse(w http.ResponseWriter) {
	http.Error(w, "{\"status\":\"error\", \"message\":\"No such http route\"}", http.StatusNotFound)
}

func initEndpoints(conf *CraigConfig) {
	http.HandleFunc("/searches", getListSearchesHandler(conf))
	http.HandleFunc("/search", getRunSearchHandler(conf))
	http.HandleFunc("/health", healthHandler)
}

func startServer(conf *CraigConfig) {
	log.Print("Starting server...")
	initEndpoints(conf)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
