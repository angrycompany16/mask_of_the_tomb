package trigger

import (
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/triggerenv"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/utils"
)

// Represents an object that raises an event whenever another object
// intersects with this one. Can add masks and stuff as well
type Trigger struct {
	*transform2D.Transform2D
	trigger     *triggerenv.Trigger
	OnCollision *eventsv2.Event
}

func (t *Trigger) Update(cmd *engine.Commands) {
	t.Transform2D.Update(cmd)
	if ok, info := cmd.TriggerEnv().CheckCollision(t.trigger); ok {
		t.OnCollision.WithData("otherName", info.OtherName).Raise()
	}
}

func defaultTrigger(transform2D *transform2D.Transform2D) *Trigger {
	return &Trigger{
		Transform2D: transform2D,
		trigger:     triggerenv.NewTrigger(maths.NewRect(0, 0, 8, 8), ""),
		OnCollision: eventsv2.NewEvent(),
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
