package menu

import (
	"mask_of_the_tomb/internal/game/UI/textbox"
	"mask_of_the_tomb/internal/maths"
	"os"

	"gopkg.in/yaml.v3"
)

type Menu struct {
	Name        string                `yaml:"Name"`
	Textboxes   []*textbox.Textbox    `yaml:"Textboxes"`
	Selectables []*textbox.Selectable `yaml:"Selectables"`
	SelectorPos int
}

func (m *Menu) Update(dirInput int) {
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

func (m *Menu) Draw() {
	for _, textbox := range m.Textboxes {
		textbox.Draw()
	}
	for _, selectable := range m.Selectables {
		selectable.Draw()
	}
}

func (m *Menu) GetConfirmed() (chart map[string]bool) {
	chart = make(map[string]bool)
	for _, selectable := range m.Selectables {
		chart[selectable.Name] = selectable.GetConfirm()
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
	}

	for _, textbox := range menu.Textboxes {
		textbox.Color.LoadColorPair()
	}
	return menu, nil
}

// func newMenu(textboxes []*Textbox, selectables []*Selectable) *Menu {
// 	return &Menu{
// 		Textboxes:   textboxes,
// 		Selectables: selectables,
// 		selectorPos: 0,
// 	}
// }
