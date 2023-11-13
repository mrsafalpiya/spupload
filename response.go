package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"time"

	"github.com/gabriel-vasile/mimetype"
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
	Filetype         string `json:"filetype"`
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

func jsonFileDetails(w http.ResponseWriter, path string, fileInfo fs.FileInfo) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var filetype string
	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		filetype = "Unknown"
	} else {
		filetype = mtype.String()
	}

	p := fileDetailsJson{
		Filename:         fileInfo.Name(),
		Size:             fileInfo.Size(),
		ModificationTime: fileInfo.ModTime().UTC().Format(time.RFC3339),
		Filetype:         filetype,
	}
	json.NewEncoder(w).Encode(p)
}
