package speechbubble

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/assetloader/assettypes"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SpeechBubble struct {
	graphic     *speechBubbleGraphic
	textDisplay *speechBubbleText
	// vocalizer   *vocalizer
	currentLine int
	lines       []string
}

func (sb *SpeechBubble) Update() {
	sb.graphic.Update()
	newChar, c := sb.textDisplay.Update()
	if newChar {
		Vocalize(string(c))
	}

	sb.textDisplay.x = sb.graphic.rect.Left() + sb.textDisplay.paddingX
	sb.textDisplay.y = sb.graphic.rect.Top() + sb.textDisplay.paddingY

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		sb.LoadActiveLine()
		w, h := sb.textDisplay.Size()
		sb.graphic.targetWidth = w + 2*sb.textDisplay.paddingX
		sb.graphic.targetHeight = h + 2*sb.textDisplay.paddingY
	}
}

func (sb *SpeechBubble) LoadActiveLine() {
	if sb.textDisplay.revealIndex != len(sb.textDisplay.text) {
		sb.textDisplay.revealIndex = len(sb.textDisplay.text) - 1
		return
	} else if sb.currentLine == len(sb.lines) {
		sb.textDisplay.revealIndex = 0
		sb.textDisplay.text = ""
		return
	}
	sb.textDisplay.text = sb.lines[sb.currentLine]
	sb.currentLine++
	sb.textDisplay.revealIndex = 0
	// Expand speech bubble graphic if needed
}

func (sb *SpeechBubble) Draw(ctx rendering.Ctx) {
	sb.graphic.Draw(ctx)
	sb.textDisplay.Draw(ctx)
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
	newSpeechBubble.textDisplay = &speechBubbleText{
		x:        anchorX,
		y:        anchorY,
		paddingX: 32,
		paddingY: 32,
		font:     assetloader.GetFont("JSE_AmigaAMOS"),
		fontSize: 24,

		text:           "",
		revealIndex:    0,
		revealTicker:   time.NewTicker(time.Duration(revealPeriod * float64(time.Second))),
		relLineSpacing: 1.2,
	}
	// newSpeechBubble.vocalizer = newVocalizer()
	newSpeechBubble.lines = []string{
		"Heisann!",
		"...",
		"Lorem ipsum\nmultiline",
		"AAAAAAAAA",
		"Dad? Pop?",
		"Let's talk for real man. \nThis is gonna be a long sentence,\nand I'm sorry about that.",
	}
	return &newSpeechBubble
}
