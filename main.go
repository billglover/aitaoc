package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	cols     int     = 18
	rows     int     = 20
	size     int     = 4096
	padding  int     = 0
	maxAngle int     = 45
	bgColor  bool    = false
	fgColor  bool    = true
	fgAlpha  float64 = 0.8
	theme    string  = "eXpresso"
)

// Gradients from https://uigradients.com/#Relay
var themes = map[string]GradientTable{
	"scooter":           GradientTable{{MustParseHex("#36d1dc"), 0.0}, {MustParseHex("#5b86e5"), 1.0}},
	"visionsOfGrandeur": GradientTable{{MustParseHex("#000046"), 0.0}, {MustParseHex("#1CB5E0"), 1.0}},
	"blueSkies":         GradientTable{{MustParseHex("#56CCF2"), 0.0}, {MustParseHex("#2F80ED"), 1.0}},
	"darkOcean":         GradientTable{{MustParseHex("#373B44"), 0.0}, {MustParseHex("#4286f4"), 1.0}},
	"yoda":              GradientTable{{MustParseHex("#FF0099"), 0.0}, {MustParseHex("#493240"), 1.0}},
	"amin":              GradientTable{{MustParseHex("#8E2DE2"), 0.0}, {MustParseHex("#4A00E0"), 1.0}},
	"harvey":            GradientTable{{MustParseHex("#1f4037"), 0.0}, {MustParseHex("#99f2c8"), 1.0}},
	"flare":             GradientTable{{MustParseHex("#f12711"), 0.0}, {MustParseHex("#f5af19"), 1.0}},
	"ultraViolet":       GradientTable{{MustParseHex("#654ea3"), 0.0}, {MustParseHex("#eaafc8"), 1.0}},
	"sinCityRed":        GradientTable{{MustParseHex("#ED213A"), 0.0}, {MustParseHex("#93291E"), 1.0}},
	"eveningNight":      GradientTable{{MustParseHex("#005AA7"), 0.0}, {MustParseHex("#FFFDE4"), 1.0}},
	"eXpresso":          GradientTable{{MustParseHex("#3c1053"), 0.0}, {MustParseHex("#ad5389"), 1.0}},
	"coolSky":           GradientTable{{MustParseHex("#2980B9"), 0.0}, {MustParseHex("#6DD5FA"), 0.5}, {MustParseHex("#FFFFFF"), 1.0}},
	"moonlitAsteroid":   GradientTable{{MustParseHex("#0f2027"), 0.0}, {MustParseHex("#203a43"), 0.5}, {MustParseHex("#2c5364"), 1.0}},
	"jShine":            GradientTable{{MustParseHex("#12c2e9"), 0.0}, {MustParseHex("#c471ed"), 0.5}, {MustParseHex("#f64f59"), 1.0}},
}

func main() {
	for k := range themes {
		generateImage(k, 0.2)
		generateImage(k, 0.4)
		generateImage(k, 0.6)
		generateImage(k, 0.8)
	}
}

func generateImage(theme string, fgAlpha float64) {
	name := fmt.Sprintf("%s_%d_%03.f.png", theme, size, fgAlpha*100)
	fmt.Println(name)

	rand.Seed(time.Now().UnixNano())

	dc := gg.NewContext(size, size)
	dc.SetLineWidth(1)

	// find the longest edge
	max := cols + 2
	if cols < rows {
		max = rows + 3
	}

	sqSize := size / max

	// calculate offset to ensure horizontal centring
	imgOffsetX := sqSize
	if cols < rows {
		imgOffsetX = (size - cols*sqSize) / 2
	}

	// add a background background
	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawRectangle(0, 0, float64(size), float64(size))
	dc.Fill()

	if bgColor {
		for y := 0; y < size; y++ {
			c := themes[theme].GetInterpolatedColorFor(float64(y) / float64(size))
			dc.SetRGB(c.R, c.G, c.B)
			dc.DrawLine(0, float64(y), float64(size), float64(y))
			dc.Stroke()
		}
	}

	for y := 0; y < rows; y++ {

		for x := 0; x < cols; x++ {
			dc.Push()

			xCoord := ((x + 1) * sqSize) + imgOffsetX
			yCoord := (y + 1) * sqSize

			midX := float64(xCoord + ((sqSize - 2*padding) / 2.0))
			midY := float64(yCoord + ((sqSize - 2*padding) / 2.0))

			// rotate each square by an random amount
			offsetPct := float64(y) / float64(rows)
			rotDeg := offsetPct * (rand.Float64()*float64(maxAngle*2.0) - float64(maxAngle))
			dc.RotateAbout(gg.Radians(rotDeg), midX, midY)

			// offset each square by a random amount
			xOffset := offsetPct/5*rand.Float64()*float64(sqSize*2.0) - float64(sqSize)

			c := colorful.Color{R: 1.0, G: 1.0, B: 1.0}
			if fgColor {
				c = themes[theme].GetInterpolatedColorFor(float64(y) / float64(rows))
			}

			dc.SetRGBA(c.R, c.G, c.B, 1)
			dc.DrawRectangle(float64(xCoord+padding)+xOffset, float64(yCoord+padding), float64(sqSize-2*padding), float64(sqSize-2*padding))
			dc.Stroke()

			dc.SetRGBA(c.R, c.G, c.B, fgAlpha)
			dc.DrawRectangle(float64(xCoord+padding)+xOffset, float64(yCoord+padding), float64(sqSize-2*padding), float64(sqSize-2*padding))
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
	sPosX := float64(cols*sqSize) + float64(imgOffsetX) - w
	sPosY := (float64(rows)+2.5)*float64(sqSize) - h

	dc.SetRGBA(0.5, 0.5, 0.5, 0.5)
	dc.DrawString(caption, sPosX, sPosY)

	err = dc.SavePNG("img/" + name)
	if err != nil {
		log.Fatal("unable to load font:", err)
	}
}

// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (self GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(self)-1; i++ {
		c1 := self[i]
		c2 := self[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return self[len(self)-1].Col
}

// This is a very nice thing Golang forces you to do!
// It is necessary so that we can write out the literal of the colortable below.
func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}
