package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func getProperAvailableFilepath(path string) string {
	filepathExtension := filepath.Ext(path)
	filepathBase := strings.TrimSuffix(path, filepathExtension)

	properAvailableFilepath := path
	curRevision := 1
	for true {
		if _, err := os.Stat(properAvailableFilepath); err != nil {
			break
		}
		properAvailableFilepath = fmt.Sprintf("%s-%d%s", filepathBase, curRevision, filepathExtension)
		curRevision++
	}
	return properAvailableFilepath
}

func filepathToDownloadUrl(path string) (string, error) {
	return url.JoinPath(getFullHostname(), strings.TrimPrefix(path, UPLOADS_DIR))
}

func getFullHostname() string {
	if PORT == 80 || PORT == 443 {
		return HOST
	}
	return fmt.Sprintf("%s:%d", HOST, PORT)
}
