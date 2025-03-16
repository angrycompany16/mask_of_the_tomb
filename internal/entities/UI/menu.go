package ui

import (
	"mask_of_the_tomb/internal/maths"
)

type menu struct {
	textboxes   []*textbox
	selectables []*selectable
	selectorPos int
}

func (m *menu) update(dirInput int) {
	if len(m.selectables) == 0 {
		return
	}

	m.selectorPos += dirInput
	m.selectorPos = maths.Mod(m.selectorPos, len(m.selectables))

	// Important note - iterating through map is random
	for i, selectable := range m.selectables {
		if i == m.selectorPos {
			selectable.setSelected()
		} else {
			selectable.setDeselected()
		}
		i++
	}
}

func (m *menu) draw() {
	for _, textbox := range m.textboxes {
		textbox.draw()
	}
	for _, selectable := range m.selectables {
		selectable.draw()
	}
}

func (m *menu) getConfirmed() (chart map[string]bool) {
	chart = make(map[string]bool)
	for _, selectable := range m.selectables {
		chart[selectable.name] = selectable.getConfirm()
	}
	return chart
}

func newMenu(textboxes []*textbox, selectables []*selectable) *menu {
	return &menu{
		textboxes:   textboxes,
		selectables: selectables,
		selectorPos: 0,
	}
}
