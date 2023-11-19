package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

const HOST = "http://localhost"
const PORT = 5432
const UPLOADS_DIR = "/home/safal/uploads"

func main() {
	router := httprouter.New()
	router.GET("/*filepath", serveFile())
	router.POST("/*filepath", uploadFile())

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), handler))
}

func serveFile() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fileReq := p.ByName("filepath")
		// disable / or /assets/childdirective/../
		if fileReq == "/" || (len(fileReq) > 1 && fileReq[len(fileReq)-1:] == "/") {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 page not found")
			return
		}

		switch viewMode := r.URL.Query().Get("view"); strings.ToLower(viewMode) {
		case "detail":
			fileLocationInDisk := filepath.Join(UPLOADS_DIR, fileReq)
			fileInfo, err := os.Stat(fileLocationInDisk)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					jsonFileNotFound(w)
					return
				} else {
					jsonErrorResponse(w, fmt.Sprintf("Internal server error: %s", err), http.StatusInternalServerError)
					return
				}
			}
			jsonFileDetails(w, fileLocationInDisk, fileInfo)
			return
		default:
			r.URL.Path = fileReq
			fileServer := http.FileServer(http.Dir(UPLOADS_DIR))
			fileServer.ServeHTTP(w, r)
		}
	}
}

func uploadFile() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		customFilename := r.PostFormValue("filename")
		replaceFile := r.PostFormValue("replace") == "true"
		disableFileOptimization := r.PostFormValue("disable-file-optimization") == "true"

		fileReq := p.ByName("filepath")
		uploadLocation := filepath.Join(UPLOADS_DIR, fileReq)

		fileUploaded, fileUploadedHeader, err := r.FormFile("file")
		if err != nil {
			jsonErrorResponse(w, fmt.Sprintf("Couldn't parse file from the request: %s", err), http.StatusInternalServerError)
			return
		}
		defer fileUploaded.Close()

		var outputFilepath string
		if customFilename != "" {
			fileExtension := filepath.Ext(fileUploadedHeader.Filename)
			outputFilepath = filepath.Join(uploadLocation, fmt.Sprintf("%s%s", customFilename, fileExtension))
		} else {
			outputFilepath = filepath.Join(uploadLocation, fileUploadedHeader.Filename)
		}

		outputFileDir := filepath.Dir(outputFilepath)
		err = os.MkdirAll(outputFileDir, os.ModePerm)
		if err != nil {
			jsonErrorResponse(w, fmt.Sprintf("Couldn't create output directory: %s", err), http.StatusInternalServerError)
			return
		}

		uploadedBuffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(uploadedBuffer, fileUploaded); err != nil {
			jsonErrorResponse(w, fmt.Sprintf("Internal server error: %s", err), http.StatusInternalServerError)
			return
		}

		if !disableFileOptimization {
			optimizeFile(&outputFilepath, uploadedBuffer)
		}

		if !replaceFile {
			outputFilepath = getProperAvailableFilepath(outputFilepath)
		}

		fileOut, err := os.OpenFile(outputFilepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			jsonErrorResponse(w, fmt.Sprintf("Couldn't create output file: %s", err), http.StatusInternalServerError)
			return
		}
		io.Copy(fileOut, uploadedBuffer)

		downloadUrl, err := filepathToDownloadUrl(outputFilepath)
		if err != nil {
			jsonErrorResponse(w, fmt.Sprintf("Couldn't generate download link: %s", err), http.StatusInternalServerError)
			return
		}

		jsonFileCreatedResponse(w, downloadUrl)
	}
}
