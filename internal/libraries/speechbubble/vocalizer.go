package speechbubble

import (
	"mask_of_the_tomb/internal/core/sound_v2"
	"strings"
)

func Vocalize(c string) {
	switch strings.ToLower(c) {
	case "a":
		sound_v2.PlaySound("vowelA", 1.0, "main")
	case "e":
		sound_v2.PlaySound("vowelE", 1.0, "main")
	case "i":
		sound_v2.PlaySound("vowelI", 1.0, "main")
	case "o":
		sound_v2.PlaySound("vowelO", 1.0, "main")
	case "u":
		sound_v2.PlaySound("vowelU", 1.0, "main")
	}
}
