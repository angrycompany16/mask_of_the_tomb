package resetlistener

import (
	"mask_of_the_tomb/internal/backend/node"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/game/actors/slamboxactor"
	"mask_of_the_tomb/internal/game/gamestate"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ResetListener struct {
	*nodeactor.Node
}

func (r *ResetListener) Init(cmd *commands.Commands) {
	r.Node.Init(cmd)

	cmd.InputHandler.InputSchemes["PlayerControls"].AddBinding("Reset", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeyR)
	})
}

func (r *ResetListener) Update(cmd *commands.Commands) {
	r.Node.Update(cmd)

	if cmd.InputHandler.InputSchemes["PlayerControls"].PollAction("Reset") {
		scene, _ := commands.Get[engine.Scene](cmd)
		gamestate, _ := commands.Get[gamestate.GameState](cmd)
		levelstate := gamestate.LevelStates[scene.GetName()]

		slamboxes := scene.GetRoot().GetChildrenFunc(
			func(n *node.Node[engine.Actor]) bool {
				_, ok := engine.As[*slamboxactor.Slambox](n.GetValue())
				return ok
			},
		)

		for _, slambox := range slamboxes {
			slamboxactor, _ := engine.As[*slamboxactor.Slambox](slambox.GetValue())
			if player, ok := engine.As[*player.Player](slambox.GetValue()); ok {
				player.SetPos(levelstate.PlayerSpawnPos.X, levelstate.PlayerSpawnPos.Y)
				player.Direction = levelstate.PlayerSpawnDir
				continue
			}
			slamboxPos := levelstate.SlamboxPositions[slambox.GetID()]
			slamboxactor.SetPos(slamboxPos.X, slamboxPos.Y)
		}
	}
}

func NewResetListener(node *nodeactor.Node) *ResetListener {
	return &ResetListener{
		Node: node,
	}
}
