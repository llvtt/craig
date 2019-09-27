package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func listSearches(w http.ResponseWriter, r *http.Request) {
	config := parseConfig(CONFIG_FILE_NAME)
	searches, _ := json.Marshal(config.Searches)
	w.Write(searches)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"status\":\"healthy\"}"))
}

func initEndpoints() {
	http.HandleFunc("/searches", listSearches)
	http.HandleFunc("/health", health)
}

func startServer() {
	initEndpoints()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
