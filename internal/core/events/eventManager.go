package events

// Maybe consider removing the singleton pattern
// TODO: We will actually do this! It shall be replaced by the event bus.
var (
	_eventManager = eventManager{}
)

type eventManager struct {
	events []*Event
}

func Update() {
	for _, event := range _eventManager.events {
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
