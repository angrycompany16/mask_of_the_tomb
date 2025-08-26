package scenes

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/core/sound"
	"mask_of_the_tomb/internal/libraries/musicplayer"
	"mask_of_the_tomb/internal/libraries/node"
	save "mask_of_the_tomb/internal/libraries/savesystem"
)

type BaseScene struct {
	musicPlayer *musicplayer.MusicPlayer
	lock        bool
}

func (m *BaseScene) Init() {
	gameData := errs.Must(save.GetSaveAsset("saveData"))
	resources.Settings = gameData.Settings

	selectSoundStream := errs.Must(assettypes.GetOggStream("selectSound"))
	node.SelectSound = &sound.EffectPlayer{errs.Must(sound.FromStream(selectSoundStream)), 1.0}

	dialogueSoundStream := errs.Must(assettypes.GetOggStream("dialogueSound"))
	node.DialogueSound = &sound.EffectPlayer{errs.Must(sound.FromStream(dialogueSoundStream)), 1.0}
}

func (m *BaseScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	// Not sure if this should really be here
	// m.musicPlayer.ResetMusic()
	// if resources.DebugMode {
	// 	ebitenutil.DebugPrint(rendering.ScreenLayers.Overlay, fmt.Sprintf("TPS: %0.2f \nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
	// }
	events.Update()

	if !m.lock {
		m.lock = true
		return &scene.SceneTransition{
			Kind:       scene.Push,
			OtherScene: MakeMenuScene(),
		}, true
	}
	return nil, false
}

func (m *BaseScene) Draw()           {}
func (m *BaseScene) GetName() string { return "musicScene" }
func MakeBaseScene() *BaseScene      { return &BaseScene{musicplayer.NewMusicPlayer(), false} }
