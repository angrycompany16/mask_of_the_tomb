package search

import (
	"mask_of_the_tomb/internal/game/UI/textbox"
	"mask_of_the_tomb/internal/maths"
)

type FileSearch struct {
	Name        string              `yaml:"Name"`
	Searchfield []*textbox.Inputbox `yaml:"Searchfield"`
	Results     []*textbox.Selectable
	SelectorPos int
}

func (f *FileSearch) UpdateSelectables(dirInput int) {
	if len(f.Selectables) == 0 {
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

func UpdateSearchResults() {

}
