package key

import (
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/missile"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/globaldata"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

// Here it would be nice to just create a bundle. Idea for how it could work:
// Key -- has logic, collision detection etc.
// -> Sprite
// Thats's it

type KeyState int

const (
	IDLE KeyState = iota
	CHASING
)

type Key struct {
	*missile.Missile
	OnPickupEv *events.Event
	OnPickup   *events.EventBus
	Hitbox     *maths.Rect
	EntityIid  string
	DoorIid    string
	PickedUp   bool
	state      KeyState
	Sprite     *sprite.Sprite
}

func (k *Key) Init(cmd *commands.Commands) {
	k.Missile.Init(cmd)
	globaldata, _ := commands.Get[globaldata.GlobalData](cmd)

	if _, ok := globaldata.Persist.Profile.Inventory.Keys[k.EntityIid]; ok {
		k.PickedUp = true
		k.Sprite.Hidden = true
	}
}

func (k *Key) Update(cmd *commands.Commands) {
	k.Missile.Update(cmd)
	if k.PickedUp {
		return
	}

	switch k.state {
	case IDLE:
		if data, raised := k.OnPickup.Poll(); raised {
			k.PickedUp = true
			k.Missile.Active = true
			k.state = CHASING
			collectorTransform := data["Transform"].(*transform2D.Transform2D)
			k.TargetTransform = collectorTransform
		}
	case CHASING:
	}
}

func defaultKey() *Key {
	onPickupEv := events.NewEvent()

	return &Key{
		Missile:    missile.NewMissile(),
		OnPickupEv: onPickupEv,
		OnPickup:   events.NewBusFrom(onPickupEv),
		PickedUp:   false,
		state:      IDLE,
	}
}

func NewKey(entity *ebitenLDTK.Entity) *Key {
	missile := missile.NewMissile(
		missile.WithTransform(
			transform2D.NewTransform2D(transform2D.WithPos(entity.Px[0], entity.Px[1])),
		),
		missile.WithSpeed(4),
		missile.WithActive(false),
		missile.WithCircularOffset(40, 1, maths.RandomRange(0, 10)),
	)

	key := defaultKey()
	key.Missile = missile

	key.EntityIid = entity.Iid

	key.Hitbox = maths.NewRect(entity.Px[0], entity.Px[1], entity.Width, entity.Height)

	opensField := utils.Must(entity.GetFieldByName("Opens"))
	key.DoorIid = ebitenLDTK.As[ebitenLDTK.EntityRef](opensField).EntityIid

	return key
}
