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
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

var HOST string
var PORT int64
var UPLOADS_DIR string
var API_KEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[ERROR] Couldn't load the .env file: %s", err)
	}

	HOST = os.Getenv("HOST")
	if HOST == "" {
		log.Fatal("[ERROR] The environment variable 'HOST' is not set")
	}

	portStr := os.Getenv("PORT")
	if HOST == "" {
		log.Fatal("[ERROR] The environment variable 'PORT' is not set")
	}
	PORT, err = strconv.ParseInt(portStr, 10, 32)
	if portStr == "" || err != nil {
		log.Fatal("[ERROR] The environment variable 'PORT' is not a number")
	}

	UPLOADS_DIR = os.Getenv("UPLOADS_DIR")
	if UPLOADS_DIR == "" {
		log.Fatal("[ERROR] The environment variable 'UPLOADS_DIR' is not set")
	}

	API_KEY = os.Getenv("API_KEY")
	if API_KEY == "" {
		log.Fatal("[ERROR] The environment variable 'API_KEY' is not set")
	}
}

func main() {
	router := httprouter.New()
	router.GET("/*filepath", serveFile())
	router.POST("/*filepath", uploadFile())

	handler := cors.Default().Handler(router)
	log.Printf("Listening on port %d\n", PORT)
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
		headerAPIKey := r.Header.Get("x-spupload-api-key")
		if headerAPIKey != API_KEY {
			jsonErrorResponse(w, fmt.Sprintf("Invalid 'x-spupload-api-key' header value"), http.StatusBadRequest)
			return
		}

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
