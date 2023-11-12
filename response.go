package main

import (
	"encoding/json"
	"net/http"
)

type errorJson struct {
	Message string `json:"message"`
}

type fileCreatedJson struct {
	Url string `json:"url"`
}

func jsonErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	p := errorJson{message}
	json.NewEncoder(w).Encode(p)
}

func jsonFileCreatedResponse(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p := fileCreatedJson{url}
	json.NewEncoder(w).Encode(p)
}
