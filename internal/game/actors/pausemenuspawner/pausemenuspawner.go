package pausemenuspawner

import (
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/bundles"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseMenuSpawner struct {
	*nodeactor.Node
	active bool
}

func (p *PauseMenuSpawner) Init(cmd *commands.Commands) {
	p.Node.Init(cmd)
	UIControls := cmd.InputHandler.InputSchemes["UIControls"]
	UIControls.AddBinding("Pause", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
	})
}

func (p *PauseMenuSpawner) Update(cmd *commands.Commands) {
	p.Node.Update(cmd)
	// TODO: Disable gameplay when pause menu is active
	// Q: What to disable?
	// Idea:
	// - Slow down game time
	// - Disable input
	UIControls := cmd.InputHandler.InputSchemes["UIControls"]
	playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
	if UIControls.PollAction("Pause") {
		scene, _ := commands.Get[engine.Scene](cmd)
		pauseMenu, exists := p.getPauseMenu(scene)
		if exists {
			playerControls.Active = true
			p.active = false
			scene.Delete(pauseMenu)
		} else {
			playerControls.Active = false
			scene.SpawnBundleV2(cmd, bundles.MakePauseMenuBundle())
			p.active = true
		}
	}
}

func (p *PauseMenuSpawner) getPauseMenu(scene *engine.Scene) (*engine.Node, bool) {
	pauseScreenRoot, ok := scene.GetNodeByName("PauseScreenBundle")
	return pauseScreenRoot, ok
}

func NewPauseMenuSpawner(node *nodeactor.Node) *PauseMenuSpawner {
	return &PauseMenuSpawner{
		Node: node,
	}
}
