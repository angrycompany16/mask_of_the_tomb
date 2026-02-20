package speechbubble

import (
	"image/color"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

const padding = 4

type speechBubbleGraphic struct {
	rect                      *maths.Rect
	tickX, tickY              float64
	anchorX, anchorY          float64
	targetWidth, targetHeight float64
}

func (sg *speechBubbleGraphic) Update() {
	paddedRect := maths.NewRect(
		sg.rect.Left()-padding,
		sg.rect.Top()-padding,
		sg.rect.Width()+padding*2,
		sg.rect.Height()+padding*2,
	)

	var targetRect maths.Rect
	if !paddedRect.Contains(sg.anchorX, sg.anchorY) {
		newPaddedRect := paddedRect.Reach(sg.anchorX, sg.anchorY)
		targetRect = *maths.NewRect(
			newPaddedRect.Left()+padding,
			newPaddedRect.Top()+padding,
			sg.rect.Width(),
			sg.rect.Height(),
		)
	} else if sg.rect.Contains(sg.anchorX, sg.anchorY) {
		targetRect = sg.rect.Reach(sg.anchorX, sg.anchorY)
	}

	newRect := sg.rect.Lerp(&targetRect, 0.05)
	sg.rect = &newRect
}

func (sg *speechBubbleGraphic) Draw(ctx rendering.Ctx) {
	vector.StrokeRect(ctx.Dst, float32(sg.rect.Left()), float32(sg.rect.Top()), float32(sg.rect.Width()), float32(sg.rect.Height()), 1.0, color.RGBA{255, 0, 0, 255}, false)
	vector.FillCircle(ctx.Dst, float32(sg.tickX), float32(sg.tickY), 2, color.RGBA{0, 255, 0, 255}, true)
}

// Function that sets width / height based on text
func (sg *speechBubbleGraphic) FitText() {

}
