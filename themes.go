package main

// Code for generating colour themes is adapted from the example provided in the documentation of the go-colorful package

import (
	"log"

	colorful "github.com/lucasb-eyer/go-colorful"
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
	"mono":              GradientTable{{MustParseHex("#000000"), 0.0}, {MustParseHex("#ffffff"), 1.0}},
	"black":             GradientTable{{MustParseHex("#000000"), 0.0}, {MustParseHex("#000000"), 1.0}},
	"white":             GradientTable{{MustParseHex("#ffffff"), 0.0}, {MustParseHex("#ffffff"), 1.0}},
	"coolSky":           GradientTable{{MustParseHex("#2980B9"), 0.0}, {MustParseHex("#6DD5FA"), 0.5}, {MustParseHex("#FFFFFF"), 1.0}},
	"moonlitAsteroid":   GradientTable{{MustParseHex("#0f2027"), 0.0}, {MustParseHex("#203a43"), 0.5}, {MustParseHex("#2c5364"), 1.0}},
	"jShine":            GradientTable{{MustParseHex("#12c2e9"), 0.0}, {MustParseHex("#c471ed"), 0.5}, {MustParseHex("#f64f59"), 1.0}},
}

// MustParseHex takes a HEX color value and parses it to ensure
// it is valid. If invalid it terminates indicating the cause
// of the error.
func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		log.Fatal("MustParseHex:", err)
	}
	return c
}

// GradientTable holds "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

// GetInterpolatedColor returns a HCL-blend between the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (gt GradientTable) GetInterpolatedColor(t float64) colorful.Color {
	for i := 0; i < len(gt)-1; i++ {
		c1 := gt[i]
		c2 := gt[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	return gt[len(gt)-1].Col
}
