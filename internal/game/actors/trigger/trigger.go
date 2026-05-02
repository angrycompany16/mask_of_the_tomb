package trigger

import (
	"fmt"
	"image/color"
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/triggerenv"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type TriggerState int

const (
	DISJOINT TriggerState = iota
	COLLIDING
)

// Represents an object that raises an event whenever another object
// intersects with this one. Can add masks and stuff as well
type Trigger struct {
	*graphic.Graphic
	trigger           *triggerenv.Trigger
	OnCollision       *events.Event
	OnCollisionEnter  *events.Event
	OnCollisionExit   *events.Event
	otherColliderName string
	state             TriggerState
	gizmosImage       *ebiten.Image
}

func (t *Trigger) Init(cmd *commands.Commands) {
	t.Graphic.Init(cmd)
	triggerenv, ok := commands.Get[triggerenv.TriggerEnv](cmd)
	if !ok {
		fmt.Println("Missing triggerenv (Trigger)")
		return
	}

	triggerenv.AddTrigger(t.trigger)
	if ok, info := triggerenv.CheckCollision(t.trigger); ok {
		switch t.state {
		case DISJOINT:
			t.state = COLLIDING
			t.otherColliderName = info.OtherName
		}
	}
	vector64.FillRect(t.gizmosImage, 0, 0, t.trigger.Rect.Width, t.trigger.Rect.Height, color.RGBA{40, 100, 100, 50}, false)
}

func (t *Trigger) Update(cmd *commands.Commands) {
	t.Graphic.Update(cmd)

	t.trigger.Rect.SetPos(t.Transform2D.GetPos(false))

	triggerenv, ok := commands.Get[triggerenv.TriggerEnv](cmd)
	if !ok {
		panic("Missing triggerenv (Trigger)")
	}

	// fmt.Println("Checking collisions:", t.trigger.Name)
	if ok, info := triggerenv.CheckCollision(t.trigger); ok {
		switch t.state {
		case DISJOINT:
			t.OnCollision.WithData("otherName", info.OtherName).Raise()
			t.OnCollisionEnter.WithData("otherName", info.OtherName).Raise()
			t.state = COLLIDING
			t.otherColliderName = info.OtherName
		case COLLIDING:
			// fmt.Println("Colliding", t.trigger.Name)
			t.OnCollision.WithData("otherName", info.OtherName).Raise()
		}
	} else {
		switch t.state {
		case COLLIDING:
			t.OnCollisionExit.WithData("otherName", t.otherColliderName).Raise()
			t.otherColliderName = ""
			t.state = DISJOINT
		}
	}
}

func (t *Trigger) DrawGizmo(cmd *commands.Commands) {
	t.Graphic.DrawGizmo(cmd)
	gPosX, gPosY := t.GetPos(false)
	camX, camY := t.GetCamera().WorldToCam(gPosX, gPosY, false)
	cmd.Renderer.Request(opgen.Pos(t.gizmosImage, camX, camY, 0, 0), t.gizmosImage, renderer.RenderTarget{
		renderer.SCREEN,
		"Overlay",
	}, 1)
}

func (t *Trigger) GetRect() maths.Rect {
	return *t.trigger.Rect
}

func defaultTrigger(graphic *graphic.Graphic) *Trigger {
	return &Trigger{
		Graphic:          graphic,
		trigger:          triggerenv.NewTrigger(maths.NewRect(0, 0, 8, 8), ""),
		OnCollision:      events.NewEvent(),
		OnCollisionEnter: events.NewEvent(),
		OnCollisionExit:  events.NewEvent(),
		gizmosImage:      ebiten.NewImage(8, 8),
	}
}

func NewTrigger(graphic *graphic.Graphic, options ...utils.Option[Trigger]) *Trigger {
	newTrigger := defaultTrigger(graphic)

	for _, option := range options {
		option(newTrigger)
	}

	return newTrigger
}

func WithRect(rect *maths.Rect) utils.Option[Trigger] {
	return func(t *Trigger) {
		t.trigger.Rect = rect
		t.gizmosImage = ebiten.NewImage(int(rect.Width), int(rect.Height))
	}
}

func WithName(name string) utils.Option[Trigger] {
	return func(t *Trigger) {
		t.trigger.Name = name
	}
}
