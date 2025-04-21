package node

import (
	"mask_of_the_tomb/editor/fileio"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type FileSearch struct {
	NodeData `yaml:",inline"`
	Name     string `yaml:"Name"`
}

func (f *FileSearch) Update(confirmations map[string]bool) {
	// Update the search field
	searchField := f.Children[0].Node.(*InputField)
	searchField.Update(confirmations)
	confirmations[f.Name] = inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if !searchField.selected {
		searchField.SetSelected()
	}

	// Update the result list
	selectList := f.Children[1].Node.(*SelectList)
	selectList.Update(confirmations)
	f.UpdateSearchResults()
}

func (f *FileSearch) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(f.Width, f.Height, parentWidth, parentHeight)
	f.DrawChildren(offsetX+f.PosX, offsetY+f.PosY, w, h)
}

// TODO: Only update on key stroke
func (f *FileSearch) UpdateSearchResults() {
	results := make([]string, 0)
	searchField := f.Children[0].Node.(*InputField)
	results = fileio.FindFiles(searchField.Text, results)

	selectList := f.Children[1].Node.(*SelectList)

	if len(results) == len(selectList.Children) {
		return
	}

	selectList.Children = make([]NodeContainer, 0)

	offset := 0.0
	for _, result := range results {
		resultTextBox := searchField.Button
		resultTextBox.Text = result
		resultTextBox.PosY = offset

		selectList.AddChild(NodeContainer{&resultTextBox})
		offset += resultTextBox.FontSize
	}
}

func (f *FileSearch) Reset() {
	f.ResetChildren()
}
