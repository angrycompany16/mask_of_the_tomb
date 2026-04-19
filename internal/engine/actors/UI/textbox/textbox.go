package textbox

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/colors"
	eventsv2 "mask_of_the_tomb/internal/backend/events_v2"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/container"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Textbox struct {
	*container.Container
	Text           string
	Color          colors.ColorPair
	srcPath        string
	FontRef        *assetloader.AssetRef[text.GoTextFaceSource]
	FontSize       float64
	LineSpacing    float64
	PrimaryAlign   text.Align
	SecondaryAlign text.Align
	ShadowX        float64
	ShadowY        float64
	Image          *ebiten.Image
	target         renderer.RenderTarget
	drawOrder      int
	OnResize       *eventsv2.EventBus
	pivotX, pivotY float64
}

func (t *Textbox) OnTreeAdd(node *engine.Node, cmd *commands.Commands) {
	t.Container.OnTreeAdd(node, cmd)
	t.FontRef = assetloader.StageAsset[text.GoTextFaceSource](
		cmd.AssetLoader,
		t.srcPath,
		assettypes.NewFontAsset(t.srcPath),
	)
}

func (t *Textbox) Init(cmd *commands.Commands) {
	t.Container.Init(cmd)
	t.OnResize = eventsv2.NewEventBus(t.Container.OnResize)
}

func (t *Textbox) Update(cmd *commands.Commands) {
	t.Container.Update(cmd)

	if data, raised := t.OnResize.Poll(); raised {
		newRect := data["Rect"].(maths.Rect)
		t.Image = ebiten.NewImage(int(newRect.Width), int(newRect.Height))
	}

	t.Image.Clear()
	opText := &text.DrawOptions{}
	opText.LayoutOptions.LineSpacing = t.LineSpacing
	opText.LayoutOptions.PrimaryAlign = t.PrimaryAlign
	opText.LayoutOptions.SecondaryAlign = t.SecondaryAlign
	opText.GeoM.Translate(t.ShadowX, t.ShadowY)

	if t.PrimaryAlign == text.AlignCenter {
		opText.GeoM.Translate(t.Rect.Width/2, t.Rect.Height/2)
	}

	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.Color.DarkColor)

	text.Draw(t.Image,
		t.Text,
		&text.GoTextFace{
			Source: t.FontRef.Value(),
			Size:   t.FontSize,
		}, opText)

	opText.ColorScale = ebiten.ColorScale{}
	opText.ColorScale.ScaleWithColor(t.Color.BrightColor)

	opText.GeoM.Translate(-t.ShadowX, -t.ShadowY)
	text.Draw(t.Image, t.Text, &text.GoTextFace{
		Source: t.FontRef.Value(),
		Size:   t.FontSize,
	}, opText)

	x, y := t.GetAbsPos()
	cmd.Renderer.Request(opgen.Pos(t.Image, x, y, t.pivotX, t.pivotY), t.Image, t.target, t.drawOrder)
}

func defaultTextbox(container *container.Container) *Textbox {
	return &Textbox{
		Container: container,
		Text:      "Lorem ipsum",
		Color: colors.ColorPair{
			BrightColor: color.RGBA{255, 255, 255, 255},
			DarkColor:   color.RGBA{35, 35, 35, 255},
		},
		FontSize:       64,
		LineSpacing:    4,
		PrimaryAlign:   text.AlignCenter,
		SecondaryAlign: text.AlignCenter,
		ShadowX:        4,
		ShadowY:        4,
		Image:          ebiten.NewImage(int(container.Rect.Width), int(container.Rect.Height)),
		target: renderer.RenderTarget{
			Type: renderer.SCREEN,
			Name: "ScreenUI",
		},
		drawOrder: 0,
		pivotX:    0,
		pivotY:    0,
	}
}

func NewTextBox(container *container.Container, fontPath string, options ...utils.Option[Textbox]) *Textbox {
	textbox := defaultTextbox(container)

	textbox.srcPath = fontPath

	for _, option := range options {
		option(textbox)
	}

	return textbox
}

func WithText(text string) utils.Option[Textbox] {
	return func(t *Textbox) {
		t.Text = text
	}
}

func WithPivot(x, y float64) utils.Option[Textbox] {
	return func(t *Textbox) {
		t.pivotX = x
		t.pivotY = y
	}
}
