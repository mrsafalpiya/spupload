package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

func optimizeFile(outputFilepath *string, fileBytes *bytes.Buffer) {
	mtype := mimetype.Detect(fileBytes.Bytes())

	mtypeString := mtype.String()
	switch {
	case strings.HasPrefix(mtypeString, "image/"):
		if mtypeString == "image/webp" {
			return
		}

		var img image.Image
		var err error
		switch {
		case strings.HasSuffix(mtypeString, "jpeg"):
			img, err = jpeg.Decode(fileBytes)
		case strings.HasSuffix(mtypeString, "png"):
			img, err = png.Decode(fileBytes)
		default:
			return // Unsupported image format
		}
		if err != nil {
			return
		}

		webpOptions, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
		if err != nil {
			return
		}
		if err := webp.Encode(fileBytes, img, webpOptions); err != nil {
			return
		}
		*outputFilepath = replaceExtension(*outputFilepath, ".webp")
	default:
		return
	}
}

func replaceExtension(path string, newExt string) string {
	currentExt := filepath.Ext(path)
	pathWithoutExt := strings.TrimSuffix(path, currentExt)
	return pathWithoutExt + newExt
}
