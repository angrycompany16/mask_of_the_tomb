package menu

import (
	"mask_of_the_tomb/internal/game/UI/search"
	"mask_of_the_tomb/internal/game/UI/textbox"
	"mask_of_the_tomb/internal/maths"
	"os"

	"gopkg.in/yaml.v3"
)

type Menu struct {
	Name        string                `yaml:"Name"`
	Textboxes   []*textbox.Textbox    `yaml:"Textboxes"`
	Selectables []*textbox.Selectable `yaml:"Selectables"`
	Inputboxes  []*textbox.Inputbox   `yaml:"Inputboxes"`
	FileSearch  *search.FileSearch    `yaml:"FileSearch"`
	SelectorPos int
}

// A problem: Right now FileSearch is a pointer, which means that it can be nil sometimes.
// However, in these cases we cannot update it. However, we don't want to check every time
// if the filesearch is nil as that's kinda just ugly
// Might have to figure out something smart

// TODO: Include select for input boxes
func (m *Menu) UpdateSelectables(dirInput int) {
	if len(m.Selectables) == 0 {
		return
	}

	m.SelectorPos += dirInput
	m.SelectorPos = maths.Mod(m.SelectorPos, len(m.Selectables))

	for i, selectable := range m.Selectables {
		if i == m.SelectorPos {
			selectable.SetSelected()
		} else {
			selectable.SetDeselected()
		}
		i++
	}
}

func (m *Menu) UpdateInputboxes() {
	for _, inputbox := range m.Inputboxes {
		inputbox.Update()
	}
}

func (m *Menu) Draw() {
	for _, textbox := range m.Textboxes {
		textbox.Draw()
	}
	for _, selectable := range m.Selectables {
		selectable.Draw()
	}
	for _, inputbox := range m.Inputboxes {
		inputbox.Draw()
	}
	if m.FileSearch != nil {
		m.FileSearch.Draw()
	}
}

func (m *Menu) GetConfirmed() map[string]bool {
	chart := make(map[string]bool)
	for _, selectable := range m.Selectables {
		chart[selectable.Name] = selectable.GetConfirmed()
	}
	return chart
}

func (m *Menu) GetSubmitted() map[string]string {
	chart := make(map[string]string)
	for _, inputbox := range m.Inputboxes {
		chart[inputbox.Name] = inputbox.Read()
	}
	return chart
}

func FromFile(path string) (*Menu, error) {
	menu := &Menu{}
	file, err := os.Open(path)
	if err != nil {
		return menu, err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(menu)
	if err != nil {
		return menu, err
	}
	menu.SelectorPos = 0

	for _, selectable := range menu.Selectables {
		selectable.NormalColor.LoadColorPair()
		selectable.HoverColor.LoadColorPair()
		selectable.Textbox.SetFont()
		selectable.Textbox.Color = selectable.NormalColor
	}

	for _, textbox := range menu.Textboxes {
		textbox.Color.LoadColorPair()
		textbox.SetFont()
	}

	for _, inputbox := range menu.Inputboxes {
		inputbox.Textbox.SetFont()
		inputbox.Textbox.Color.LoadColorPair()
	}

	if menu.FileSearch != nil {
		menu.FileSearch.Searchfield.Textbox.SetFont()
		menu.FileSearch.Searchfield.Textbox.Color.LoadColorPair()
		menu.FileSearch.ResultNormalColor.LoadColorPair()
		menu.FileSearch.ResultSelectedColor.LoadColorPair()
	}

	return menu, nil
}
