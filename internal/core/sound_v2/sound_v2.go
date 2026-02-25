package sound_v2

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/errs"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/solarlune/resound"
	"github.com/solarlune/resound/effects"
)

type AudioFormat int

const (
	Mp3 AudioFormat = iota
	Wav
	Ogg
)

const sampleRate = 48000
const requestBufferSize = 4

var dspChannels = DSPChannels{
	mtx:      &sync.Mutex{},
	channels: make(map[string]*resound.DSPChannel, 0),
}

var playRequestChan = make(chan playRequest)
var stopRequestChan = make(chan string)
var dspChannelEffectChan = make(chan effectRequest)
var dspChannelEditEffectChan = make(chan editEffectRequest)

type SoundData struct {
	Path        string
	Looping     bool
	format      AudioFormat
	QueueSize   int
	DSPChannels []string
}

type playRequest struct {
	name               string
	DSPChannelName     string
	pitchRandomization float64
}

type effectRequest struct {
	channelName string
	effectName  string
	effect      resound.IEffect
}

type editEffectRequest struct {
	channelName string
	effectName  string
	action      func(effect resound.IEffect) error
}

type finishedPlayer struct {
	p       *resound.Player
	looping bool
}

func PlaySound(name string, DSPChannel string, pitchRandomization float64) {
	playRequestChan <- playRequest{name, DSPChannel, pitchRandomization}
}

// TODO: Add a small fadeout when stopping a sound
func StopSound(name string) {
	stopRequestChan <- name
}

func AddDSPChannelEffect(channelName string, effectName string, effect resound.IEffect) {
	dspChannelEffectChan <- effectRequest{channelName, effectName, effect}
}

// Be a bit careful with this. action is a method that will be evaluated
// on the desired effect, so one should be careful not to break things or cause race conditions
func EditDSPChannelEffect(channelName string, effectName string, action func(effect resound.IEffect) error) {
	dspChannelEditEffectChan <- editEffectRequest{channelName, effectName, action}
}

func SoundServer(
	// Data - Always pass by value so we don't get shared memory errors
	soundCatalogue map[string]SoundData,
	DSPChannelNames []string,
) {
	audio.NewContext(sampleRate)
	playChans := make(map[string]chan playRequest)
	stopChans := make(map[string]chan int)
	dspChannels.makeDSPChannels(DSPChannelNames)

	for name, soundData := range soundCatalogue {
		playerChan := make(chan finishedPlayer, soundData.QueueSize)
		playChan := make(chan playRequest, requestBufferSize)
		stopChan := make(chan int)

		playChans[name] = playChan
		stopChans[name] = stopChan

		go player(playChan, playerChan, stopChan)
		go worker(playerChan, soundData.Path, soundData.Looping, soundData.format)
	}

	for {
		select {
		case audioRequest := <-playRequestChan:
			if playChans[audioRequest.name] == nil {
				fmt.Printf("player named [%s] not found!\n", audioRequest.name)
				continue
			}
			playChans[audioRequest.name] <- audioRequest
		case effectRequest := <-dspChannelEffectChan:
			dspChannels.addEffect(effectRequest.channelName, effectRequest.effectName, effectRequest.effect)
		case editEffectRequest := <-dspChannelEditEffectChan:
			dspChannels.editEffect(editEffectRequest.channelName, editEffectRequest.effectName, editEffectRequest.action)
		case name := <-stopRequestChan:
			if stopChans[name] == nil {
				fmt.Printf("player named [%s] not found!\n", name)
				continue
			}
			stopChans[name] <- 1
		}
	}
}

func player(
	// Inputs
	playChan <-chan playRequest,
	playerChan <-chan finishedPlayer,
	stopChan <-chan int,
) {
	isPlaying := false
	var activePlayer *resound.Player
	// For now - Let's only enable pausing for
	// infinite loop players
	for {
		select {
		case rq := <-playChan:
			finishedPlayer := <-playerChan
			player := finishedPlayer.p

			if finishedPlayer.looping {
				if isPlaying {
					return
				}
				isPlaying = true
			}
			dspChannels.addPlayer(player, rq.DSPChannelName)

			pitchShift := computePitchShift(rq.pitchRandomization)
			pitchShiftEffect := effects.NewPitchShift(2048).SetSource(player.Source).SetPitch(pitchShift)
			player.AddEffect("pitch", pitchShiftEffect)

			player.Play()
			activePlayer = player
		case <-stopChan:
			if !isPlaying {
				fmt.Println("Attempted to pause a non-playing player.")
				continue
			}
			isPlaying = false
			activePlayer.Pause()
		}
	}
}

func worker(
	// Outputs
	playerChan chan<- finishedPlayer,

	// Data
	path string,
	looping bool,
	format AudioFormat,
) {
	f, err := assets.FS.Open(path)
	if err != nil {
		log.Fatal("Could not open audio file:", err)
	}

	soundBytes, err := decodeFile(f, format)
	if err != nil {
		log.Fatal("Could not decode audio file:", err)
	}

	for {
		var player *resound.Player
		if looping {
			byteStream := bytes.NewReader(soundBytes)
			loopedStream := audio.NewInfiniteLoop(byteStream, int64(byteStream.Len()))

			player, err = resound.NewPlayer(0, loopedStream)
			if err != nil {
				log.Fatal("Infinite loop player failed:", err)
			}

			player.SetBufferSize(500 * time.Millisecond)
		} else {
			player = errs.Must(resound.NewPlayer(0, bytes.NewReader(soundBytes)))
			player.SetBufferSize(50 * time.Millisecond)
		}

		playerChan <- finishedPlayer{
			p:       player,
			looping: looping,
		}
	}
}

func decodeFile(f fs.File, format AudioFormat) ([]byte, error) {
	var soundBytes []byte
	var mp3Stream *mp3.Stream
	var oggStream *vorbis.Stream
	var wavStream *wav.Stream
	var err error
	switch format {
	case Mp3:
		mp3Stream, err = mp3.DecodeWithoutResampling(f)
		soundBytes, err = io.ReadAll(mp3Stream)
	case Ogg:
		oggStream, err = vorbis.DecodeWithoutResampling(f)
		soundBytes, err = io.ReadAll(oggStream)
	case Wav:
		wavStream, err = wav.DecodeWithoutResampling(f)
		soundBytes, err = io.ReadAll(wavStream)
	}

	return soundBytes, err
}

func computePitchShift(randomization float64) float64 {
	return 1.0 + randomization*(2*rand.Float64()-1)
}

type DSPChannels struct {
	mtx      *sync.Mutex
	channels map[string]*resound.DSPChannel
}

func (d *DSPChannels) addPlayer(player *resound.Player, channelName string) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.channels[channelName] == nil {
		fmt.Println("Failed to assign player to channel.")
		fmt.Printf("Channel with name %s does not exist!\n", channelName)
		return
	}
	// fmt.Println(player)
	// fmt.Println(d.channels[channelName])
	player.SetDSPChannel(d.channels[channelName])
}

func (d *DSPChannels) makeDSPChannels(names []string) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	for _, name := range names {
		d.channels[name] = resound.NewDSPChannel()
	}
}

func (d *DSPChannels) addEffect(channelName string, effectName string, effect resound.IEffect) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.channels[channelName] == nil {
		fmt.Println("Failed to add effect to channel.")
		fmt.Printf("Channel with name %s does not exist!\n", channelName)
		return
	}
	d.channels[channelName].AddEffect(effectName, effect)
}

func (d *DSPChannels) editEffect(channelName string, effectName string, action func(resound.IEffect) error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.channels[channelName] == nil {
		fmt.Println("Failed to edit effect on channel.")
		fmt.Printf("Channel with name %s does not exist!\n", channelName)
		return
	}
	effect := d.channels[channelName].Effects[effectName]
	if effect == nil {
		fmt.Println("Failed to edit effect on channel.")
		fmt.Printf("Channel %s has no effect named %s!\n", channelName, effectName)
		return
	}
	err := action(effect)
	if err != nil {
		fmt.Println("Effect editing failed with error:", err)
	}
}
