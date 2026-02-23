package sound_v2

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"mask_of_the_tomb/assets"
	resound "mask_of_the_tomb/resound_butwithbytes"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type AudioFormat int

const (
	Mp3 AudioFormat = iota
	Wav
	Ogg
)

const sampleRate = 48000
const requestBufferSize = 4

var playChan = make(chan playRequest)
var stopChan = make(chan string)
var volumeChan = make(chan volumeRequest)

type SoundData struct {
	Path      string
	Looping   bool
	format    AudioFormat
	QueueSize int
}

type playRequest struct {
	name   string
	volume float64
}

type volumeRequest struct {
	name   string
	volume float64
}

func PlaySound(name string, volume float64) {
	playChan <- playRequest{name, volume}
}

func StopSound(name string) {
	stopChan <- name
}

func SetVolume(name string, volume float64) {
	volumeChan <- volumeRequest{name, volume}
}

func SoundServer(
	// Data - Always pass by value so we don't get shared memory errors
	soundCatalogue map[string]SoundData,
) {
	// We will need to uncomment this when we fully deprecate the old sound
	// lib
	// audio.NewContext(sampleRate)
	playChans := make(map[string]chan playRequest)
	stopChans := make(map[string]chan int)
	volumeChans := make(map[string]chan float64)

	for name, soundData := range soundCatalogue {
		playerChan := make(chan *resound.Player, soundData.QueueSize)
		playChan := make(chan playRequest, requestBufferSize)
		stopChan := make(chan int, requestBufferSize)
		volumeChan := make(chan float64, requestBufferSize)

		playChans[name] = playChan
		stopChans[name] = stopChan
		volumeChans[name] = volumeChan

		go player(playChan, stopChan, playerChan, volumeChan)
		go worker(playerChan, soundData.Path, soundData.Looping, soundData.format)
	}

	for {
		select {
		case audioRequest := <-playChan:
			playChans[audioRequest.name] <- audioRequest
		case name := <-stopChan:
			stopChans[name] <- 1
		case volumeRequest := <-volumeChan:
			volumeChans[volumeRequest.name] <- volumeRequest.volume
		}
	}
}

// Player thread that takes out and plays an elt from the queue on request
// Should work ait but idk
// TODO: Important: need to be able to change the volume of looping players on the fly
func player(
	// Inputs
	playChan <-chan playRequest,
	stopChan <-chan int,
	playerChan <-chan *resound.Player,
	volumeChan <-chan float64,
) {
	// players := make([]*audio.Player, 0)

	for rq := range playChan {
		player := <-playerChan
		player.SetVolume(rq.volume)
		player.Play()
	}

	// for {
	// 	select {
	// 	case rq := <-playChan:
	// 		player := <-playerChan
	// 		player.SetVolume(rq.volume)
	// 		player.Play()
	// 		// players = append(players, player)
	// 		// case <-stopChan:
	// 		// 	for _, player := range players {
	// 		// 		player.Pause()
	// 		// 	}
	// 		// case volume := <-volumeChan:
	// 		// 	for _, player := range players {
	// 		// 		player.SetVolume(volume)
	// 		// 	}
	// 		// default:
	// 		// 	// Check which players are finished, and remove those.
	// 		// 	activePlayers := make([]*audio.Player, len(players))
	// 		// 	activePlayers = slices.Collect(func(yield func(*audio.Player) bool) {
	// 		// 		for _, player := range players {
	// 		// 			if player.IsPlaying() {
	// 		// 				if !yield(player) {
	// 		// 					return
	// 		// 				}
	// 		// 			}
	// 		// 		}
	// 		// 	})
	// 		// 	players = activePlayers
	// 	}
	// }
}

// Worker thread that fills the queue up again
func worker(
	// Outputs
	playerChan chan<- *resound.Player,

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

	// Hmmm. If we want to use resound, how will we add our sound players to channels
	// in a thread-safe way?
	// This seems to be a non trivial problem...
	for {
		var player *resound.Player
		if looping {
			byteStream := bytes.NewReader(soundBytes)
			loopedStream := audio.NewInfiniteLoopF32(byteStream, int64(byteStream.Len()))
			player, err := resound.NewPlayer(0, loopedStream)
			// player, err = audio.CurrentContext().NewPlayerF32(loopedStream)
			if err != nil {
				log.Fatal("Infinite loop player failed:", err)
			}
			player.SetBufferSize(500 * time.Millisecond)
			// done := AddStream(player)
			// <-done
		} else {
			player = resound.NewPlayerFromBytes(0, soundBytes)
			player.SetBufferSize(50 * time.Millisecond)
		}
		playerChan <- player
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
		mp3Stream, err = mp3.DecodeF32(f)
		soundBytes, err = io.ReadAll(mp3Stream)
	case Ogg:
		oggStream, err = vorbis.DecodeF32(f)
		soundBytes, err = io.ReadAll(oggStream)
	case Wav:
		wavStream, err = wav.DecodeF32(f)
		soundBytes, err = io.ReadAll(wavStream)
	}

	return soundBytes, err
}

type playerRequest struct {
	player *resound.Player
	done   chan int
}

type effectRequest struct {
	effect resound.IEffect
	name   string
}

var playerRequestChan chan playerRequest
var effectRequestChan chan effectRequest

func AddStream(player *resound.Player) chan int {
	doneChan := make(chan int)
	playerRequestChan <- playerRequest{
		player: player,
		done:   doneChan,
	}
	return doneChan
}

func AddEffect(effect resound.IEffect, name string) {
	effectRequestChan <- effectRequest{
		effect: effect,
		name:   name,
	}
}

// Thread-safe resound DSP channel
func DSPChannel(
	streamChan <-chan playerRequest,
	effectChan <-chan effectRequest,
) {
	dspChannel := resound.NewDSPChannel()

	for {
		select {
		case playerRequest := <-streamChan:
			playerRequest.player.SetDSPChannel(dspChannel)
			playerRequest.done <- 1
		case effectRequest := <-effectChan:
			dspChannel.AddEffect(effectRequest.name, effectRequest.effect)
		}
	}
}

// func newPlayerFromBytes() *resound.Player {
// 	cp := &resound.Player{
// 		id:      id,
// 		Source:  sourceStream,
// 		effects: map[any]IEffect{},
// 	}

// 	player, err := audio.CurrentContext().NewPlayer(cp)

// 	if err != nil {
// 		return nil, err
// 	}

// 	cp.Player = player

// 	return cp, nil
// }
