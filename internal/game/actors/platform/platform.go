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

type Platform struct {
	*graphic.Graphic
	Hitbox       *maths.Rect
	direction    maths.Direction
	OnPlayerMove *eventsv2.EventBus
	active       bool
}

func (p *Platform) Init(cmd *commands.Commands) {
	p.Graphic.Init(cmd)
	slamboxenv, ok := commands.Get[slambox.SlamboxEnvironment](cmd)
	if !ok {
		panic("Missing slambox env (Platform)")
	}

	slamboxenv.AddEnvironmentRect(p.Hitbox)

	scene, _ := commands.Get[engine.Scene](cmd)
	playerNode, ok := engine.GetNodeByType[*player.Player](scene)
	if !ok {
		panic("Unable to find player")
	}

	playerActor, _ := engine.As[*player.Player](playerNode.GetValue())
	p.OnPlayerMove = eventsv2.NewEventBus(playerActor.OnMove)
}

func (p *Platform) Update(cmd *commands.Commands) {
	p.Graphic.Update(cmd)
	if data, raised := p.OnPlayerMove.Poll(); raised {
		direction := data["Direction"]
		slamboxenv, _ := commands.Get[slambox.SlamboxEnvironment](cmd)
		if direction == p.direction {
			p.active = false
			slamboxenv.RemoveEnvironmentRect(p.Hitbox)
		} else if p.active == false {
			p.active = true
			slamboxenv.AddEnvironmentRect(p.Hitbox)
		}
	}
}

func NewPlatform(
	graphic *graphic.Graphic,
	entity *ebitenLDTK.Entity,
) *Platform {
	newPlatform := &Platform{
		Graphic: graphic,
		active:  true,
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
