package sequence

type ChannelTypes int

type Sequence struct {
	EventChannels map[int]chan any
	Running       bool
	FinishChannel chan int
	action        func(seq *Sequence)
}

func (s *Sequence) Start() {
	if !s.Running {
		// fmt.Println("Running is:", s.Running)
		// fmt.Println("Setting running to true")
		s.Running = true
		go s.action(s)
		// go s.monitor()
	}
}

func (s *Sequence) Stop() {
	s.Running = false
	// Terminate action thread
}

func (s *Sequence) PollEvent(event int) any {
	eventChan, ok := s.EventChannels[event]
	if !ok {
		return nil
	}

	select {
	case event := <-eventChan:
		return event
	default:
		return nil
	}
}

// func (s *Sequence) monitor() {
// 	<-s.FinishChannel
// 	fmt.Println("Stop running")
// 	s.Running = false
// }

func NewSequence(eventChannels map[int]chan any, action func(seq *Sequence)) *Sequence {
	return &Sequence{
		EventChannels: eventChannels,
		Running:       false,
		FinishChannel: make(chan int),
		action:        action,
	}
}
