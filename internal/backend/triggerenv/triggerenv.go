package triggerenv

import (
	"mask_of_the_tomb/internal/backend/maths"
)

type CollisionInfo struct {
	OtherRect *maths.Rect
	OtherName string
}

// it might be better to just simply give a direct reference to the
// Trigger struct instead
type Trigger struct {
	Rect *maths.Rect
	// ID   string
	Name string
}

func NewTrigger(rect *maths.Rect, name string) *Trigger {
	return &Trigger{
		Rect: rect,
		Name: name,
	}
}

type TriggerEnv struct {
	triggers []*Trigger
}

func (t *TriggerEnv) CheckCollision(trigger *Trigger) (bool, CollisionInfo) {
	for _, otherTrigger := range t.triggers {
		if otherTrigger == trigger {
			continue
		}

		if otherTrigger.Rect.Overlapping(trigger.Rect) {
			return true, CollisionInfo{
				OtherRect: otherTrigger.Rect,
				OtherName: otherTrigger.Name,
			}
		}
	}
	return false, CollisionInfo{}
}

func (t *TriggerEnv) AddTrigger(trigger *Trigger) {
	t.triggers = append(t.triggers, trigger)
}

func NewTriggerEnv() *TriggerEnv {
	return &TriggerEnv{
		triggers: make([]*Trigger, 0),
	}
}
