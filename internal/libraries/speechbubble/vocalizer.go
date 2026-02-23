package speechbubble

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/sound"
	"mask_of_the_tomb/internal/libraries/sound_v2"
	"strings"
)

type vocalizer struct {
	playerA sound.EffectPlayer
	playerE sound.EffectPlayer
	playerI sound.EffectPlayer
	playerO sound.EffectPlayer
	playerU sound.EffectPlayer
}

func (v *vocalizer) Vocalize(c string) {
	switch strings.ToLower(c) {
	case "a":
		sound_v2.PlaySound("vowelA", 1.0)
		// v.playerA.Play()
	case "e":
		v.playerE.Play()
	case "i":
		v.playerI.Play()
	case "o":
		v.playerO.Play()
	case "u":
		v.playerU.Play()
	}
}

func newVocalizer() *vocalizer {
	streamA := errs.Must(assettypes.GetWavStream("vowelA"))
	streamE := errs.Must(assettypes.GetWavStream("vowelE"))
	streamI := errs.Must(assettypes.GetWavStream("vowelI"))
	streamO := errs.Must(assettypes.GetWavStream("vowelO"))
	streamU := errs.Must(assettypes.GetWavStream("vowelU"))
	newVocalizer := vocalizer{
		playerA: sound.EffectPlayer{errs.Must(sound.FromStream(streamA)), 1.0},
		playerE: sound.EffectPlayer{errs.Must(sound.FromStream(streamE)), 1.0},
		playerI: sound.EffectPlayer{errs.Must(sound.FromStream(streamI)), 1.0},
		playerO: sound.EffectPlayer{errs.Must(sound.FromStream(streamO)), 1.0},
		playerU: sound.EffectPlayer{errs.Must(sound.FromStream(streamU)), 1.0},
	}
	return &newVocalizer
}
