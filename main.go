package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

const HOST = "http://localhost"
const PORT = 5555
const UPLOADS_DIR = "/home/safal/uploads"

func main() {
	router := httprouter.New()
	router.GET("/*filepath", serveFile())
	router.POST("/*filepath", uploadFile())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), router))
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

		r.URL.Path = fileReq
		fileServer := http.FileServer(http.Dir(UPLOADS_DIR))
		fileServer.ServeHTTP(w, r)
	}
}

func uploadFile() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		customFilename := r.PostFormValue("filename")
		replaceFile := r.PostFormValue("replace") == "true"

		fileReq := p.ByName("filepath")
		uploadLocation := filepath.Join(UPLOADS_DIR, fileReq)

		fileUploaded, fileUploadedHeader, err := r.FormFile("file")
		if err != nil {
			jsonErrorResponse(w, fmt.Sprintf("Couldn't parse file from the request: %s", err), http.StatusInternalServerError)
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
		}

		if !replaceFile {
			outputFilepath = getProperAvailableFilepath(outputFilepath)
		}
		fileOut, err := os.OpenFile(outputFilepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			jsonErrorResponse(w, fmt.Sprintf("Couldn't create output file: %s", err), http.StatusInternalServerError)
		}
		io.Copy(fileOut, fileUploaded)

		downloadUrl, err := filepathToDownloadUrl(outputFilepath)
		if err != nil {
			jsonErrorResponse(w, fmt.Sprintf("Couldn't generate download link: %s", err), http.StatusInternalServerError)
		}

		jsonFileCreatedResponse(w, downloadUrl)
	}
}
