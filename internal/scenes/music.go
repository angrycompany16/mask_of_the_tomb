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

func (m *MusicScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	// Not sure if this should really be here
	// m.musicPlayer.ResetMusic()
	// if resources.DebugMode {
	// 	ebitenutil.DebugPrint(rendering.ScreenLayers.Overlay, fmt.Sprintf("TPS: %0.2f \nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	// }
	return &scene.SceneTransition{
		Kind:       scene.Push,
		OtherScene: &MenuScene{},
	}, true
}

func (m *MusicScene) Draw()           {}
func (m *MusicScene) GetName() string { return "musicScene" }
