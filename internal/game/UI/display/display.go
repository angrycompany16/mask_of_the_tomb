package display

import "mask_of_the_tomb/internal/game/UI/node"

// How to make the ui system more systematic

// Idea: Node hierarchy
// Nodes have position and size
// Textboxes, Input fields, buttons, sliders, search result lists, etc... are nodes
// with extra behaviour on top

// Thus each menu becomes a hierarchy of nodes. There should always be one node which
// is the root

// IDEA: Make interfaces such as selectable,

type Display struct {
	Name string    `yaml:"Name"`
	Root node.Node `yaml:"Root"`
}

// Node types:
// Selectable list
// Input field
// more

// type Menu struct {
// 	Name        string                `yaml:"Name"`
// 	Textboxes   []*textbox.Textbox    `yaml:"Textboxes"`
// 	Selectables []*textbox.Selectable `yaml:"Selectables"`
// 	Inputboxes  []*textbox.Inputbox   `yaml:"Inputboxes"`
// 	FileSearch  *search.FileSearch    `yaml:"FileSearch"`
// 	SelectorPos int
// }

// A problem: Right now FileSearch is a pointer, which means that it can be nil sometimes.
// However, in these cases we cannot update it. However, we don't want to check every time
// if the filesearch is nil as that's kinda just ugly
// Might have to figure out something smart

// TODO: Include select for input boxes
func (d *Display) UpdateSelectables(dirInput int) {
	// 	if len(m.Selectables) == 0 {
	// 		return
	// 	}

	// 	m.SelectorPos += dirInput
	// 	m.SelectorPos = maths.Mod(m.SelectorPos, len(m.Selectables))

	//	for i, selectable := range m.Selectables {
	//		if i == m.SelectorPos {
	//			selectable.SetSelected()
	//		} else {
	//			selectable.SetDeselected()
	//		}
	//		i++
	//	}
}

func (d *Display) UpdateInputboxes() {
	//	for _, inputbox := range m.Inputboxes {
	//		inputbox.Update()
	//	}
}

func (d *Display) Draw() {
	d.Root.Draw(0, 0)
}

// func (d *Display) Draw() {
// 	for _, textbox := range m.Textboxes {
// 		textbox.Draw()
// 	}
// 	for _, selectable := range m.Selectables {
// 		selectable.Draw()
// 	}
// 	for _, inputbox := range m.Inputboxes {
// 		inputbox.Draw()
// 	}
// 	if m.FileSearch != nil {
// 		m.FileSearch.Draw()
// 	}
// }

func (d *Display) GetConfirmed() map[string]bool {
	// chart := make(map[string]bool)
	//
	//	for _, selectable := range m.Selectables {
	//		chart[selectable.Name] = selectable.GetConfirmed()
	//	}
	//
	// return chart
	return nil
}

func (d *Display) GetSubmitted() map[string]string {
	// chart := make(map[string]string)
	//
	//	for _, inputbox := range m.Inputboxes {
	//		chart[inputbox.Name] = inputbox.Read()
	//	}
	//
	// return chart
	return nil
}

// // Want: Some way to group all these functions, which should take in a path and
// // return a pointer to the asset, pool them together and then load them
// // but alas the question is how

func FromFile(path string) (*Display, error) {
	// 	menu := &Menu{}
	// 	file, err := os.Open(path)
	// 	if err != nil {
	// 		return menu, err
	// 	}
	// 	defer file.Close()

	// 	err = yaml.NewDecoder(file).Decode(menu)
	// 	if err != nil {
	// 		return menu, err
	// 	}
	// 	menu.SelectorPos = 0

	// 	// TODO: Resolve SetFont and LoadColorPair
	// 	for _, selectable := range menu.Selectables {
	// 		selectable.NormalColor.LoadColorPair()
	// 		selectable.HoverColor.LoadColorPair()
	// 		selectable.Textbox.SetFont()
	// 		selectable.Textbox.Color = selectable.NormalColor
	// 	}

	// 	for _, textbox := range menu.Textboxes {
	// 		textbox.Color.LoadColorPair()
	// 		textbox.SetFont()
	// 	}

	// 	for _, inputbox := range menu.Inputboxes {
	// 		inputbox.Textbox.SetFont()
	// 		inputbox.Textbox.Color.LoadColorPair()
	// 	}

	// 	if menu.FileSearch != nil {
	// 		menu.FileSearch.Searchfield.Textbox.SetFont()
	// 		menu.FileSearch.Searchfield.Textbox.Color.LoadColorPair()
	// 		menu.FileSearch.ResultNormalColor.LoadColorPair()
	// 		menu.FileSearch.ResultSelectedColor.LoadColorPair()
	// 	}

	// return menu, nil
	return &Display{}, nil
}
