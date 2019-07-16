package main

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

// convert rgb to ycbr
func RgbaToYCbCr(pixCol color.Color) (y int, cb int, cr int) {
	rInt, gInt, bInt, _ := pixCol.RGBA()

	// y8, cb8, cr8 := color.RGBToYCbCr(uint8(rInt), uint8(gInt), uint8(bInt))
	// return int(y8), int(cb8), int(cr8)

	r := float32(rInt / 256)
	g := float32(gInt / 256)
	b := float32(bInt / 256)

	return int(16.0 + 0.256788*r + 0.504129*g +  0.097905*b),
		int(128.0 - 0.148223*r - 0.290992*g +  0.439215*b),
		int(128.0 + 0.439215*r - 0.367788*g -  0.071427*b)
}

// Convert image to gocv.Mat
// https://github.com/hybridgroup/gocv/issues/228
func ImageToRGB8Mat(img image.Image) (gocv.Mat, error) {
	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()
	bytes := make([]byte, 0, x*y*3)

	//don't get surprised of reversed order everywhere below
	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
		}
	}
	return gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
}
