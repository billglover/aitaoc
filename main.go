package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

type config struct {
	numCols     int
	numRows     int
	width       int
	height      int
	padding     int
	strokeWidth float64
	maxRotation float64
	maxOffset   float64
	alpha       float64
	theme       string
}

func main() {
	var numCols, numRows, width, height, padding int
	var strokeWidth, maxRotation, maxOffset, alpha float64
	var theme string
	flag.IntVar(&numCols, "numCols", 12, "number of columns in the image")
	flag.IntVar(&numRows, "numRows", 24, "number of rows in the image")
	flag.IntVar(&width, "width", 4096, "width of the image in pixels")
	flag.IntVar(&height, "height", 4096, "height of the image in pixels")
	flag.IntVar(&padding, "padding", 0, "padding between blocks in pixels")
	flag.Float64Var(&strokeWidth, "strokeWidth", 1, "stroke width in pixels")
	flag.Float64Var(&maxRotation, "maxRotation", 90, "maximum angle of rotation in degrees")
	flag.Float64Var(&maxOffset, "maxOffset", 1, "maximum offset as a fraction of block size")
	flag.Float64Var(&alpha, "alpha", 0.2, "alpha value for boxes")
	flag.StringVar(&theme, "theme", "mono", "colour theme to use")
	flag.Parse()

	cfg := config{
		numCols:     numCols,
		numRows:     numRows,
		width:       width,
		height:      height,
		strokeWidth: strokeWidth,
		padding:     padding,
		maxRotation: maxRotation,
		maxOffset:   maxOffset,
		alpha:       alpha,
		theme:       theme,
	}

	generateImage(cfg, "eveningNight", 0.2)
}

func generateImage(cfg config, theme string, fgAlpha float64) {

	rand.Seed(time.Now().UnixNano())

	name := fmt.Sprintf("%s_%dx%d_%03.f.png", theme, cfg.width, cfg.height, fgAlpha*100)

	dc := gg.NewContext(cfg.width, cfg.height)

	// figure out if the image is wide or tall
	hSize := cfg.width / (cfg.numCols + 2)
	vSize := cfg.height / (cfg.numRows + 3)

	sqSize := hSize
	if hSize > vSize {
		sqSize = vSize
	}

	// calculate offset to ensure horizontal centring
	imgOffsetX := (cfg.width - (cfg.numCols+2)*sqSize) / 2

	// add a white background
	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawRectangle(0, 0, float64(cfg.width), float64(cfg.height))
	dc.Fill()

	for y := 0; y < cfg.numRows; y++ {

		for x := 0; x < cfg.numCols; x++ {
			dc.Push()

			xCoord := ((x + 1) * sqSize) + imgOffsetX
			yCoord := (y + 1) * sqSize

			midX := float64(xCoord + ((sqSize - 2*cfg.padding) / 2.0))
			midY := float64(yCoord + ((sqSize - 2*cfg.padding) / 2.0))

			// rotate each square by an random amount up to a pre-defined max
			offsetPct := float64(y) / float64(cfg.numRows)
			rotDeg := offsetPct * (rand.Float64()*cfg.maxRotation*2.0 - cfg.maxRotation)
			dc.RotateAbout(gg.Radians(rotDeg), midX, midY)

			// offset each square by a random amount
			xOffset := cfg.maxOffset * offsetPct * (rand.Float64()*float64(sqSize*2) - float64(sqSize))

			// calculate colour based on theme and row
			c := themes[cfg.theme].GetInterpolatedColor(float64(y) / float64(cfg.numRows))

			// draw the box outline
			dc.SetRGBA(c.R, c.G, c.B, 1)
			dc.SetLineWidth(cfg.strokeWidth)
			dc.DrawRectangle(float64(xCoord+cfg.padding)+xOffset, float64(yCoord+cfg.padding), float64(sqSize-2*cfg.padding), float64(sqSize-2*cfg.padding))
			dc.Stroke()

			// fill in the box
			dc.SetRGBA(c.R, c.G, c.B, cfg.alpha)
			dc.DrawRectangle(float64(xCoord+cfg.padding)+xOffset, float64(yCoord+cfg.padding), float64(sqSize-2*cfg.padding), float64(sqSize-2*cfg.padding))
			dc.Fill()

			dc.Pop()
		}
	}

	err := dc.LoadFontFace("Go-Mono.ttf", 18.0)
	if err != nil {
		log.Fatal("unable to load font:", err)
	}

	hSum := sha256.New()
	dc.EncodePNG(hSum)
	caption := fmt.Sprintf("%x // @BillGlover", hSum.Sum(nil))

	w, h := dc.MeasureString(caption)
	sPosX := float64(cfg.numCols*sqSize) + float64(imgOffsetX) - w
	sPosY := (float64(cfg.numRows)+2.5)*float64(sqSize) - h

	dc.SetRGBA(0.5, 0.5, 0.5, 0.5)
	dc.DrawString(caption, sPosX, sPosY)

	err = dc.SavePNG("img/" + name)
	if err != nil {
		log.Fatal("unable to load font:", err)
	}
}
