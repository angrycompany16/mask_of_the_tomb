package events

type EventListener struct {
	notified bool
	event    *Event
}

func NewEventListener(event *Event) *EventListener {
	eventListener := EventListener{notified: false, event: event}
	event.listeners = append(event.listeners, &eventListener)
	return &eventListener
}

func (el *EventListener) Poll() (EventInfo, bool) {
	if el.notified {
		return el.event.info, true
	}
	return EventInfo{}, false
}
