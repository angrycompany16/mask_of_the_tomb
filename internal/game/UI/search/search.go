package search

import (
	"mask_of_the_tomb/editor/fileio"
	"mask_of_the_tomb/internal/game/UI/colorpair"
	"mask_of_the_tomb/internal/game/UI/textbox"
	"mask_of_the_tomb/internal/maths"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type FileSearch struct {
	Name                string               `yaml:"Name"`
	Searchfield         *textbox.Inputbox    `yaml:"Searchfield"`
	ResultNormalColor   *colorpair.ColorPair `yaml:"NormalColor"`
	ResultSelectedColor *colorpair.ColorPair `yaml:"HoverColor"`
	Results             []*textbox.Selectable
	SelectorPos         int
}

func (f *FileSearch) Update(dirInput int) {
	f.Searchfield.Update()

	if len(f.Results) == 0 {
		return
	}

	f.SelectorPos += dirInput
	f.SelectorPos = maths.Mod(f.SelectorPos, len(f.Results))

	for i, selectable := range f.Results {
		if i == f.SelectorPos {
			selectable.SetSelected()
		} else {
			selectable.SetDeselected()
		}
		i++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		f.Searchfield.Textbox.Text = f.Results[f.SelectorPos].Textbox.Text
	}
}

func (f *FileSearch) Draw() {
	f.Searchfield.Draw()
	for _, result := range f.Results {
		result.Draw()
	}
}

func (f *FileSearch) UpdateSearchResults() {
	// TODO: Listen to event
	// WTF is this
	results := make([]string, 0)
	results = fileio.FindFiles(f.Searchfield.Textbox.Text, results)

	if len(results) == len(f.Results) {
		return
	}

	// fmt.Println("New results")
	f.Results = make([]*textbox.Selectable, 0)

	pos := f.Searchfield.Textbox.PosY + f.Searchfield.Textbox.FontSize
	for _, result := range results {
		resultTextBox := *f.Searchfield.Textbox
		resultTextBox.Text = result
		resultTextBox.PosY = pos

		f.Results = append(f.Results, &textbox.Selectable{
			Textbox:     &resultTextBox,
			NormalColor: *f.ResultNormalColor,
			HoverColor:  *f.ResultSelectedColor,
			Name:        result,
		})
		pos += resultTextBox.FontSize
	}
}
