package assets

import (
	_ "embed"
	"path/filepath"
)

var (
	EnvironmentTilemapFolder = filepath.Join("assets", "sprites", "environment", "tilemaps", "export")
	PlayerFolder             = filepath.Join("assets", "sprites", "player", "export")

	//go:embed fonts/JSE_AmigaAMOS.ttf
	JSE_AmigaAMOS_ttf []byte

	//go:embed music/placeholder/jungle-ish-beat-for-video-games-314073.mp3
	Menu_mp3 []byte

	//go:embed music/homemade/actually_basement.wav
	Basement_wav []byte

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

	//go:embed shaders/pixel_lights.kage
	Pixel_lights_kage []byte

	//go:embed sprites/environment/tilemaps/export/slambox_tilemap.png
	Slambox_tilemap []byte

	//go:embed sprites/player/export/player.png
	Player_sprite []byte
)
