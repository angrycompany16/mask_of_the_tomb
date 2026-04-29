package cursor

import (
	"mask_of_the_tomb/internal/engine/actors/UI/uigraphic"
	"mask_of_the_tomb/internal/engine/commands"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO: BUG that is not too problematic:
// Cursor position upon launching the game somehow gets offset,
// which causes a small jump in the cursor movement as soon
// as the player starts to move the cursor...

// TODO: OnMove, OnHover, OnClick
type Cursor struct {
	*uigraphic.UIGraphic
}

func (c *Cursor) Init(cmd *commands.Commands) {
	c.UIGraphic.Init(cmd)
	mouseX, mouseY := ebiten.CursorPosition()
	c.UIGraphic.SetPos(float64(mouseX), float64(mouseY))
}

// TODO: Here's an annoying thing: To avoid flickering we need to set the
// mouse pos also in Init, which is kind of ugly... would be nice if
// we could somehow defer rendering to the end of the Update loop always...
func (c *Cursor) Update(cmd *commands.Commands) {
	c.UIGraphic.Update(cmd)
	mouseX, mouseY := ebiten.CursorPosition()
	c.UIGraphic.SetPos(float64(mouseX), float64(mouseY))
}

func NewCursor(uigraphic *uigraphic.UIGraphic) *Cursor {
	return &Cursor{
		UIGraphic: uigraphic,
	}
}
