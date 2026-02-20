package speechbubble

import (
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
)

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
}

func NewSpeechBubble(anchorX, anchorY, width, height float64) *SpeechBubble {
	speechBubbleTick := errs.Must(assettypes.GetImageAsset("textBoxTickSprite"))
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
			tickSprite:   speechBubbleTick,
		},
	}
	return &newSpeechBubble
}
