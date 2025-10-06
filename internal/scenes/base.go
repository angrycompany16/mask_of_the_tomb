package scenes

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/core/sound"
	"mask_of_the_tomb/internal/libraries/musicplayer"
	"mask_of_the_tomb/internal/libraries/node"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type BaseScene struct {
	musicPlayer *musicplayer.MusicPlayer
	lock        bool
	initTime    time.Time
}

func (b *BaseScene) Init() {
	gameData := errs.Must(save.GetSaveAsset("saveData"))
	resources.Settings = gameData.Settings

	selectSoundStream := errs.Must(assettypes.GetOggStream("selectSound"))
	node.SelectSound = &sound.EffectPlayer{errs.Must(sound.FromStream(selectSoundStream)), 1.0}

	dialogueSoundStream := errs.Must(assettypes.GetOggStream("dialogueSound"))
	node.DialogueSound = &sound.EffectPlayer{errs.Must(sound.FromStream(dialogueSoundStream)), 1.0}
}

func (b *BaseScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	b.musicPlayer.ResetMusicVolume()
	events.Update()
	resources.Time = time.Since(b.initTime).Seconds()

	if !b.lock {
		b.lock = true
		return &scene.SceneTransition{
			Kind:       scene.Push,
			OtherScene: MakeMenuScene(),
		}, true
	}
	return nil, false
}

func (b *BaseScene) Draw() {
	if resources.DebugMode {
		fmt.Println(ebiten.ActualFPS())
		fmt.Println(ebiten.ActualTPS())
	}
}
func (b *BaseScene) GetName() string { return "musicScene" }
func MakeBaseScene() *BaseScene      { return &BaseScene{musicplayer.NewMusicPlayer(), false, time.Now()} }
