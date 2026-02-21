package speechbubble

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"time"
)

type SpeechBubble struct {
	graphic *speechBubbleGraphic
	text    *speechBubbleText
}

func (sb *SpeechBubble) Update() {
	sb.graphic.Update()
	sb.text.Update()

	sb.text.x = sb.graphic.rect.Left() + sb.text.paddingX
	sb.text.y = sb.graphic.rect.Top() + sb.text.paddingY
}

func (sb *SpeechBubble) Draw(ctx rendering.Ctx) {
	sb.graphic.Draw(ctx)
	sb.text.Draw(ctx)
}

func (sb *SpeechBubble) SetAnchor(x, y float64) {
	sb.graphic.anchorX = x
	sb.graphic.anchorY = y
}

func NewSpeechBubble(anchorX, anchorY, width, height float64) *SpeechBubble {
	speechBubbleTick := errs.Must(assettypes.GetImageAsset("textBoxTickSprite"))
	newSpeechBubble := SpeechBubble{}

	newSpeechBubble.graphic = &speechBubbleGraphic{
		// This will change
		// Or will it...
		rect:         maths.NewRect(anchorX, anchorY, width*rendering.PIXEL_SCALE, height*rendering.PIXEL_SCALE),
		anchorX:      anchorX,
		anchorY:      anchorY,
		tickX:        anchorX,
		tickY:        anchorY,
		targetWidth:  width * rendering.PIXEL_SCALE,
		targetHeight: height * rendering.PIXEL_SCALE,
		tickSprite:   speechBubbleTick,
	}
	newSpeechBubble.text = &speechBubbleText{
		x:        anchorX,
		y:        anchorY,
		paddingX: 32,
		paddingY: 32,
		font:     assetloader.GetFont("JSE_AmigaAMOS"),
		fontSize: 32,
		lines: []string{
			"HEisann!",
			"På degsann!",
			"Lorem ipsum",
		},
		text:         "",
		revealIndex:  0,
		revealTicker: time.NewTicker(time.Duration(revealPeriod * float64(time.Second))),
	}
	return &newSpeechBubble
}
