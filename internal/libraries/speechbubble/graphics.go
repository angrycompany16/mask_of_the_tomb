package speechbubble

import (
	"image/color"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/maths"
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/core/threads"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const interactPadding = 4 * rendering.PIXEL_SCALE

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
		sg.rect.Left()-interactPadding,
		sg.rect.Top()-interactPadding,
		sg.rect.Width()+interactPadding*2,
		sg.rect.Height()+interactPadding*2,
	)

	var targetRect maths.Rect
	targetRect.SetSize(sg.targetWidth, sg.targetHeight)
	if !paddedRect.Contains(sg.anchorX, sg.anchorY) {
		newPaddedRect := paddedRect.Reach(sg.anchorX, sg.anchorY)
		targetRect.SetPos(newPaddedRect.Left()+interactPadding, newPaddedRect.Top()+interactPadding)
	} else if sg.rect.Contains(sg.anchorX, sg.anchorY) {
		reachRect := sg.rect.Reach(sg.anchorX, sg.anchorY)
		targetRect.SetPos(reachRect.Left(), reachRect.Top())
	} else {
		targetRect.SetPos(sg.rect.Left(), sg.rect.Top())
	}
	newRect := sg.rect.Lerp(&targetRect, 0.05)
	sg.rect = &newRect

	closestX, closestY := sg.rect.FindClosestPointOnEdge(sg.anchorX, sg.anchorY)
	sg.tickX = closestX
	sg.tickY = closestY
	s := sg.tickSprite.Bounds().Size()

	// This ended up being very illogical
	if sg.rect.Contains(sg.anchorX, sg.anchorY) {
		switch sg.rect.GetInteriorSection(sg.anchorX, sg.anchorY) {
		case maths.TOP:
			sg.tickRotation = maths.DirDown
			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+float64(s.X)*rendering.PIXEL_SCALE, sg.rect.Right()-rendering.PIXEL_SCALE)
			sg.tickY += 2
		case maths.LEFT:
			sg.tickRotation = maths.DirRight
			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+rendering.PIXEL_SCALE, sg.rect.Bottom()-float64(s.Y)*rendering.PIXEL_SCALE)
			sg.tickX += 2
		case maths.RIGHT:
			sg.tickRotation = maths.DirLeft
			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+float64(s.Y)*rendering.PIXEL_SCALE, sg.rect.Bottom()+rendering.PIXEL_SCALE)
			sg.tickX -= 2
		case maths.BOTTOM:
			sg.tickRotation = maths.DirUp
			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+rendering.PIXEL_SCALE, sg.rect.Right()-float64(s.X)*rendering.PIXEL_SCALE)
			sg.tickY -= 2
		}
	} else {
		switch sg.rect.GetExteriorSection(sg.anchorX, sg.anchorY) {
		case maths.TOP_LEFT, maths.TOP_MIDDLE, maths.TOP_RIGHT:
			sg.tickRotation = maths.DirDown
			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+float64(s.X)*rendering.PIXEL_SCALE, sg.rect.Right()-rendering.PIXEL_SCALE)
			sg.tickY += 2
		case maths.MIDDLE_LEFT:
			sg.tickRotation = maths.DirRight
			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+rendering.PIXEL_SCALE, sg.rect.Bottom()-float64(s.Y)*rendering.PIXEL_SCALE)
			sg.tickX += 2
		case maths.MIDDLE_RIGHT:
			sg.tickRotation = maths.DirLeft
			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+float64(s.Y)*rendering.PIXEL_SCALE, sg.rect.Bottom()+rendering.PIXEL_SCALE)
			sg.tickX -= 2
		case maths.BOTTOM_LEFT, maths.BOTTOM_MIDDLE, maths.BOTTOM_RIGHT:
			sg.tickRotation = maths.DirUp
			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+rendering.PIXEL_SCALE, sg.rect.Right()-float64(s.X)*rendering.PIXEL_SCALE)
			sg.tickY -= 2
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
	vector.StrokeLine(ctx.Dst, l+rendering.PIXEL_SCALE/2, t, r-rendering.PIXEL_SCALE/2, t, rendering.PIXEL_SCALE, borderColor, false)
	vector.StrokeLine(ctx.Dst, l, t+rendering.PIXEL_SCALE/2, l, b-rendering.PIXEL_SCALE/2, rendering.PIXEL_SCALE, borderColor, false)
	vector.StrokeLine(ctx.Dst, l+rendering.PIXEL_SCALE/2, b, r-rendering.PIXEL_SCALE/2, b, rendering.PIXEL_SCALE, borderColor, false)
	vector.StrokeLine(ctx.Dst, r, t+rendering.PIXEL_SCALE/2, r, b-rendering.PIXEL_SCALE/2, rendering.PIXEL_SCALE, borderColor, false)

	vector.FillRect(ctx.Dst, l+rendering.PIXEL_SCALE/2, t+rendering.PIXEL_SCALE/2, w-rendering.PIXEL_SCALE, h-rendering.PIXEL_SCALE, interiorColor, false)

	ebitenrenderutil.DrawAtRotatedScaled(sg.tickSprite, ctx.Dst, sg.tickX, sg.tickY, maths.DirToRadians(sg.tickRotation), rendering.PIXEL_SCALE, rendering.PIXEL_SCALE)
}

// TODO: Figure out text breaking and box sizing
// Then figure out how to play audio along with
// the text.
// Then profit
// Then code the whole NPC system
// Then we possible also need some more customization

const revealPeriod = 0.08

type speechBubbleText struct {
	x, y               float64
	paddingX, paddingY float64
	text               string
	revealIndex        int
	font               *text.GoTextFaceSource
	fontSize           float64
	revealTicker       *time.Ticker
	relLineSpacing     float64
}

func (st *speechBubbleText) Update() (bool, byte) {
	if _, raised := threads.Poll(st.revealTicker.C); raised {
		st.revealIndex++
		isNew := st.revealIndex < len(st.text)+1
		st.revealIndex = maths.Clamp(st.revealIndex, 0, len(st.text))
		if st.revealIndex > 0 {
			return isNew, st.text[st.revealIndex-1]
		}
	}
	return false, 0
}

func (st *speechBubbleText) GetRevealed() string {
	return st.text[0:st.revealIndex]
}

func (st *speechBubbleText) Size() (float64, float64) {
	return text.Measure(st.text, &text.GoTextFace{
		Source: st.font,
		Size:   st.fontSize,
	}, st.fontSize*st.relLineSpacing)
}

func (st *speechBubbleText) Draw(ctx rendering.Ctx) {
	opText := &text.DrawOptions{}
	// Just need to figure out how we're supposed to find out the position
	// of this thing
	opText.GeoM.Translate(st.x, st.y)
	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.SetR(float32(borderColor.R))
	opText.ColorScale.SetG(float32(borderColor.G))
	opText.ColorScale.SetB(float32(borderColor.B))
	opText.LineSpacing = st.fontSize * st.relLineSpacing

	text.Draw(ctx.Dst, st.GetRevealed(), &text.GoTextFace{
		Source: st.font,
		Size:   st.fontSize,
	}, opText)
}
