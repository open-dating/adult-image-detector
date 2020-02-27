package main

import (
	"errors"
	"gocv.io/x/gocv"
	"image"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"os"
	"path/filepath"
)

// get score from yahoo open nsfw https://github.com/yahoo/open_nsfw
func GetOpenNsfwScore(filePath string) (score float32, err error) {
	img := gocv.IMRead(filePath, gocv.IMReadColor)
	if img.Empty() {
		RemoveFile(filePath)
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
		RemoveFile(filePath)
		return 0, errors.New("Invalid blob")
	}
	defer blob.Close()

	protoPath, _ := filepath.Abs("./models/open_nsfw/nsfw_model/deploy.prototxt")
	modelPath, _ := filepath.Abs("./models/open_nsfw/nsfw_model/resnet_50_1by2_nsfw.caffemodel")

	net := gocv.ReadNetFromCaffe(
		protoPath,
		modelPath,
	)
	if net.Empty() {
		RemoveFile(filePath)
		return 0, errors.New("Invalid net")
	}
	defer net.Close()

	net.SetInput(blob, "")

	detBlob := net.Forward("")
	defer detBlob.Close()

	return detBlob.GetFloatAt(0, 1), nil
}

// get result from An Algorithm for Nudity Detection by Rigan Ap-apid
func getAnAlgorithmForNudityDetectionResult(filePath string, debug bool) (result bool, err error)  {
	existingImageFile, err := os.Open(filePath)
	if err != nil {
		return true, errors.New("Cant open file")
	}
	defer existingImageFile.Close()

	imageData, _, err := image.Decode(existingImageFile)
	if err != nil {
		RemoveFile(filePath)
		return true, errors.New("Decode err")
	}

	anAlg := AnAlgorithm{
		img:imageData,
		debug:debug,
	}
	return anAlg.IsNude()
}
