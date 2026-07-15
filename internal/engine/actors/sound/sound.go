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
	eventBus *events.EventBus
}

func (s *SoundPlayer) Init (cmd *commands.Commands) {
	s.Node.Init(cmd)

	if s.loop {
		sound_v2.AddNewSound(s.name, sound_v2.Loop(s.filePath))
	} else {
		sound_v2.AddNewSound(s.name, sound_v2.Oneshot(s.filePath, 1))
	}
}

func (s *SoundPlayer) Update(cmd *commands.Commands) {
	if s.eventBus == nil {
		return
	}

	if _, ok := s.eventBus.Poll(); ok{
		s.Play()
	}
}

func (s *SoundPlayer) Play() {
	sound_v2.PlaySound(s.name, s.dspChannelName, s.pitchRandomization)
}

func defaultSoundPlayer(node *nodeactor.Node) *SoundPlayer {
	return &SoundPlayer{
		Node: node,
		filePath: "sfx/vietnamese_talking.wav",
		loop: false,
		name: "test-sound",
		dspChannelName:  "master",
		pitchRandomization: 0,
		eventBus: nil,
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

func WithEventBus(event *events.Event) utils.Option[SoundPlayer] {
	return func(s *SoundPlayer) {
		s.eventBus = events.NewBusFrom(event)
	}
}

