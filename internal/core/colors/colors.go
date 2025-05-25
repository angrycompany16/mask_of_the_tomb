package colors

import (
	"image/color"
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

			parseColor := func(i int) (uint8, error) {
				color, err := strconv.ParseInt(colorValue[i].Value, 10, 16)
				return uint8(color), err
			}

			switch content.Value {
			case "Bright":
				r, err := parseColor(0)
				g, err := parseColor(1)
				b, err := parseColor(2)
				a, err := parseColor(3)
				if err != nil {
					return err
				}
				cp.BrightColor = color.RGBA{r, g, b, a}
			case "Dark":
				r, err := parseColor(0)
				g, err := parseColor(1)
				b, err := parseColor(2)
				a, err := parseColor(3)
				if err != nil {
					return err
				}
				cp.DarkColor = color.RGBA{r, g, b, a}
			}
		}
	}
	return nil
}

type YAMLColor struct {
	color.Color
}

func (c *YAMLColor) UnmarshalYAML(node *yaml.Node) error {
	parseColor := func(i int) (uint8, error) {
		color, err := strconv.ParseInt(node.Content[i].Value, 10, 16)
		return uint8(color), err
	}

	r, err := parseColor(0)
	g, err := parseColor(0)
	b, err := parseColor(0)
	a, err := parseColor(0)
	if err != nil {
		return err
	}
	c.Color = color.RGBA{r, g, b, a}
	return nil
}

// Takes an input string of the form #RRGGBB and converts it into an RGBA value
func HexToRGB(hexInput string) (color.RGBA, error) {
	// TODO: Regex to check that the string is correct

	R, err := strconv.ParseInt(hexInput[1:3], 16, 8)
	G, err := strconv.ParseInt(hexInput[3:5], 16, 8)
	B, err := strconv.ParseInt(hexInput[5:7], 16, 8)

	return color.RGBA{uint8(R), uint8(G), uint8(B), 1}, err
}
