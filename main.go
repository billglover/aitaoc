package main

import (
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	cols     int = 12
	rows     int = 24
	size     int = 1024
	padding  int = 0
	maxAngle int = 45
)

func main() {

	rand.Seed(time.Now().UnixNano())

	dc := gg.NewContext(size, size)

	keypoints := GradientTable{
		{MustParseHex("#9e0142"), 0.0},
		{MustParseHex("#d53e4f"), 0.1},
		{MustParseHex("#f46d43"), 0.2},
		{MustParseHex("#fdae61"), 0.3},
		{MustParseHex("#fee090"), 0.4},
		{MustParseHex("#ffffbf"), 0.5},
		{MustParseHex("#e6f598"), 0.6},
		{MustParseHex("#abdda4"), 0.7},
		{MustParseHex("#66c2a5"), 0.8},
		{MustParseHex("#3288bd"), 0.9},
		{MustParseHex("#5e4fa2"), 1.0},
	}

	// find the longest edge
	max := cols
	if max < rows {
		max = rows
	}

	sqSize := size / (max + 2)

	imgOffsetX := (size - (cols+2)*sqSize) / 2

	// add a white background
	dc.SetRGBA(1, 1, 1, 1)
	dc.DrawRectangle(0, 0, float64(size), float64(size))
	dc.Fill()

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
			xOffset := offsetPct/10*rand.Float64()*float64(sqSize*2.0) - float64(sqSize)

			dc.SetRGBA(0, 0, 0, 0.5)
			dc.DrawRectangle(float64(xCoord+padding)+xOffset, float64(yCoord+padding), float64(sqSize-2*padding), float64(sqSize-2*padding))
			dc.Stroke()

			c := keypoints.GetInterpolatedColorFor(float64(y) / float64(rows))
			dc.SetRGBA(c.R, c.G, c.B, 0.2)
			dc.DrawRectangle(float64(xCoord+padding)+xOffset, float64(yCoord+padding), float64(sqSize-2*padding), float64(sqSize-2*padding))
			dc.Fill()

			dc.Pop()
		}
	}

	dc.SavePNG("out.png")
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
