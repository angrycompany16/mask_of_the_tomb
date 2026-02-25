package sound_v2

import (
	"bytes"
	"io"
	"io/fs"
	"log"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/errs"
	resound "mask_of_the_tomb/resound_butwithbytes"
	"sync"
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

var dspChannels = DSPChannels{
	mtx:      &sync.Mutex{},
	channels: make(map[string]*resound.DSPChannel, 0),
}

var playChan = make(chan playRequest)
var stopChan = make(chan string)
var volumeChan = make(chan volumeRequest)
var makeDSPChannelsChan = make(chan []string, 1)
var DSPChannelEffectChan = make(chan effectRequest)

type SoundData struct {
	Path        string
	Looping     bool
	format      AudioFormat
	QueueSize   int
	DSPChannels []string
}

type playRequest struct {
	name           string
	volume         float64
	DSPChannelName string
}

type effectRequest struct {
	name   string
	effect resound.IEffect
}

type volumeRequest struct {
	name   string
	volume float64
}

func PlaySound(name string, volume float64, DSPChannel string) {
	playChan <- playRequest{name, volume, DSPChannel}
}

func MakeDSPChannels(names []string) {
	makeDSPChannelsChan <- names
}

func AddDSPChannelEffect(name string, effect resound.IEffect) {
	DSPChannelEffectChan <- effectRequest{name, effect}
}

// DSP channel effects can be added either when initializing the server or at
// runtime.
// Hmmm. How do we make sure that the players created in worker() are actually
// made to have the correct DSP channel?
// I don't like the current method...
func SoundServer(
	// Data - Always pass by value so we don't get shared memory errors
	soundCatalogue map[string]SoundData,
	DSPChannelNames []string,
) {
	// We will need to uncomment this when we fully deprecate the old sound
	// lib
	// audio.NewContext(sampleRate)
	playChans := make(map[string]chan playRequest)
	dspChannels.MakeDSPChannels(DSPChannelNames)

	for name, soundData := range soundCatalogue {
		playerChan := make(chan *resound.Player, soundData.QueueSize)
		playChan := make(chan playRequest, requestBufferSize)

		playChans[name] = playChan

		go player(playChan, playerChan)
		go worker(playerChan, soundData.Path, soundData.Looping, soundData.format)
	}

	for {
		select {
		case audioRequest := <-playChan:
			playChans[audioRequest.name] <- audioRequest
		case effectRequest := <-DSPChannelEffectChan:
			dspChannels.AddEffect(effectRequest.name, effectRequest.effect)
		}
	}
}

// Player thread that takes out and plays an elt from the queue on request
// Should work ait but idk
// TODO: Important: need to be able to change the volume of looping players on the fly
func player(
	// Inputs
	playChan <-chan playRequest,
	playerChan <-chan *resound.Player,
) {
	for rq := range playChan {
		player := <-playerChan
		dspChannels.AddPlayer(player, rq.DSPChannelName)
		player.Play()

		// So this effect works? Mamma mia
		// Why would this one work, but others not???
		// TEST END
	}
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
			player, err := resound.NewPlayer("KAnjeg", loopedStream)
			// player, err = audio.CurrentContext().NewPlayerF32(loopedStream)
			if err != nil {
				log.Fatal("Infinite loop player failed:", err)
			}
			player.SetBufferSize(500 * time.Millisecond)
			// done := AddStream(player)
			// <-done
		} else {
			player = errs.Must(resound.NewPlayer(0, bytes.NewReader(soundBytes)))
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
