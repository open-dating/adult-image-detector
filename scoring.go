package main

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gocv.io/x/gocv"
)

// get score from yahoo open nsfw https://github.com/yahoo/open_nsfw
func GetOpenNsfwScore(filePath string, net gocv.Net) (score float32, err error) {
	img := gocv.IMRead(filePath, gocv.IMReadColor)
	if img.Empty() {
		return 0, errors.New("Invalid image")
	}
	defer img.Close()

	blob := gocv.BlobFromImage(
		img,
		1.0,
		image.Pt(224, 224),
		gocv.NewScalar(104, 117, 123, 0),
		true,
		false,
	)
	if blob.Empty() {
		return 0, errors.New("Invalid blob")
	}
	defer blob.Close()

	net.SetInput(blob, "")

	detBlob := net.Forward("")
	defer detBlob.Close()

	return detBlob.GetFloatAt(0, 1), nil
}

// get result from An Algorithm for Nudity Detection by Rigan Ap-apid
func getAnAlgorithmForNudityDetectionResult(filePath string, debug bool) (result bool, err error) {
	existingImageFile, err := os.Open(filePath)
	if err != nil {
		return true, errors.New("Cant open file")
	}
	defer existingImageFile.Close()

	imageData, _, err := image.Decode(existingImageFile)
	if err != nil {
		return true, errors.New("Decode err")
	}

	anAlg := AnAlgorithm{
		img:   imageData,
		debug: debug,
	}
	return anAlg.IsNude()
}
