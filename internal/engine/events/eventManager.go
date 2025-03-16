package events

import "mask_of_the_tomb/internal/engine/entities"

var eventMangerSingleton eventManager

type eventManager struct {
	events []*Event
}

func InitEventManager() {
	eventMangerSingleton = eventManager{
		events: make([]*Event, 0),
	}

	entities.RegisterEntity(&eventMangerSingleton, "EventManager")
}

func (evm *eventManager) PreUpdate() {
	for _, event := range evm.events {
		for _, listener := range event.listeners {
			listener.notified = false
		}

		if event.raised {
			for _, listener := range event.listeners {
				listener.notified = true
			}
		}
		event.raised = false
	}
}

func (evm *eventManager) Update() {}
