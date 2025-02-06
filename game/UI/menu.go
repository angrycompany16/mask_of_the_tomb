package ui

import "mask_of_the_tomb/utils"

type menu struct {
	selectables []*selectable
	selectorPos int
}

func (m *menu) update(dirInput int, confirmInput bool) {
	if len(m.selectables) == 0 {
		return
	}

	m.selectorPos += dirInput
	m.selectorPos = utils.Mod(m.selectorPos, len(m.selectables))

	for i, selectable := range m.selectables {
		if i == m.selectorPos {
			selectable.setSelected()
		} else {
			selectable.setDeselected()
		}
	}
}

func (m *menu) draw() {
	for _, selectable := range m.selectables {
		selectable.draw()
	}
}

func newMenu(selectables []*selectable) *menu {
	return &menu{
		selectables: selectables,
		selectorPos: 0,
	}
}
