package platform

import (
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

// TODO: Consider using an autotile-setup for the graphics

type Platform struct {
	*graphic.Graphic
	Hitbox             *maths.Rect
	direction          maths.Direction
	OnPlayerMove       *eventsv2.EventBus
	OnPlayerMoveFinish *eventsv2.EventBus
	active             bool
}

func (p *Platform) Init(cmd *commands.Commands) {
	p.Graphic.Init(cmd)

	scene, _ := commands.Get[engine.Scene](cmd)
	playerNode, ok := engine.GetNodeByType[*player.Player](scene)
	if !ok {
		panic("Unable to find player")
	}

	playerActor, _ := engine.As[*player.Player](playerNode.GetValue())
	p.OnPlayerMove = eventsv2.NewEventBus(playerActor.OnMove)
	p.OnPlayerMoveFinish = eventsv2.NewEventBus(playerActor.OnMoveFinishEv)
}

// TODO: Bug: When the player uses a buffered input, OnPlayerMove
// is actually called *before* OnPlayerMoveFinish, which leads to
// a bug where the player can go through more platforms than he is
// supposed to.
func (p *Platform) Update(cmd *commands.Commands) {
	p.Graphic.Update(cmd)
	slamboxenv, _ := commands.Get[slambox.SlamboxEnvironment](cmd)
	if data, raised := p.OnPlayerMove.Poll(); raised {
		direction := data["Direction"]
		if direction == maths.Opposite(p.direction) {
			p.active = true
			slamboxenv.AddEnvironmentRect(p.Hitbox)
		} else if p.active == true {
			p.active = false
			slamboxenv.RemoveEnvironmentRect(p.Hitbox)
		}
	}

	if _, raised := p.OnPlayerMoveFinish.Poll(); raised {
		if p.active == true {
			slamboxenv.RemoveEnvironmentRect(p.Hitbox)
		}
	}
}

func NewPlatform(
	graphic *graphic.Graphic,
	entity *ebitenLDTK.Entity,
) *Platform {
	newPlatform := &Platform{
		Graphic: graphic,
		active:  false,
	}
	newPlatform.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)

	directionField := utils.Must(entity.GetFieldByName("Direction"))
	newPlatform.direction = maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

	return newPlatform
}
