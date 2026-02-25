package sound_v2

import (
	resound "mask_of_the_tomb/resound_butwithbytes"
	"sync"
)

type DSPChannels struct {
	mtx      *sync.Mutex
	channels map[string]*resound.DSPChannel
}

func (d *DSPChannels) GetDSPChannel(name string) *resound.DSPChannel {
	d.mtx.Lock()
	return d.channels[name]
}

func (d *DSPChannels) ReturnDSPChannel(name string) {
	d.mtx.Unlock()
}

func (d *DSPChannels) AddPlayer(player *resound.Player, name string) {
	d.mtx.Lock()
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

func (d *DSPChannels) AddEffect(channelName string, effectName string, effect resound.IEffect) {
	d.mtx.Lock()
	d.channels[channelName].AddEffect(effectName, effect)
	d.mtx.Unlock()
}
