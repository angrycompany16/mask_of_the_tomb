package scenes

import (
	"fmt"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/scene"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DebugScene struct {
	actualFPS       float64
	actualTPS       float64
	renderDebugInfo ebiten.DebugInfo
}

func (d *DebugScene) Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool) {
	if inpututil.IsKeyJustReleased(ebiten.KeyTab) {
		return &scene.SceneTransition{
			Kind: scene.PopName,
			Name: d.GetName(),
		}, true
	}

	// Fetch interesting debug data
	d.actualFPS = ebiten.ActualFPS()
	d.actualTPS = ebiten.ActualTPS()
	ebiten.ReadDebugInfo(&d.renderDebugInfo)
	return &scene.SceneTransition{}, false
}

func (d *DebugScene) Draw() {
	debugString := fmt.Sprintf(`
		FPS: %f,
		TPS: %f,
		Rendering info: %+v
	`,
		d.actualFPS,
		d.actualTPS,
		d.renderDebugInfo)
	ebitenutil.DebugPrint(rendering.ScreenLayers.ScreenUI, debugString)
}
func (d *DebugScene) Init()           {}
func (d *DebugScene) GetName() string { return "debugScene" }

func MakeDebugScene() *DebugScene {
	return &DebugScene{}
}
