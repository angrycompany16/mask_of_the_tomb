package selectable

import (
	"mask_of_the_tomb/internal/backend/colors"
	"mask_of_the_tomb/internal/engine/actors/UI/textbox"
	"mask_of_the_tomb/internal/engine/commands"
)

type Selectable struct {
	*textbox.Textbox
	NormalColor   colors.ColorPair
	SelectedColor colors.ColorPair
	Selected      bool
}

func (b *Selectable) Update(cmd *commands.Commands) {
	b.Textbox.Update(cmd)
}

func (b *Selectable) SetSelected(suppressSound bool) {
	// if !b.selected && !suppresSound {
	// 	sound_v2.PlaySound("selectUI", "sfxMaster", 0)
	// }

	b.Selected = true
	b.Textbox.Color = b.SelectedColor
	// b.Textbox.Color = colors.ColorPair{
	// 	BrightColor: color.RGBA64{255, 255, 255, 255},
	// 	DarkColor:   color.RGBA64{255, 200, 255, 255},
	// }
	// colors.ColorPair{
	// 	BrightColor: color.RGBA{255, 255, 255, 255},
	// 	DarkColor:   color.RGBA{255, 200, 255, 255},
	// }
}

func (b *Selectable) SetDeselected() {
	b.Selected = false
	b.Textbox.Color = b.NormalColor
	// b.Textbox.Color = colors.ColorPair{
	// 	BrightColor: color.RGBA64{255, 0, 0, 255},
	// 	DarkColor:   color.RGBA64{150, 50, 50, 255},
	// }

	// colors.ColorPair{
	// 	BrightColor: color.RGBA{255, 0, 255, 255},
	// 	DarkColor:   color.RGBA{255, 0, 255, 255},
	// }
}

func NewSelectable(textbox *textbox.Textbox, normalColor colors.ColorPair, selectedColor colors.ColorPair) *Selectable {
	return &Selectable{
		Textbox:       textbox,
		NormalColor:   normalColor,
		SelectedColor: selectedColor,
		Selected:      false,
	}
}
