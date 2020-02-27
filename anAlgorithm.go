// An Algorithm for Nudity Detection
// by Rigan Ap-apid
// http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.96.9872&rep=rep1&type=pdf
package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"image/draw"
	"sort"
)

const (
	SkinCbMin = 80
	SkinCbMax = 120
	SkinCrMin = 133
	SkinCrMax = 173
)

type AnAlgorithm struct {
	img                  image.Image
	height               int
	width                int
	skinMap              image.Image
	backgroundPixelCount int
	skinPixelCount       int
	regions              [] AnAlgorithmRegion
	regionsContours      [][] image.Point
	boundsPoly           AnAlgorithmBoundsPolygon
	debug 				 bool
}

type AnAlgorithmRegion struct {
	area    float64
	contour [] image.Point
}

type AnAlgorithmBoundsPolygon struct {
	area              float64
	contour           [] image.Point
	hue               float64
	skinPixels        [] color.Color
	skinPixelCount    int
	avgSkinsIntensity float64
	image             image.Image
	height            int
	width             int
}

// find skins on image and create mask with white/black regions
// do manually, becouse have some trouble with gocv.InRange https://github.com/hybridgroup/gocv/issues/159
func (a *AnAlgorithm) maskSkinAndCountSkinPixels() () {
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

	// out, _ := os.Create("uploads/0.maskSkinAndCountSkinPixels.jpg")
	// jpeg.Encode(out, a.skinMap, nil)
}

func (a *AnAlgorithm) findRegions() {
	img, _ := ImageToRGB8Mat(a.skinMap)

	mask := gocv.NewMat()
	gocv.InRangeWithScalar(img, gocv.NewScalar(0.0, 0.0, 0.0, 0.0), gocv.NewScalar(1.0, 1.0, 1.0, 0.0), &mask)

	contours := gocv.FindContours(mask, gocv.RetrievalList, gocv.ChainApproxNone)
	a.regionsContours = contours

	for i := range contours {
		a.regions = append(a.regions, AnAlgorithmRegion{
			area: gocv.ContourArea(contours[i]),
			contour: contours[i],
		})
	}

	sort.Slice(a.regions, func(i, j int) bool {
		return a.regions[i].area > a.regions[j].area
	})

	// if we havent 2rd region
	if len(a.regions) == 1 {
		a.regions = append(a.regions, a.regions[0])
	}

	// if we havent 3rd region
	if len(a.regions) == 2 {
		a.regions = append(a.regions, a.regions[1])
	}
}


// Identify the leftmost, the uppermost, the rightmost,
// and the lowermost skin pixels of the three largest
// skin regions. Use these points as the corner points
func (a *AnAlgorithm) findBoundsPolyCorners() {
	leftmost := a.width
	uppermost := a.height
	rightmost := 0
	lowermost := 0

	for i, region := range a.regions {
		if i > 2 {
			break
		}

		// find corners
		for _, p := range region.contour {
			if p.X < leftmost {
				leftmost = p.X
			}
			if p.X > rightmost {
				rightmost = p.X
			}
			if p.Y < uppermost {
				uppermost = p.Y
			}
			if p.Y > lowermost {
				lowermost = p.Y
			}
		}
	}

	width := rightmost - leftmost
	height := lowermost - uppermost

	a.boundsPoly = AnAlgorithmBoundsPolygon{
		area: float64(width * height),
		contour: []image.Point{
			image.Point{X: leftmost, Y:uppermost},
			image.Point{X: rightmost, Y:lowermost},
		},
		skinPixelCount: 0,
	}

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	a.boundsPoly.image = image.NewRGBA(image.Rectangle{upLeft, lowRight})
}

// create poly from images, cacl pixel count
// and save pixels fro find avg in we need it
func (a *AnAlgorithm) createBoundsPolyAndCalcSkins() {

	xBig := 0
	yBig := 0

	for x := a.boundsPoly.contour[0].X; x < a.boundsPoly.contour[1].X; x++ {
		for y := a.boundsPoly.contour[0].Y; y < a.boundsPoly.contour[1].Y; y++ {
			pixCol := a.img.At(x, y)

			if a.yCbCrSkinDetector(pixCol) == true {
				a.boundsPoly.skinPixelCount++
				a.boundsPoly.skinPixels = append(a.boundsPoly.skinPixels, pixCol)
			}

			a.boundsPoly.image.(draw.Image).Set(xBig, yBig, pixCol)

			yBig++
		}
		yBig = 0
		xBig++
	}

	// out, _ := os.Create("uploads/5.createBoundsPolyAndCalcSkins.jpg")
	// jpeg.Encode(out, a.boundsPoly.image, nil)
}

// detect is skin
func (a *AnAlgorithm) yCbCrSkinDetector(pixCol color.Color) bool {
	_, cb, cr := RgbaToYCbCr(pixCol)

	return cb >= SkinCbMin && cb <= SkinCbMax && cr >= SkinCrMin && cr <= SkinCrMax
}


// find avg intensity in boundsPoly skins region
func (a *AnAlgorithm) findAverageSkinsIntensityInBoundsPoly() {

	a.boundsPoly.avgSkinsIntensity = 0

	skinsLen := len(a.boundsPoly.skinPixels)
	if skinsLen == 0 {
		return
	}

	var cbSum int = 0
	var crSum int = 0
	for i := 0 ; i < skinsLen; i++ {
		_, cb, cr := RgbaToYCbCr(a.boundsPoly.skinPixels[i])

		cbSum += cb
		crSum += cr
	}

	avgColorVal := float64((cbSum + crSum) / skinsLen)
	if avgColorVal == 0 {
		return
	}

	a.boundsPoly.avgSkinsIntensity = float64(SkinCbMax - SkinCbMin + SkinCrMax - SkinCrMin) / avgColorVal
}

// return is nude image or not
func (a *AnAlgorithm) IsNude() (bool, error) {
	a.width = a.img.Bounds().Max.X
	a.height = a.img.Bounds().Max.Y

	a.maskSkinAndCountSkinPixels()

	totalPixelCount := a.skinPixelCount + a.backgroundPixelCount

	totalSkinPortion := float32(a.skinPixelCount) / float32(totalPixelCount)

	if totalPixelCount == 0 {
		if a.debug {
			fmt.Println("No pixels found")
		}
		return false, nil
	}

	// Criteria (a)
	if a.debug {
		fmt.Println("a: totalSkinPortion=", totalSkinPortion, " < 0.15")
	}
	if totalSkinPortion < 0.15 {
		return false, nil
	}

	a.findRegions()
	largestRegionPortion := 0.0
	nextRegionPortion := 0.0
	thirdRegionPortion := 0.0

	if len(a.regions) > 0 {
		largestRegionPortion = a.regions[0].area / float64(a.skinPixelCount)
		nextRegionPortion = a.regions[1].area / float64(a.skinPixelCount)
		thirdRegionPortion = a.regions[2].area / float64(a.skinPixelCount)
	}

	// Criteria (b)
	if a.debug {
		fmt.Println("b: largestRegionPortion=", largestRegionPortion, " < 0.35 && nextRegionPortion=", nextRegionPortion, " < 0.30 && thirdRegionPortion=", thirdRegionPortion, " < 0.30")
	}
	if largestRegionPortion < 0.35 && nextRegionPortion < 0.30 && thirdRegionPortion < 0.30 {
		return false, nil
	}

	// Criteria (c)
	if a.debug {
		fmt.Println("c: largestRegionPortion=", largestRegionPortion, " < 0.45")
	}
	if largestRegionPortion < 0.45 {
		return false, nil
	}

	// Criteria (d)
	a.findBoundsPolyCorners()
	a.createBoundsPolyAndCalcSkins()

	if a.debug {
		fmt.Println("d: totalSkinPortion=", totalSkinPortion, " < 0.30")
	}
	if totalSkinPortion < 0.30 {
		boundsPolySkinPortion := float64(a.boundsPoly.skinPixelCount) / a.boundsPoly.area

		if a.debug {
			fmt.Println("d: boundsPolySkinPortion=", boundsPolySkinPortion, " < 0.55")
		}
		if boundsPolySkinPortion < 0.55 {
			return false, nil
		}
	}

	// Criteria (e)
	if a.debug {
		fmt.Println("e: len(a.regions)=", len(a.regions), " > 60")
	}
	if len(a.regions) > 60 {
		a.findAverageSkinsIntensityInBoundsPoly()

		if a.debug {
			fmt.Println("e: boundsPoly.avgSkinsIntensity=", a.boundsPoly.avgSkinsIntensity, " < 0.25")
		}
		if a.boundsPoly.avgSkinsIntensity < 0.25 {
			return false, nil
		}
	}

	return true, nil
}


