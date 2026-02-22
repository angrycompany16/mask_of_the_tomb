package sound

import (
	// Illegal reference - Need to move this file out of core.
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 48000

// Fuck it- process oriented
func SoundServer(
	// Outputs
	audioQueues map[string]chan<- *audio.Player,

) {
	audio.NewContext(sampleRate)
	for name, queue := range audioQueues {
		// Later we will pick different buffer sizes for the channels
		// i.e. 1 for looping sounds, maybe 2 for non-looping ones
		go generatePlayers(queue, name, 1)
	}

	for {

	}
}

// We assume that audiocontext is thread-safe and performant
func generatePlayers(
	// Inputs
	assetChan <-chan assetloader.Asset,

	// Outputs
	queue chan<- *audio.Player,
	assetRequest chan<- string,

	path string,
	bufferSize int,
	audioFormat assettypes.AudioFormat) {

	// hm. This was painful
	audioCtx := audio.CurrentContext()
	// assetRequest <- path
	// asset := <- assetloader.assetChan

	// Convert asset to the correct type of stream
	for {
		// Now we can make our audio players. Nice

		// queue <-
	}
}
