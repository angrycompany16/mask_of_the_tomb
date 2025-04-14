package colors

import (
	"image/color"
	"mask_of_the_tomb/internal/errs"
	"strconv"
)

// Takes an input string of the form #RRGGBB and converts it into an RGBA value
func HexToRGB(hexInput string) color.RGBA {
	// TODO: Regex to check that the string is correct
	R := uint8(errs.Must(strconv.ParseInt(hexInput[1:3], 16, 8)))
	G := uint8(errs.Must(strconv.ParseInt(hexInput[3:5], 16, 8)))
	B := uint8(errs.Must(strconv.ParseInt(hexInput[5:7], 16, 8)))

	return color.RGBA{R, G, B, 1}
}
