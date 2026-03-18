package eventsv2

type Event struct {
	buses []*EventBus
}

func (e *Event) Raise() *Event {
	// Notify all buses
	for _, bus := range e.buses {
		bus.notif = true
	}
	return e
}

func (e *Event) WithData(name string, data any) *Event {
	// Notify all buses
	for _, bus := range e.buses {
		bus.data[name] = data
	}
	return e
}

func NewEvent() *Event {
	return &Event{
		buses: make([]*EventBus, 0),
	}
}

type EventBus struct {
	notif bool
	data  map[string]any
}

func (e *EventBus) Poll() (map[string]any, bool) {
	if e.notif {
		e.notif = false
		return e.data, true
	}
	return e.data, false
}

func NewEventBus(event *Event) *EventBus {
	eventbus := &EventBus{
		notif: false,
		data:  make(map[string]any),
	}
	event.buses = append(event.buses, eventbus)

	return eventbus
}
