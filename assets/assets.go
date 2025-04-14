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
	Basement_mp3 []byte

	//go:embed music/placeholder/game-music-puzzle-strategy-arcade-technology-301226.mp3
	Library_mp3 []byte

	//go:embed sfx/dash.wav
	Dash_wav []byte
)
