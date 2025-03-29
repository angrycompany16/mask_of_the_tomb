package events

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
