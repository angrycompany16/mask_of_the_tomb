package sound_v2

import (
	"io"
	"io/fs"
	"log"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"slices"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 48000
const requestBufferSize = 4

var requestChan = make(chan audioRequest)
var stopChan = make(chan string)

type SoundData struct {
	Path        string
	Looping     bool
	AudioFormat assettypes.AudioFormat
	QueueSize   int
}

type audioRequest struct {
	name   string
	volume float64
}

func PlaySound(name string, volume float64) {
	requestChan <- audioRequest{name, volume}
}

func SoundServer(
	// Data - Always pass by value so we don't get shared memory errors
	soundCatalogue map[string]SoundData,
) {
	// We will need to uncomment this when we fully deprecate the old sound
	// lib
	// audio.NewContext(sampleRate)
	rqChans := make(map[string]chan audioRequest)
	stopChans := make(map[string]chan int)

	for name, soundData := range soundCatalogue {
		playerChan := make(chan *audio.Player, soundData.QueueSize)
		rqChan := make(chan audioRequest, requestBufferSize)
		stopChan := make(chan int, requestBufferSize)

		rqChans[name] = rqChan
		stopChans[name] = stopChan

		go player(rqChan, stopChan, playerChan)
		go worker(playerChan, soundData.Path, soundData.Looping, soundData.AudioFormat)
	}

	for {
		select {
		case rq := <-requestChan:
			rqChans[rq.name] <- rq
		case name := <-stopChan:
			stopChans[name] <- 1
		}
	}
}

// Player thread that takes out and plays an elt from the queue on request
// Should work ait but idk
func player(
	// Inputs
	requestChan <-chan audioRequest,
	stopChan <-chan int,
	playerChan <-chan *audio.Player,
) {
	players := make([]*audio.Player, 0)

	for {
		select {
		case rq := <-requestChan:
			player := <-playerChan
			player.SetVolume(rq.volume)
			player.Play()
			players = append(players, player)
		case <-stopChan:
			for _, player := range players {
				player.Pause()
			}
		default:
			// Check which players are finished, and remove those.
			activePlayers := make([]*audio.Player, len(players))
			activePlayers = slices.Collect(func(yield func(*audio.Player) bool) {
				for _, player := range players {
					if player.IsPlaying() {
						if !yield(player) {
							return
						}
					}
				}
			})
			players = activePlayers
		}
	}
}

// Worker thread that fills the queue up again
func worker(
	// Outputs
	playerChan chan<- *audio.Player,

	// Data
	path string,
	looping bool,
	format assettypes.AudioFormat,
) {
	// bytes, err := assets.FS.ReadFile(path)
	f, err := assets.FS.Open(path)
	if err != nil {
		log.Fatal("Could not open audio file:", err)
	}

	bytes, err := decodeFile(f, format)
	if err != nil {
		log.Fatal("Could not decode audio file:", err)
	}

	for {
		var player *audio.Player
		player = audio.CurrentContext().NewPlayerF32FromBytes(bytes)
		playerChan <- player
	}
}

func decodeFile(f fs.File, format assettypes.AudioFormat) ([]byte, error) {
	var soundBytes []byte
	var mp3Stream *mp3.Stream
	var oggStream *vorbis.Stream
	var wavStream *wav.Stream
	var err error
	switch format {
	case assettypes.Mp3:
		mp3Stream, err = mp3.DecodeF32(f)
		soundBytes, err = io.ReadAll(mp3Stream)
	case assettypes.Ogg:
		oggStream, err = vorbis.DecodeF32(f)
		soundBytes, err = io.ReadAll(oggStream)
	case assettypes.Wav:
		wavStream, err = wav.DecodeF32(f)
		soundBytes, err = io.ReadAll(wavStream)
	}

	return soundBytes, err
}
