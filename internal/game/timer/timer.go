package timer

import (
	"mask_of_the_tomb/internal/game/events"
	"time"
)

type timerState int

const (
	Paused timerState = iota
	Running
	Timedout
)

type Timer struct {
	startTime     time.Time
	timeout       time.Duration
	state         timerState
	TimedoutEvent *events.Event
}

func (t *Timer) Update() {
	switch t.state {
	case Paused:
	case Running:
		if time.Since(t.startTime) > t.timeout {
			t.TimedoutEvent.Raise(events.EventInfo{})
			t.state = Timedout
		}
	case Timedout:
	}
}

func (t *Timer) Reset() {
	t.state = Running
	t.startTime = time.Now()
}

func (t *Timer) Start() {
	t.state = Running
}

func (t *Timer) Pause() {
	t.state = Paused
}

func NewTimer(timeout time.Duration) *Timer {
	return &Timer{
		startTime:     time.Now(),
		timeout:       timeout,
		state:         Running,
		TimedoutEvent: events.NewEvent(),
	}
}
