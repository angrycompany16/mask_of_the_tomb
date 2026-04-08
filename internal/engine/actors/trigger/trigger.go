package trigger

import (
	"fmt"
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/triggerenv"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/utils"
)

type TriggerState int

const (
	DISJOINT TriggerState = iota
	COLLIDING
)

// Represents an object that raises an event whenever another object
// intersects with this one. Can add masks and stuff as well
type Trigger struct {
	*transform2D.Transform2D
	trigger           *triggerenv.Trigger
	OnCollision       *eventsv2.Event
	OnCollisionEnter  *eventsv2.Event
	OnCollisionExit   *eventsv2.Event
	otherColliderName string
	state             TriggerState
}

func (t *Trigger) Init(cmd *engine.Commands) {
	t.Transform2D.Init(cmd)
	cmd.TriggerEnv().AddTrigger(t.trigger)

	// gX, gY := t.Transform2D.GetPos(false)
	// t.trigger.Rect.SetPos(gX, gY)

	if ok, info := cmd.TriggerEnv().CheckCollision(t.trigger); ok {
		switch t.state {
		case DISJOINT:
			t.state = COLLIDING
			t.otherColliderName = info.OtherName
		}
	}
	fmt.Println(t.state, t.trigger.Name)
}

func (t *Trigger) Update(cmd *engine.Commands) {
	t.Transform2D.Update(cmd)
	if ok, info := cmd.TriggerEnv().CheckCollision(t.trigger); ok {
		switch t.state {
		case DISJOINT:
			// fmt.Println("Enter:", t.trigger.Name, t.otherColliderName)
			t.OnCollision.WithData("otherName", info.OtherName).Raise()
			t.OnCollisionEnter.WithData("otherName", info.OtherName).Raise()
			t.state = COLLIDING
			t.otherColliderName = info.OtherName
		case COLLIDING:
			// fmt.Println("Currently colliding")
			t.OnCollision.WithData("otherName", info.OtherName).Raise()
		}
	} else {
		switch t.state {
		case COLLIDING:
			// fmt.Println("Exit:", t.trigger.Name, t.otherColliderName)
			t.OnCollisionExit.WithData("otherName", t.otherColliderName).Raise()
			t.otherColliderName = ""
			t.state = DISJOINT
		}
	}

	// gX, gY := t.Transform2D.GetPos(false)
	// t.trigger.Rect.SetPos(gX, gY)
}

// Automatically set the rect's position to the global position
// of our object
func (t *Trigger) SetPos(x, y float64) {
	// gX, gY := t.Transform2D.GetPos(false)
	// t.trigger.Rect.SetPos(gX, gY)
}

func (t *Trigger) GetRectPos() (float64, float64) {
	return t.trigger.Rect.Left(), t.trigger.Rect.Top()
}

func defaultTrigger(transform2D *transform2D.Transform2D) *Trigger {
	return &Trigger{
		Transform2D:      transform2D,
		trigger:          triggerenv.NewTrigger(maths.NewRect(0, 0, 8, 8), ""),
		OnCollision:      eventsv2.NewEvent(),
		OnCollisionEnter: eventsv2.NewEvent(),
		OnCollisionExit:  eventsv2.NewEvent(),
	}
}

func NewTrigger(transform2D *transform2D.Transform2D, options ...utils.Option[Trigger]) *Trigger {
	newTrigger := defaultTrigger(transform2D)

	for _, option := range options {
		option(newTrigger)
	}

	return newTrigger
}

func WithRect(rect *maths.Rect) utils.Option[Trigger] {
	return func(t *Trigger) {
		t.trigger.Rect = rect
	}
}

func WithName(name string) utils.Option[Trigger] {
	return func(t *Trigger) {
		t.trigger.Name = name
	}
}
