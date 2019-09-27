package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func listSearchesHandler(w http.ResponseWriter, r *http.Request) {
	config := parseConfig(CONFIG_FILE_NAME)
	searches, _ := json.Marshal(config.Searches)
	w.Write(searches)
}

func runSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		methodNotFoundResponse(w)
	}

	search()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"status\":\"healthy\"}"))
}

func methodNotFoundResponse(w http.ResponseWriter) {
	http.Error(w, "{\"status\":\"error\", \"message\":\"No such http route\"}", http.StatusNotFound)
}

func initEndpoints() {
	http.HandleFunc("/searches", listSearchesHandler)
	http.HandleFunc("/search", runSearchHandler)
	http.HandleFunc("/health", healthHandler)
}

func startServer() {
	log.Print("Starting server...")
	initEndpoints()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
