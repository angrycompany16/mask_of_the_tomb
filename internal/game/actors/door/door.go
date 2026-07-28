package door

import (
	"fmt"
	"image/color"
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/backend/triggerenv"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/animatedsprite"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/trigger"
	"mask_of_the_tomb/internal/game/globaldata"
	"mask_of_the_tomb/internal/game/sceneswitch"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DoorState int

const (
	IDLE DoorState = iota
	CLOSING
	OPENING
)

const (
	otherSideFieldName = "OtherSide"
	directionFieldName = "Direction"
)

type Door struct {
	*graphic.Graphic
	Trigger            *trigger.Trigger
	SpriteTransform    *transform2D.Transform2D
	AnimatedSprite     *animatedsprite.AnimatedSprite
	LockAnim           *animatedsprite.AnimatedSprite
	EntityIid          string
	OtherSideLevelIid  string
	OtherSideEntityIid string
	Hitbox             *maths.Rect
	isReady            bool
	gizmosImage        *ebiten.Image
	Direction          maths.Direction
	OnCollision        *events.EventBus
	OnClipFinished     *events.EventBus
	OnOpen             *events.Event
	OnTryOpenLocked    *events.Event
	OnUnlockEv         *events.Event
	State              DoorState
	biome              string
	Locked             bool
}

func (d *Door) Init(cmd *commands.Commands) {
	d.Graphic.Init(cmd)
	sceneswitch, _ := commands.Get[sceneswitch.SceneSwitch](cmd)

	if sceneswitch.SpawnEntityIid == d.EntityIid {
		d.State = CLOSING
	}

	globaldata, _ := commands.Get[globaldata.GlobalData](cmd)

	if unlocked, ok := globaldata.Persist.Profile.UnlockedDoors[d.EntityIid]; ok {
		d.Locked = !unlocked
		fmt.Println("Setting locked state to", !unlocked)
	}

	if !d.Locked {
		d.LockAnim.Hidden = true
	}

	d.OnClipFinished = events.NewBusFrom(d.AnimatedSprite.OnClipFinished)

	slamboxenv, _ := commands.Get[slambox.SlamboxEnvironment](cmd)

	playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
	playerControls.RegisterAction("DoorInteract", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeySpace)
	})
	slamboxenv.AddEnvironmentRect(d.Hitbox)
	d.OnCollision = events.NewBusFrom(d.Trigger.OnCollision)
}

func (d *Door) Update(cmd *commands.Commands) {
	d.Transform2D.Update(cmd)

	switch d.State {
	case CLOSING:
		d.AnimatedSprite.SwitchClip("Close")
		if _, raised := d.OnClipFinished.Poll(); raised {
			d.State = IDLE
			cmd.InputHandler.InputSchemes["PlayerControls"].Active = true
		}
	case OPENING:
		if value, raised := d.OnClipFinished.Poll(); raised && value["clip"] == "Open" {
			fmt.Println("Switch scene!")
			scenemanager, _ := commands.Get[engine.SceneManager](cmd)

			sceneswitch, ok := commands.Get[sceneswitch.SceneSwitch](cmd)
			if !ok {
				panic("Missing scene switch (Door)")
			}
			sceneswitch.SpawnEntityIid = d.OtherSideEntityIid
			sceneswitch.SpawnDirection = maths.Opposite(d.Direction)
			sceneswitch.PreviousBiome = d.biome
			// TODO: There is a much better way to do this - Include an OnDestroy method that gets called whenever a scene gets destroyed.
			slamboxenv, _ := commands.Get[slambox.SlamboxEnvironment](cmd)
			slamboxenv.Reset()

			triggerenv, _ := commands.Get[triggerenv.TriggerEnv](cmd)
			triggerenv.Reset()

			scenemanager.SpawnScene(d.OtherSideLevelIid, cmd)
		}
	case IDLE:
		d.AnimatedSprite.SwitchClip("Idle")
		if value, raised := d.OnCollision.Poll(); raised && value["otherName"] == "Player" {
			d.isReady = true
		} else {
			d.isReady = false
		}

		playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
		if playerControls.PollAction("DoorInteract") && d.isReady {
			if d.Locked {
				d.OnTryOpenLocked.Raise()
				d.LockAnim.SwitchClip("TryUnlock")
				d.LockAnim.Restart()
			} else if !d.Locked {
				d.OnOpen.Raise()
				d.AnimatedSprite.SwitchClip("Open")
				d.State = OPENING
			}
		}

		if playerControls.PollAction("Use") && d.isReady {
			globaldata, _ := commands.Get[globaldata.GlobalData](cmd)
			inventory := globaldata.Persist.Profile.Inventory
			var keyIid string
			haskeyThisSide, keyIidThisSide := inventory.HasKey(d.EntityIid)
			hasKeyOtherSide, keyIidOtherSide := inventory.HasKey(d.OtherSideEntityIid)
			if haskeyThisSide {
				keyIid = keyIidThisSide
			} else if hasKeyOtherSide {
				keyIid = keyIidOtherSide
			}

			if d.Locked && (haskeyThisSide || hasKeyOtherSide) {
				d.OnUnlockEv.Raise()
				d.Locked = false
				globaldata.Persist.Profile.UnlockedDoors[d.EntityIid] = true
				globaldata.Persist.Profile.UnlockedDoors[d.OtherSideEntityIid] = true
				d.LockAnim.SwitchClip("Unlock")
				inventory.Keys[keyIid].Used = true
			}
		}
	}
}

func (d *Door) DrawGizmo(cmd *commands.Commands) {
	d.Graphic.DrawGizmo(cmd)
	d.gizmosImage.Clear()
	vector64.StrokeRect(d.gizmosImage, 0, 0, d.Hitbox.Width-1, d.Hitbox.Height-1, 1, color.RGBA{255, 0, 0, 255}, false)

	camX, camY := d.GetCamera().WorldToCam(d.Hitbox.Left(), d.Hitbox.Top(), false)

	cmd.Renderer.Request(opgen.Pos(d.gizmosImage, camX, camY), d.gizmosImage, renderer.RenderTarget{
		Type: renderer.SCREEN,
		Name: "Overlay",
	}, 0)
}

// Hard-coded for now. Not great but might have to do
func (d *Door) GetSpawnPos() (float64, float64) {
	cx, cy := d.Hitbox.Center()
	switch d.Direction {
	case maths.DirUp:
		return cx - 8, d.Hitbox.Top() - 16
	case maths.DirDown:
		return cx - 8, d.Hitbox.Bottom()
	case maths.DirLeft:
		return d.Hitbox.Left() - 16, cy - 8
	case maths.DirRight:
		return d.Hitbox.Right(), cy - 8
	}
	return 0, 0
}

func NewDoor(graphic *graphic.Graphic, entity *ebitenLDTK.Entity, levelLDTK *ebitenLDTK.Level) *Door {
	newDoor := Door{
		Graphic:         graphic,
		EntityIid:       entity.Iid,
		State:           IDLE,
		OnOpen:          events.NewEvent(),
		OnTryOpenLocked: events.NewEvent(),
		OnUnlockEv:      events.NewEvent(),
	}

	newDoor.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)

	biomeField := utils.Must(levelLDTK.GetFieldByName("Biome"))
	newDoor.biome = ebitenLDTK.As[ebitenLDTK.Enum](biomeField).Value

	directionField := utils.Must(entity.GetFieldByName(directionFieldName))
	newDoor.Direction = maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

	lockedField := utils.Must(entity.GetFieldByName("Locked"))
	newDoor.Locked = ebitenLDTK.As[bool](lockedField)

	doorOtherSideField := utils.Must(entity.GetFieldByName(otherSideFieldName))
	doorOtherSide := ebitenLDTK.As[ebitenLDTK.EntityRef](doorOtherSideField)

	newDoor.OtherSideLevelIid = doorOtherSide.LevelIid
	newDoor.OtherSideEntityIid = doorOtherSide.EntityIid

	newDoor.gizmosImage = ebiten.NewImage(int(entity.Width), int(entity.Height))

	return &newDoor
}
