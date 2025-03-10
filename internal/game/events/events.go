package events

var (
	GlobalEventManager = EventManager{}
)

// TODO: Stop the singleton stuff?
// Singleton class
type EventManager struct {
	events []*Event
}

func (em *EventManager) Update() {
	for _, event := range em.events {
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

type Event struct {
	raised    bool
	listeners []*EventListener
}

// TODO: add delayed Raise()
func (e *Event) Raise() {
	e.raised = true
}

func NewEvent() *Event {
	event := Event{
		raised:    false,
		listeners: make([]*EventListener, 0),
	}

	GlobalEventManager.events = append(GlobalEventManager.events, &event)
	return &event
}

type EventListener struct {
	notified bool
}

func NewEventListener(event *Event) *EventListener {
	eventListener := EventListener{notified: false}
	event.listeners = append(event.listeners, &eventListener)
	return &eventListener
}

func (el *EventListener) Poll() bool {
	return el.notified
}
