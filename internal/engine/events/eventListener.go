package events

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
