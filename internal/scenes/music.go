package scenes

import (
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/libraries/musicplayer"
)

type MusicScene struct {
	musicPlayer *musicplayer.MusicPlayer
}

func (m *MusicScene) Init() {
	m.musicPlayer.Init()
}

func (m *MusicScene) Update() {
	// Not sure if this should really be here
	// m.musicPlayer.ResetMusic()
	// if resources.DebugMode {
	// 	ebitenutil.DebugPrint(rendering.ScreenLayers.Overlay, fmt.Sprintf("TPS: %0.2f \nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	// }
}

func (m *MusicScene) Draw() {}
func (m *MusicScene) Exit() (scene.SceneTransition, bool) {
	return scene.SceneTransition{
		Kind: scene.Sibling,
		Name: "menuScene",
	}, true
}
