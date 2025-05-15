package colors

import (
	"image/color"
	"mask_of_the_tomb/internal/core/errs"
	"strconv"

	"gopkg.in/yaml.v3"
)

type ColorPair struct {
	BrightColor color.Color
	DarkColor   color.Color
}

func (cp *ColorPair) UnmarshalYAML(node *yaml.Node) error {
	for i, content := range node.Content {
		if content.Value == "Bright" || content.Value == "Dark" {
			colorValue := node.Content[i+1].Content

			parseColor := func(i int) uint8 {
				return uint8(errs.Must(strconv.ParseInt(colorValue[i].Value, 10, 16)))
			}

			switch content.Value {
			case "Bright":
				cp.BrightColor = color.RGBA{parseColor(0), parseColor(1), parseColor(2), parseColor(3)}
			case "Dark":
				cp.DarkColor = color.RGBA{parseColor(0), parseColor(1), parseColor(2), parseColor(3)}
			}
		}
	}
	return nil
}

type YAMLColor struct {
	color.Color
}

func (c *YAMLColor) UnmarshalYAML(node *yaml.Node) error {
	parseColor := func(i int) uint8 {
		return uint8(errs.Must(strconv.ParseInt(node.Content[i].Value, 10, 16)))
	}

	c.Color = color.RGBA{parseColor(0), parseColor(1), parseColor(2), parseColor(3)}
	return nil
}

// Takes an input string of the form #RRGGBB and converts it into an RGBA value
func HexToRGB(hexInput string) color.RGBA {
	// TODO: Regex to check that the string is correct
	R := uint8(errs.Must(strconv.ParseInt(hexInput[1:3], 16, 8)))
	G := uint8(errs.Must(strconv.ParseInt(hexInput[3:5], 16, 8)))
	B := uint8(errs.Must(strconv.ParseInt(hexInput[5:7], 16, 8)))

	return color.RGBA{R, G, B, 1}
}
