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
	//go:embed fonts/JSE_ZXSpectrum.ttf
	JSE_ZXSpectrum_ttf []byte
	//go:embed fonts/dude.ttf
	C_AND_C_Red_Alert_ttf []byte
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
	//go:embed shaders/pixel_lights.kage
	Pixel_lights_kage []byte
	//go:embed shaders/transition.kage
	Transition_kage []byte
	//go:embed sprites/environment/tilemaps/export/slambox_tilemap.png
	Slambox_tilemap []byte
	//go:embed sprites/player/export/player.png
	Player_sprite []byte
	//go:embed sprites/environment/entities/export/grass.png
	Grass_tiles []byte
	//go:embed sprites/environment/entities/export/turret.png
	Turret_sprite []byte
	//go:embed sprites/UI/export/level-titlecard.png
	Level_titlecard_sprite []byte
)
