package speechbubble

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine/actors/sprite"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/actors/vectorgraphic"
	"mask_of_the_tomb/internal/engine/commands"

	"github.com/hajimehoshi/ebiten/v2"
)

// How we might organize this: Speechbubble itself is a vector graphic containing the dynamically resizing box
// Then inside that we have the text, interacting with the player input and all that
// Then we also have the tick mark pointing to where the character / speaker is

// So a possible way to organize this is as follows:
// Parent ((inherits from) Vectorgraphic)
// |> Text ((inherits from) Textbox from UI?) I think that might work, which is pretty nice
// |>	Tick (inherits from sprite, has custom movement) (It's possible that we can get away with just making this a normal sprite and controlling it
//				  using the SpeechBubble)
// |> Vocalizer (inherits from soundplayer, but has five sounds I guess: One for each vowel)

const BOUNDARY_OFFSET = 16.0 // This is really annoying but if we don't use it the vector graphic doesn't render properly :(

type SpeechBubble struct {
	*vectorgraphic.VectorGraphic
	rect                      *maths.Rect
	targetWidth, targetHeight float64                  // Determine the target width and height for the SpeechBubble
	anchor                    *transform2D.Transform2D // A transform representing the position of the SpeechBubble
	interactPadding           float64                  // Determines the padding between the anchor position and the textbox
	borderColor               *color.RGBA
	fillColor                 *color.RGBA
	tickSprite                *sprite.Sprite // Child reference
	// textbox                   *textbox.Textbox // For now
}

func (s *SpeechBubble) Init(cmd *commands.Commands) {
	s.VectorGraphic.Init(cmd)
	pixelScale := cmd.Renderer.GetPixelScale()
	s.VectorGraphic.DrawFunc = func(img *ebiten.Image) {
		w := s.rect.Width * pixelScale
		h := s.rect.Height * pixelScale

		vector64.StrokeLine(img, BOUNDARY_OFFSET+pixelScale/2, BOUNDARY_OFFSET, BOUNDARY_OFFSET+w-pixelScale/2, BOUNDARY_OFFSET, pixelScale, s.borderColor, false)
		vector64.StrokeLine(img, BOUNDARY_OFFSET, BOUNDARY_OFFSET+pixelScale/2, BOUNDARY_OFFSET, BOUNDARY_OFFSET+h-pixelScale/2, pixelScale, s.borderColor, false)
		vector64.StrokeLine(img, BOUNDARY_OFFSET+pixelScale/2, BOUNDARY_OFFSET+h, BOUNDARY_OFFSET+w-pixelScale/2, BOUNDARY_OFFSET+h, pixelScale, s.borderColor, false)
		vector64.StrokeLine(img, BOUNDARY_OFFSET+w, BOUNDARY_OFFSET+pixelScale/2, BOUNDARY_OFFSET+w, BOUNDARY_OFFSET+h-pixelScale/2, pixelScale, s.borderColor, false)

		vector64.FillRect(img, BOUNDARY_OFFSET+pixelScale/2, BOUNDARY_OFFSET+pixelScale/2, w-pixelScale, h-pixelScale, s.fillColor, false)
	}
}

func (s *SpeechBubble) Update(cmd *commands.Commands) {
	s.VectorGraphic.Update(cmd)

	pixelScale := cmd.Renderer.GetPixelScale()

	paddedRect := maths.NewRect(
		s.rect.Left()-s.interactPadding,
		s.rect.Top()-s.interactPadding,
		s.rect.Width+s.interactPadding*2,
		s.rect.Height+s.interactPadding*2,
	)

	anchorX, anchorY := s.anchor.GetPos(false)
	var targetRect maths.Rect
	targetRect.SetSize(s.targetWidth, s.targetHeight)
	if !paddedRect.Contains(anchorX, anchorY) {
		newPaddedRect := paddedRect.Reach(anchorX, anchorY)
		targetRect.SetPos(newPaddedRect.Left()+s.interactPadding, newPaddedRect.Top()+s.interactPadding)
	} else if s.rect.Contains(anchorX, anchorY) {
		reachRect := s.rect.Reach(anchorX, anchorY)
		targetRect.SetPos(reachRect.Left(), reachRect.Top())
	} else {
		targetRect.SetPos(s.rect.Left(), s.rect.Top())
	}
	newRect := s.rect.Lerp(&targetRect, 0.05)
	s.rect = &newRect
	s.Transform2D.SetPos(newRect.X*pixelScale-BOUNDARY_OFFSET/pixelScale, newRect.Y*pixelScale-BOUNDARY_OFFSET/pixelScale)

	closestX, closestY := s.rect.FindClosestPointOnEdge(anchorX, anchorY)

	tickX, tickY := closestX, closestY
	tickDirection := maths.DirDown
	tickW, tickH := s.tickSprite.GetSize()

	//	// This ended up being very illogical
	//	if s.rect.Contains(anchorX, anchorY) {
	//		switch s.rect.GetInteriorSection(anchorX, anchorY) {
	//		case maths.TOP:
	//			tickDirection = maths.DirDown
	//			tickX = maths.Clamp(tickX, tickW*pixelScale, s.rect.Width-pixelScale)
	//			tickY += 2
	//		case maths.LEFT:
	//			tickDirection = maths.DirRight
	//			tickY = maths.Clamp(tickY, pixelScale, s.rect.Height-tickH*pixelScale)
	//			tickX += 2
	//		case maths.RIGHT:
	//			tickDirection = maths.DirLeft
	//			tickY = maths.Clamp(tickY, tickH*pixelScale, s.rect.Height+pixelScale)
	//			tickX -= 2
	//		case maths.BOTTOM:
	//			tickDirection = maths.DirUp
	//			tickX = maths.Clamp(tickX, pixelScale, s.rect.Width-tickW*pixelScale)
	//			tickY -= 2
	//		}
	//	} else {
	switch s.rect.GetExteriorSection(anchorX, anchorY) {
	case maths.TOP_LEFT, maths.TOP_MIDDLE, maths.TOP_RIGHT:
		tickDirection = maths.DirDown
		tickX = maths.Clamp(tickX, s.rect.Left()+tickW, s.rect.Right())
		//tickY += 2
	case maths.MIDDLE_LEFT:
		tickDirection = maths.DirRight
		tickY = maths.Clamp(tickY, s.rect.Top()+tickH, s.rect.Bottom())
		//tickX += 2
	case maths.MIDDLE_RIGHT:
		tickDirection = maths.DirLeft
		tickY = maths.Clamp(tickY, s.rect.Top()+tickH, s.rect.Bottom())
		tickX += tickH
	case maths.BOTTOM_LEFT, maths.BOTTOM_MIDDLE, maths.BOTTOM_RIGHT:
		tickDirection = maths.DirUp
		tickX = maths.Clamp(tickX, s.rect.Left()+tickW, s.rect.Right())
		tickY += tickW
	}
	//	}

	s.tickSprite.SetPos(tickX*pixelScale, tickY*pixelScale)
	s.tickSprite.SetAngle(maths.DirToRadians(tickDirection))
}

// Note to self: Set interactPadding to 4 * pixel scale
func NewSpeechBubble(anchor *transform2D.Transform2D, tickSprite *sprite.Sprite, width, height, interactPadding float64) *SpeechBubble {
	anchorX, anchorY := anchor.GetPos(false)
	return &SpeechBubble{
		VectorGraphic: vectorgraphic.NewVectorGraphic(
			vectorgraphic.WithTarget(renderer.ScreenTarget("WorldUI")),
			vectorgraphic.WithImage(600, 400),
		),
		anchor:          anchor,
		tickSprite:      tickSprite,
		rect:            maths.NewRect(anchorX, anchorY, width, height),
		targetWidth:     width,
		targetHeight:    height,
		interactPadding: interactPadding,
		fillColor:       &color.RGBA{21, 10, 31, 255},
		borderColor:     &color.RGBA{255, 253, 240, 255},
	}
}

// Seems like this part mostly handles the *logic*, such as interfacing with key presses and stuff
//
//func (sb *SpeechBubble) Update() {
//	sb.graphic.Update()
//	newChar, c := sb.textDisplay.Update()
//	if newChar {
//		Vocalize(string(c))
//	}
//
//	sb.textDisplay.x = sb.graphic.rect.Left() + sb.textDisplay.paddingX
//	sb.textDisplay.y = sb.graphic.rect.Top() + sb.textDisplay.paddingY
//
//	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
//		sb.LoadActiveLine()
//		w, h := sb.textDisplay.Size()
//		sb.graphic.targetWidth = w + 2*sb.textDisplay.paddingX
//		sb.graphic.targetHeight = h + 2*sb.textDisplay.paddingY
//	}
//}
//
//func (sb *SpeechBubble) LoadActiveLine() {
//	if sb.textDisplay.revealIndex != len(sb.textDisplay.text) {
//		sb.textDisplay.revealIndex = len(sb.textDisplay.text) - 1
//		return
//	} else if sb.currentLine == len(sb.lines) {
//		sb.textDisplay.revealIndex = 0
//		sb.textDisplay.text = ""
//		return
//	}
//	sb.textDisplay.text = sb.lines[sb.currentLine]
//	sb.currentLine++
//	sb.textDisplay.revealIndex = 0
//	// Expand speech bubble graphic if needed
//}
//
//func (sb *SpeechBubble) Draw(ctx rendering.Ctx) {
//	if !sb.hidden {
//		sb.graphic.Draw(ctx)
//		sb.textDisplay.Draw(ctx)
//	}
//}
//
//func (sb *SpeechBubble) SetAnchor(x, y float64) {
//	sb.graphic.anchorX = x
//	sb.graphic.anchorY = y
//}
//
//// Very simple for now - These can be expanded later
//func (sb *SpeechBubble) Hide() {
//	sb.hidden = true
//}
//
//func (sb *SpeechBubble) Reveal() {
//	sb.hidden = false
//}
//
//func (sb *SpeechBubble) ToggleVisibility() {
//	sb.hidden = !sb.hidden
//}
//
//func (sb *SpeechBubble) SetLines(lines []string) {
//	sb.textDisplay.revealIndex = 0
//	sb.currentLine = 0
//	sb.lines = lines
//}
//
//func NewSpeechBubble(anchorX, anchorY, width, height float64, hidden bool) *SpeechBubble {
//	speechBubbleTick := errs.Must(assettypes.GetImageAsset("textBoxTickSprite"))
//	newSpeechBubble := SpeechBubble{}
//
//	newSpeechBubble.graphic = &speechBubbleGraphic{
//		rect:         maths.NewRect(anchorX, anchorY, width*rendering.PIXEL_SCALE, height*rendering.PIXEL_SCALE),
//		anchorX:      anchorX,
//		anchorY:      anchorY,
//		tickX:        anchorX,
//		tickY:        anchorY,
//		targetWidth:  width * rendering.PIXEL_SCALE,
//		targetHeight: height * rendering.PIXEL_SCALE,
//		tickSprite:   speechBubbleTick,
//	}
//	newSpeechBubble.textDisplay = &speechBubbleText{
//		x:        anchorX,
//		y:        anchorY,
//		paddingX: 32,
//		paddingY: 32,
//		font:     assetloader.GetFont("JSE_AmigaAMOS"),
//		fontSize: 24,
//
//		text:           "",
//		revealIndex:    0,
//		revealTicker:   time.NewTicker(time.Duration(revealPeriod * float64(time.Second))),
//		relLineSpacing: 1.2,
//	}
//	// newSpeechBubble.vocalizer = newVocalizer()
//	newSpeechBubble.lines = make([]string, 0)
//	newSpeechBubble.hidden = hidden
//	return &newSpeechBubble
//}
//
//const interactPadding = 4 * rendering.PIXEL_SCALE
//
//var interiorColor = color.RGBA{21, 10, 31, 255}
//var borderColor = color.RGBA{255, 253, 240, 255}
//
//type speechBubbleGraphic struct {
//	rect                      *maths.Rect
//	tickX, tickY              float64
//	anchorX, anchorY          float64
//	targetWidth, targetHeight float64
//	tickSprite                *ebiten.Image
//	tickRotation              maths.Direction
//}
//
//func (sg *speechBubbleGraphic) Update() {
//	paddedRect := maths.NewRect(
//		sg.rect.Left()-interactPadding,
//		sg.rect.Top()-interactPadding,
//		sg.rect.Width()+interactPadding*2,
//		sg.rect.Height()+interactPadding*2,
//	)
//
//	var targetRect maths.Rect
//	targetRect.SetSize(sg.targetWidth, sg.targetHeight)
//	if !paddedRect.Contains(sg.anchorX, sg.anchorY) {
//		newPaddedRect := paddedRect.Reach(sg.anchorX, sg.anchorY)
//		targetRect.SetPos(newPaddedRect.Left()+interactPadding, newPaddedRect.Top()+interactPadding)
//	} else if sg.rect.Contains(sg.anchorX, sg.anchorY) {
//		reachRect := sg.rect.Reach(sg.anchorX, sg.anchorY)
//		targetRect.SetPos(reachRect.Left(), reachRect.Top())
//	} else {
//		targetRect.SetPos(sg.rect.Left(), sg.rect.Top())
//	}
//	newRect := sg.rect.Lerp(&targetRect, 0.05)
//	sg.rect = &newRect
//
//	closestX, closestY := sg.rect.FindClosestPointOnEdge(sg.anchorX, sg.anchorY)
//	sg.tickX = closestX
//	sg.tickY = closestY
//	s := sg.tickSprite.Bounds().Size()
//
//	// This ended up being very illogical
//	if sg.rect.Contains(sg.anchorX, sg.anchorY) {
//		switch sg.rect.GetInteriorSection(sg.anchorX, sg.anchorY) {
//		case maths.TOP:
//			sg.tickRotation = maths.DirDown
//			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+float64(s.X)*rendering.PIXEL_SCALE, sg.rect.Right()-rendering.PIXEL_SCALE)
//			sg.tickY += 2
//		case maths.LEFT:
//			sg.tickRotation = maths.DirRight
//			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+rendering.PIXEL_SCALE, sg.rect.Bottom()-float64(s.Y)*rendering.PIXEL_SCALE)
//			sg.tickX += 2
//		case maths.RIGHT:
//			sg.tickRotation = maths.DirLeft
//			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+float64(s.Y)*rendering.PIXEL_SCALE, sg.rect.Bottom()+rendering.PIXEL_SCALE)
//			sg.tickX -= 2
//		case maths.BOTTOM:
//			sg.tickRotation = maths.DirUp
//			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+rendering.PIXEL_SCALE, sg.rect.Right()-float64(s.X)*rendering.PIXEL_SCALE)
//			sg.tickY -= 2
//		}
//	} else {
//		switch sg.rect.GetExteriorSection(sg.anchorX, sg.anchorY) {
//		case maths.TOP_LEFT, maths.TOP_MIDDLE, maths.TOP_RIGHT:
//			sg.tickRotation = maths.DirDown
//			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+float64(s.X)*rendering.PIXEL_SCALE, sg.rect.Right()-rendering.PIXEL_SCALE)
//			sg.tickY += 2
//		case maths.MIDDLE_LEFT:
//			sg.tickRotation = maths.DirRight
//			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+rendering.PIXEL_SCALE, sg.rect.Bottom()-float64(s.Y)*rendering.PIXEL_SCALE)
//			sg.tickX += 2
//		case maths.MIDDLE_RIGHT:
//			sg.tickRotation = maths.DirLeft
//			sg.tickY = maths.Clamp(sg.tickY, sg.rect.Top()+float64(s.Y)*rendering.PIXEL_SCALE, sg.rect.Bottom()+rendering.PIXEL_SCALE)
//			sg.tickX -= 2
//		case maths.BOTTOM_LEFT, maths.BOTTOM_MIDDLE, maths.BOTTOM_RIGHT:
//			sg.tickRotation = maths.DirUp
//			sg.tickX = maths.Clamp(sg.tickX, sg.rect.Left()+rendering.PIXEL_SCALE, sg.rect.Right()-float64(s.X)*rendering.PIXEL_SCALE)
//			sg.tickY -= 2
//		}
//	}
//}
//
//func (sg *speechBubbleGraphic) Draw(ctx rendering.Ctx) {
//	l := float32(sg.rect.Left())
//	r := float32(sg.rect.Right())
//	t := float32(sg.rect.Top())
//	b := float32(sg.rect.Bottom())
//	w := float32(sg.rect.Width())
//	h := float32(sg.rect.Height())
//	vector.StrokeLine(ctx.Dst, l+rendering.PIXEL_SCALE/2, t, r-rendering.PIXEL_SCALE/2, t, rendering.PIXEL_SCALE, borderColor, false)
//	vector.StrokeLine(ctx.Dst, l, t+rendering.PIXEL_SCALE/2, l, b-rendering.PIXEL_SCALE/2, rendering.PIXEL_SCALE, borderColor, false)
//	vector.StrokeLine(ctx.Dst, l+rendering.PIXEL_SCALE/2, b, r-rendering.PIXEL_SCALE/2, b, rendering.PIXEL_SCALE, borderColor, false)
//	vector.StrokeLine(ctx.Dst, r, t+rendering.PIXEL_SCALE/2, r, b-rendering.PIXEL_SCALE/2, rendering.PIXEL_SCALE, borderColor, false)
//
//	vector.FillRect(ctx.Dst, l+rendering.PIXEL_SCALE/2, t+rendering.PIXEL_SCALE/2, w-rendering.PIXEL_SCALE, h-rendering.PIXEL_SCALE, interiorColor, false)
//
//	ebitenrenderutil.DrawAtRotatedScaled(sg.tickSprite, ctx.Dst, sg.tickX, sg.tickY, maths.DirToRadians(sg.tickRotation), rendering.PIXEL_SCALE, rendering.PIXEL_SCALE)
//}
//
//// TODO: Figure out a fade-in/fade-out animation
//
//const revealPeriod = 0.08
//
//type speechBubbleText struct {
//	x, y               float64
//	paddingX, paddingY float64
//	text               string
//	revealIndex        int
//	font               *text.GoTextFaceSource
//	fontSize           float64
//	revealTicker       *time.Ticker
//	relLineSpacing     float64
//}
//
//func (st *speechBubbleText) Update() (bool, byte) {
//	if _, raised := threads.Poll(st.revealTicker.C); raised {
//		st.revealIndex++
//		isNew := st.revealIndex < len(st.text)+1
//		st.revealIndex = maths.Clamp(st.revealIndex, 0, len(st.text))
//		if st.revealIndex > 0 {
//			return isNew, st.text[st.revealIndex-1]
//		}
//	}
//	return false, 0
//}
//
//func (st *speechBubbleText) GetRevealed() string {
//	return st.text[0:st.revealIndex]
//}
//
//func (st *speechBubbleText) Size() (float64, float64) {
//	return text.Measure(st.text, &text.GoTextFace{
//		Source: st.font,
//		Size:   st.fontSize,
//	}, st.fontSize*st.relLineSpacing)
//}
//
//func (st *speechBubbleText) Draw(ctx rendering.Ctx) {
//	opText := &text.DrawOptions{}
//	// Just need to figure out how we're supposed to find out the position
//	// of this thing
//	opText.GeoM.Translate(st.x, st.y)
//	opText.ColorScale = ebiten.ColorScale{}
//	opText.ColorScale.SetR(float32(borderColor.R))
//	opText.ColorScale.SetG(float32(borderColor.G))
//	opText.ColorScale.SetB(float32(borderColor.B))
//	opText.LineSpacing = st.fontSize * st.relLineSpacing
//
//	text.Draw(ctx.Dst, st.GetRevealed(), &text.GoTextFace{
//		Source: st.font,
//		Size:   st.fontSize,
//	}, opText)
//}
