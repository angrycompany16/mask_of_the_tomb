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

// TODO: Bug: If we pause while the player is currently in a level
// transition, then the completion of said transition will re-enable
// the player input scheme, meaning that we can control the player
// when in reality we should not
// A solution would be to somehow stop the update loop for the game.
// although "live" pause screens are kind of cool i think
func (p *PauseMenuSpawner) Update(cmd *commands.Commands) {
	p.Node.Update(cmd)
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
