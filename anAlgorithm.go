// An Algorithm for Nudity Detection
// by Rigan Ap-apid
// http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.96.9872&rep=rep1&type=pdf
package main

import (
	"image"
	"image/color"
	"image/draw"
)

type AnAlgorithm struct {
	img image.Image
	height int
	width int
	skinMap image.Image
	backgroundPixelCount int
	skinPixelCount int
}

// find skins on image and mutate image to white/black regions
func (a *AnAlgorithm) mapSkinPixels() () {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{a.width, a.height}

	a.skinMap = image.NewRGBA(image.Rectangle{upLeft, lowRight})

	black := color.RGBA{0,0,0, 0xff}
	white := color.RGBA{255,255,255, 0xff}

	a.backgroundPixelCount = 0
	a.skinPixelCount = 0

	var drawPixCol color.RGBA
	for x := 0; x < a.width; x++ {
		for y := 0; y < a.height; y++ {
			pixCol := a.img.At(x, y)

			drawPixCol = white
			if a.yCbCrSkinDetector(pixCol) == true {
				a.skinPixelCount++
				drawPixCol = black
			} else {
				a.backgroundPixelCount++
			}

			a.skinMap.(draw.Image).Set(x, y, drawPixCol)
		}
	}

}

// detect is skin
func (a *AnAlgorithm) yCbCrSkinDetector(pixCol color.Color) bool {
	_, cb, cr := rgbaToYCbCr(pixCol)

	return cb >= 80 && cb <= 120 && cr >= 133 && cr <= 173
}

// return is nude image or not
func (a *AnAlgorithm) IsNude() (bool, error) {
	a.width = a.img.Bounds().Max.X
	a.height = a.img.Bounds().Max.Y

	a.mapSkinPixels()

	totalPixelCount := a.skinPixelCount + a.backgroundPixelCount

	if totalPixelCount == 0 {
		return false, nil
	}

	totalSkinPortion := float32(a.skinPixelCount) / float32(totalPixelCount)

	// Criteria (a)
	if (totalSkinPortion < 0.15) {
		return false, nil
	}

	// TODO Criteria (b)

	// TODO Criteria (c)

	// TODO Criteria (d)

	return true, nil
}

// convert rgb to ycbr
func rgbaToYCbCr(pixCol color.Color) (y int, cb int, cr int) {
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

