package sound_v2

import (
	"fmt"
	resound "mask_of_the_tomb/resound_butwithbytes"
	"sync"
)

type DSPChannels struct {
	mtx      *sync.Mutex
	channels map[string]*resound.DSPChannel
}

// But wait - This can cause stuttering
// If sounds are occupying the DSP channels, and we try to add effects
// (in the main thread), then the game stutters...
// This means that adding effects must not be done in the main thread.
// But that should be possible to achieve!
func (d *DSPChannels) GetDSPChannel(name string) *resound.DSPChannel {
	d.mtx.Lock()
	return d.channels[name]
}

// We may also consider making one mutex per DSP channel..
func (d *DSPChannels) ReturnDSPChannel(name string) {
	d.mtx.Unlock()
}

// A problem - or is it a problem? Hm.
func (d *DSPChannels) AddPlayer(player *resound.Player, name string) {
	d.mtx.Lock()
	fmt.Println(d.channels[name])
	player.SetDSPChannel(d.channels[name])
	d.mtx.Unlock()
}

func (d *DSPChannels) MakeDSPChannels(names []string) {
	d.mtx.Lock()
	for _, name := range names {
		d.channels[name] = resound.NewDSPChannel()
	}
	d.mtx.Unlock()
}

func (d *DSPChannels) AddEffect(name string, effect resound.IEffect) {
	d.mtx.Lock()
	fmt.Println("Add effect")
	d.channels[name].AddEffect(0, effect)
	d.mtx.Unlock()
}
