package events

type Event struct {
	raised    bool
	listeners []*EventListener
}

// TODO: GetEventByName() or something to allow for events to not require importing packages
// TODO: allow event to also include some data
// TODO: add delayed Raise()
func (e *Event) Raise() {
	e.raised = true
}

func NewEvent() *Event {
	event := Event{
		raised:    false,
		listeners: make([]*EventListener, 0),
	}

	eventMangerSingleton.events = append(eventMangerSingleton.events, &event)
	return &event
}
