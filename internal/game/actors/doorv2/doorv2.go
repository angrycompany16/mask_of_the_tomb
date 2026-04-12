package doorv2

import (
	"fmt"
	"image/color"
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/trigger"
	"mask_of_the_tomb/internal/game/sceneswitch"
	"mask_of_the_tomb/internal/utils"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	doorV2OtherSideFieldName = "OtherSide"
	doorDirectionFieldName   = "Direction"
)

type DoorV2 struct {
	*graphic.Graphic
	Trigger         *trigger.Trigger         // an entity ref, not inheriting anything
	SpriteTransform *transform2D.Transform2D // an entity ref, not inheriting anything
	//
	// EntityIid          string
	OtherSideLevelIid  string
	OtherSideEntityIid string
	Hitbox             *maths.Rect
	// The trigger child
	// InteractRegion     *maths.Rect
	// sprite             *ebiten.Image
	isReady     bool
	gizmosImage *ebiten.Image
	direction   maths.Direction
	OnCollision *eventsv2.EventBus
}

func (d *DoorV2) Init(cmd *commands.Commands) {
	d.Graphic.Init(cmd)

	slamboxenv, ok := commands.Get[slambox.SlamboxEnvironment](cmd)
	if !ok {
		panic("Missing slambox env (Player)")
	}

	cmd.InputHandler.RegisterAction("DoorInteract", func() bool {
		return inpututil.IsKeyJustPressed(ebiten.KeySpace)
	})
	slamboxenv.AddEnvironmentRect(d.Hitbox)
	d.OnCollision = eventsv2.NewEventBus(d.Trigger.OnCollision)
}

func (d *DoorV2) Update(cmd *commands.Commands) {
	d.Transform2D.Update(cmd)

	if value, raised := d.OnCollision.Poll(); raised && value["otherName"] == "Player" {
		d.SpriteTransform.SetAngle(maths.DirToRadians(maths.Opposite(d.direction)))
		d.isReady = true
	} else {
		d.isReady = false
		d.SpriteTransform.SetAngle(maths.DirToRadians(d.direction))
	}

	if cmd.InputHandler.PollAction("DoorInteract") && d.isReady {
		fmt.Println("Switch scene!")
		// Get the scene switch
		// Set the data
		// Load next scene
		// game.RegisterScene("gameplay", scenes.MakeGamePlayeScene("Level_3"))
		scenemanager, _ := commands.Get[engine.SceneManager](cmd)

		sceneswitch, ok := commands.Get[sceneswitch.SceneSwitch](cmd)
		if !ok {
			panic("Missing scene switch (DoorV2)")
		}
		sceneswitch.SpawnEntityIid = d.OtherSideEntityIid
		scenemanager.SpawnScene(d.OtherSideLevelIid, cmd)
		// sceneswitch.SpawnEntityIid
	}
}

func (d *DoorV2) DrawGizmo(cmd *commands.Commands) {
	d.Graphic.DrawGizmo(cmd)
	d.gizmosImage.Clear()
	vector64.StrokeRect(d.gizmosImage, 0, 0, d.Hitbox.Width()-1, d.Hitbox.Height()-1, 1, color.RGBA{255, 0, 0, 255}, false)

	camX, camY := d.GetCamera().WorldToCam(d.Hitbox.Left(), d.Hitbox.Top(), false)

	cmd.Renderer.Request(opgen.Pos(d.gizmosImage, camX, camY), d.gizmosImage, "Overlay", 0)
}

// Hard-coded for now. Not great but might have to do
func (d *DoorV2) GetSpawnPos() (float64, float64) {
	cx, cy := d.Hitbox.Center()
	switch d.direction {
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

func NewDoorV2(graphic *graphic.Graphic, entity *ebitenLDTK.Entity, levelLDTK *ebitenLDTK.Level) *DoorV2 {
	newDoor := DoorV2{
		Graphic: graphic,
	}

	newDoor.Hitbox = maths.NewRect(
		entity.Px[0],
		entity.Px[1],
		entity.Width,
		entity.Height,
	)

	directionField := utils.Must(entity.GetFieldByName(doorDirectionFieldName))
	newDoor.direction = maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

	doorOtherSideField := utils.Must(entity.GetFieldByName(doorV2OtherSideFieldName))
	doorOtherSide := ebitenLDTK.As[ebitenLDTK.EntityRef](doorOtherSideField)

	newDoor.OtherSideLevelIid = doorOtherSide.LevelIid
	newDoor.OtherSideEntityIid = doorOtherSide.EntityIid

	newDoor.gizmosImage = ebiten.NewImage(int(entity.Width), int(entity.Height))

	return &newDoor
}
