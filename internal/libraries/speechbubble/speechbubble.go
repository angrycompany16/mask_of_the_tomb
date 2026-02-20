package speechbubble

import (
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
)

// How will this system work?
// main module / interface
// speechbubble struct:
// member variables:
// lines, currentText, anchor (only needs to be used if there is something that
// "owns" the speech bubble), audioPlayer,
// functions:
// Update: Takes care of updating new text, playing sounds,
// keyboard input, etc.

// backend:
// pitchSystem: Generates sound from text
// graphicsSystem: Generates the speech bubble sprite (probably make this
// completely procedural)
// more...?

type SpeechBubble struct {
	graphic *speechBubbleGraphic
}

func (sb *SpeechBubble) Update() {
	sb.graphic.Update()
}

func (sb *SpeechBubble) Draw(ctx rendering.Ctx) {
	sb.graphic.Draw(ctx)
}

func (sb *SpeechBubble) SetAnchor(x, y float64) {
	sb.graphic.anchorX = x
	sb.graphic.anchorY = y
	sb.graphic.tickX = x
	sb.graphic.tickY = y
}

func NewSpeechBubble(anchorX, anchorY, width, height float64) *SpeechBubble {
	newSpeechBubble := SpeechBubble{
		graphic: &speechBubbleGraphic{
			// This will change
			rect:         maths.NewRect(anchorX, anchorY, width, height),
			anchorX:      anchorX,
			anchorY:      anchorY,
			tickX:        anchorX,
			tickY:        anchorY,
			targetWidth:  width,
			targetHeight: height,
		},
	}
	return &newSpeechBubble
}
