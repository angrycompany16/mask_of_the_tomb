package assets

import (
	_ "embed"
)

var (
	//go:embed fonts/JSE_AmigaAMOS.ttf
	JSE_AmigaAMOS_ttf []byte

	//go:embed music/placeholder/jungle-ish-beat-for-video-games-314073.mp3
	Menu_mp3 []byte

	//go:embed music/placeholder/Sneaky-Snitch(chosic.com).mp3
	Level_mp3 []byte
)
