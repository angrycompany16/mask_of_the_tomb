package events

type Event struct {
	raised    bool
	info      EventInfo
	listeners []*EventListener
}

// TODO: add delayed Raise()
func (e *Event) Raise(info EventInfo) {
	e.info = info
	e.raised = true
}

func NewEvent() *Event {
	event := Event{
		raised:    false,
		listeners: make([]*EventListener, 0),
	}

	_eventManager.events = append(_eventManager.events, &event)
	return &event
}
