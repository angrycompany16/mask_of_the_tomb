package speechbubble

import (
	"image/color"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const padding = 4

var interiorColor = color.RGBA{21, 10, 31, 255}
var borderColor = color.RGBA{255, 253, 240, 255}

type speechBubbleGraphic struct {
	rect                      *maths.Rect
	tickX, tickY              float64
	anchorX, anchorY          float64
	targetWidth, targetHeight float64
	tickSprite                *ebiten.Image
	tickRotation              maths.Direction
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
		newRect := sg.rect.Lerp(&targetRect, 0.05)
		sg.rect = &newRect
	} else if sg.rect.Contains(sg.anchorX, sg.anchorY) {
		targetRect = sg.rect.Reach(sg.anchorX, sg.anchorY)
		newRect := sg.rect.Lerp(&targetRect, 0.05)
		sg.rect = &newRect
	}

	closestX, closestY := sg.rect.FindClosestPointOnEdge(sg.anchorX, sg.anchorY)
	sg.tickX = closestX
	sg.tickY = closestY
	s := sg.tickSprite.Bounds().Size()

	// This ended up being very illogical
	if sg.rect.Contains(sg.anchorX, sg.anchorY) {
		switch sg.rect.GetInteriorSection(sg.anchorX, sg.anchorY) {
		case maths.TOP:
			sg.tickRotation = maths.DirDown
			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+float64(s.X), sg.rect.Right()-1)
			sg.tickY += 1
		case maths.LEFT:
			sg.tickRotation = maths.DirRight
			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+1, sg.rect.Bottom()-float64(s.Y))
			sg.tickX += 1
		case maths.RIGHT:
			sg.tickRotation = maths.DirLeft
			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+float64(s.Y), sg.rect.Bottom()+1)
			sg.tickX -= 1
		case maths.BOTTOM:
			sg.tickRotation = maths.DirUp
			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+1, sg.rect.Right()-float64(s.X))
			sg.tickY -= 1
		}
	} else {
		switch sg.rect.GetExteriorSection(sg.anchorX, sg.anchorY) {
		case maths.TOP_LEFT, maths.TOP_MIDDLE, maths.TOP_RIGHT:
			sg.tickRotation = maths.DirDown
			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+float64(s.X), sg.rect.Right()-1)
			sg.tickY += 1
		case maths.MIDDLE_LEFT:
			sg.tickRotation = maths.DirRight
			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+1, sg.rect.Bottom()-float64(s.Y))
			sg.tickX += 1
		case maths.MIDDLE_RIGHT:
			sg.tickRotation = maths.DirLeft
			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+float64(s.Y), sg.rect.Bottom()+1)
			sg.tickX -= 1
		case maths.BOTTOM_LEFT, maths.BOTTOM_MIDDLE, maths.BOTTOM_RIGHT:
			sg.tickRotation = maths.DirUp
			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+1, sg.rect.Right()-float64(s.X))
			sg.tickY -= 1
		}
	}
}

func (sg *speechBubbleGraphic) Draw(ctx rendering.Ctx) {
	l := float32(sg.rect.Left())
	r := float32(sg.rect.Right())
	t := float32(sg.rect.Top())
	b := float32(sg.rect.Bottom())
	w := float32(sg.rect.Width())
	h := float32(sg.rect.Height())
	vector.StrokeLine(ctx.Dst, l+1, t, r-1, t, 1, borderColor, false)
	vector.StrokeLine(ctx.Dst, l, t+1, l, b-1, 1, borderColor, false)
	vector.StrokeLine(ctx.Dst, l+1, b, r-1, b, 1, borderColor, false)
	vector.StrokeLine(ctx.Dst, r, t+1, r, b-1, 1, borderColor, false)
	// This shouldn't really be necessary but it seems to fix some pixel perfect bugs
	vector.StrokeRect(ctx.Dst, l+1, t+1, w-2, h-2, 1, interiorColor, false)
	vector.FillRect(ctx.Dst, l+1, t+1, w-2, h-2, interiorColor, false)

	ebitenrenderutil.DrawAtRotated(sg.tickSprite, ctx.Dst, sg.tickX, sg.tickY, maths.DirToRadians(sg.tickRotation))
}

// Function that sets width / height based on text
func (sg *speechBubbleGraphic) FitText() {

}
