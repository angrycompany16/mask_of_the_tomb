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

	//go:embed sfx/slam.wav
	Slam_wav []byte

	//go:embed sfx/death.mp3
	Death_mp3 []byte

	//go:embed sfx/select.ogg
	Select_ogg []byte

	//go:embed sfx/text-scroll.ogg
	Text_scroll_ogg []byte

	//go:embed shaders/fog.kage
	Fog_kage []byte

	//go:embed shaders/vignette.kage
	Vignette_kage []byte

	//go:embed shaders/lights_additive.kage
	Lights_additive_kage []byte

	//go:embed shaders/lights_subtractive.kage
	Lights_subtractive_kage []byte
)
