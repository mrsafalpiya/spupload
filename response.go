package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"time"
)

type errorJson struct {
	Message string `json:"message"`
}

type fileCreatedJson struct {
	Url string `json:"url"`
}

type fileDetailsJson struct {
	Filename         string `json:"filename"`
	Size             int64  `json:"size"`
	ModificationTime string `json:"modification_time"`
}

func jsonErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	p := errorJson{Message: message}
	json.NewEncoder(w).Encode(p)
}

func jsonFileCreatedResponse(w http.ResponseWriter, url string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p := fileCreatedJson{url}
	json.NewEncoder(w).Encode(p)
}

func jsonFileNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	p := errorJson{Message: "Given file doesn't exist"}
	json.NewEncoder(w).Encode(p)
}

func jsonInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	p := errorJson{Message: "Internal "}
	json.NewEncoder(w).Encode(p)
}

func jsonFileDetails(w http.ResponseWriter, fileInfo fs.FileInfo) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p := fileDetailsJson{
		Filename:         fileInfo.Name(),
		Size:             fileInfo.Size(),
		ModificationTime: fileInfo.ModTime().UTC().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(p)
}
