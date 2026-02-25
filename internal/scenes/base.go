package scenes

import (
	"fmt"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/events"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	"mask_of_the_tomb/internal/core/sound_v2"
	save "mask_of_the_tomb/internal/libraries/savesystem"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resound/effects"
)

type BaseScene struct {
	lock     bool
	initTime time.Time
}

// bruh. Scene system is kinda ass
func (b *BaseScene) Init() {
	gameData := errs.Must(save.GetSaveAsset("saveData"))
	resources.Settings = gameData.Settings

	sound_v2.AddDSPChannelEffect("sfxMaster", "vol", effects.NewVolume().SetStrength(resources.Settings.GetTotalSfxVolume()))
	sound_v2.AddDSPChannelEffect("musicMaster", "vol", effects.NewVolume().SetStrength(resources.Settings.GetTotalMusicVolume()))
}

func (b *BaseScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	events.Update()
	resources.Time = time.Since(b.initTime).Seconds()
	if inpututil.IsKeyJustReleased(ebiten.KeyTab) {
		_, _, exists := sceneStack.GetScene(MakeDebugScene().GetName())
		if !exists {
			return &scene.SceneTransition{
				Kind:       scene.Push,
				OtherScene: MakeDebugScene(),
			}, true
		}
	}

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
func MakeBaseScene() *BaseScene      { return &BaseScene{false, time.Now()} }
