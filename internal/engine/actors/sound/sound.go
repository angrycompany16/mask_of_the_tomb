package sound

import (
	"mask_of_the_tomb/internal/backend/events"
	sound_v2 "mask_of_the_tomb/internal/backend/sound"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"
)

type SoundPlayer struct {
	*nodeactor.Node
	filePath string
	loop bool
	name string
	dspChannelName string
	pitchRandomization float64
	autoplay bool
	autostop bool
	startTriggers []*events.EventBus
	endTriggers []*events.EventBus
}

func (s *SoundPlayer) Init (cmd *commands.Commands) {
	s.Node.Init(cmd)

	if s.loop {
		sound_v2.AddNewSound(s.name, sound_v2.Loop(s.filePath))
	} else {
		sound_v2.AddNewSound(s.name, sound_v2.Oneshot(s.filePath, 1))
	}

	if s.autoplay {
		s.Play()
	}
}

func (s *SoundPlayer) Update(cmd *commands.Commands) {
	s.Node.Update(cmd)
	for _, trigger := range s.startTriggers {
		if _, ok := trigger.Poll(); ok{
			s.Play()
		}
	}

	for _, trigger := range s.endTriggers {
		if _, ok := trigger.Poll(); ok{
			s.Stop()
		}
	}
}

func (s *SoundPlayer) OnDestroy(cmd *commands.Commands) {
	s.Node.OnDestroy(cmd)
	if s.autostop {
		s.Stop()
	}
}

func (s *SoundPlayer) Play() {
	sound_v2.PlaySound(s.name, s.dspChannelName, s.pitchRandomization)
}

func (s *SoundPlayer) Stop() {
	sound_v2.StopSound(s.name)
}

func defaultSoundPlayer(node *nodeactor.Node) *SoundPlayer {
	return &SoundPlayer{
		Node: node,
		filePath: "sfx/vietnamese_talking.wav",
		loop: false,
		name: "test-sound",
		dspChannelName:  "master",
		pitchRandomization: 0,
		startTriggers: make([]*events.EventBus, 0),
		endTriggers: make([]*events.EventBus, 0),
		autoplay: false,
		autostop: false,
	}
}

func NewSoundPlayer(node *nodeactor.Node, options ...utils.Option[SoundPlayer]) *SoundPlayer {
	player := defaultSoundPlayer(node)

	for _, option := range options {
		option(player)
	}

	return player
}

func WithSoundData(filePath string, loop bool, name string) utils.Option[SoundPlayer] {
	return func(s *SoundPlayer) {
		s.filePath = filePath
		s.loop = loop
		s.name = name
	}
}

func WithAutoPlay() utils.Option[SoundPlayer] {
	return func(s *SoundPlayer) {
		s.autoplay = true
	}
}

func WithAutoStop() utils.Option[SoundPlayer] {
	return func(s *SoundPlayer) {
		s.autostop = true
	}
}

func WithDspChannel(dspChannelName string) utils.Option[SoundPlayer] {
	return func(s *SoundPlayer) {
		s.dspChannelName = dspChannelName
	}
}

func WithPitchRandomization(pitchRandomization float64) utils.Option[SoundPlayer] {
	return func(s *SoundPlayer) {
		s.pitchRandomization = pitchRandomization
	}
}

func WithStartTriggers(eventlist ...*events.Event) utils.Option[SoundPlayer] {
	return func(s *SoundPlayer) {
		s.startTriggers = make([]*events.EventBus, len(eventlist))

		for i, event := range eventlist {
			s.startTriggers[i] = events.NewBusFrom(event)
		}
	}
}

func WithEndTriggers(eventlist ...*events.Event) utils.Option[SoundPlayer] {
	return func(s *SoundPlayer) {
		s.endTriggers = make([]*events.EventBus, len(eventlist))

		for i, event := range eventlist {
			s.endTriggers[i] = events.NewBusFrom(event)
		}
	}
}
