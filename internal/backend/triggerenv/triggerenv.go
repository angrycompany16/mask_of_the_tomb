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
	// fmt.Println("--- TRIGGER LIST ---")
	// for _, trigger := range t.triggers {
	// 	fmt.Println(trigger.Name)
	// }
	// fmt.Println("--------------------")
	for _, otherTrigger := range t.triggers {
		if otherTrigger == trigger {
			// fmt.Println("Skip:", otherTrigger.Name)
			continue
		}

		if otherTrigger.Rect.Overlapping(trigger.Rect) {
			// fmt.Println("Found collision with", trigger.Name)
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

func (t *TriggerEnv) Reset() {
	t.triggers = make([]*Trigger, 0)
}

func NewTriggerEnv() *TriggerEnv {
	return &TriggerEnv{
		triggers: make([]*Trigger, 0),
	}
}
