package key

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/gamestate"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

// Here it would be nice to just create a bundle. Idea for how it could work:
// Key -- has logic, collision detection etc.
// -> Sprite
// Thats's it

type Key struct {
	*transform2D.Transform2D
	OnPickupEv *events.Event
	OnPickup   *events.EventBus
	Hitbox     *maths.Rect
	EntityIid  string
	DoorIid    string
	PickedUp   bool
}

func (k *Key) Init(cmd *commands.Commands) {
	k.Transform2D.Init(cmd)
	gamestate, _ := commands.Get[gamestate.GameState](cmd)

	if _, ok := gamestate.Inventory.Keys[k.EntityIid]; ok {
		k.PickedUp = true
	}
}

func (k *Key) Update(cmd *commands.Commands) {
	k.Transform2D.Update(cmd)
	if k.PickedUp {
		return
	}

	if _, raised := k.OnPickup.Poll(); raised {
		fmt.Println("I was picked up!")
		k.PickedUp = true
	}
}

func defaultKey(transform2D *transform2D.Transform2D) *Key {
	onPickupEv := events.NewEvent()

	return &Key{
		Transform2D: transform2D,
		OnPickupEv:  onPickupEv,
		OnPickup:    events.NewBusFrom(onPickupEv),
		PickedUp:    false,
	}
}

func NewKey(entity *ebitenLDTK.Entity) *Key {
	transform := transform2D.NewTransform2D(transform2D.WithPos(entity.Px[0], entity.Px[1]))

	key := defaultKey(transform)

	key.EntityIid = entity.Iid

	key.Hitbox = maths.NewRect(entity.Px[0], entity.Px[1], entity.Width, entity.Height)

	opensField := utils.Must(entity.GetFieldByName("Opens"))
	key.DoorIid = ebitenLDTK.As[ebitenLDTK.EntityRef](opensField).EntityIid

	return key
}
