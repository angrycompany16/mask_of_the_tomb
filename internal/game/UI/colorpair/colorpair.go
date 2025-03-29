package colorpair

import "image/color"

type ColorPair struct {
	Bright      [4]uint8 `yaml:"Bright"`
	Dark        [4]uint8 `yaml:"Dark"`
	BrightColor color.Color
	DarkColor   color.Color
}

func (cp *ColorPair) LoadColorPair() {
	cp.BrightColor = color.RGBA{
		R: cp.Bright[0],
		G: cp.Bright[1],
		B: cp.Bright[2],
		A: cp.Bright[3],
	}
	cp.DarkColor = color.RGBA{
		R: cp.Dark[0],
		G: cp.Dark[1],
		B: cp.Dark[2],
		A: cp.Dark[3],
	}
}
