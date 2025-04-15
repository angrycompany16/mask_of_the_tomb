package colorpair

import (
	"image/color"
	"mask_of_the_tomb/internal/errs"
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
