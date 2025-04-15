package display

import (
	"mask_of_the_tomb/internal/game/UI/node"
	"mask_of_the_tomb/internal/game/core/rendering"
	"os"

	"gopkg.in/yaml.v3"
)

type Display struct {
	Name          string    `yaml:"Name"`
	Root          node.Root `yaml:"Root"`
	confirmations map[string]bool
}

func (d *Display) Update() {
	// pass confirmations through everywhere
	d.Root.Update(d.confirmations)
}

func (d *Display) Draw() {
	d.Root.Draw(0, 0, rendering.GameWidth*rendering.PixelScale, rendering.GameHeight*rendering.PixelScale)
}

func (d *Display) GetConfirmed() map[string]bool {
	return d.confirmations
}

func (d *Display) GetSubmitted() map[string]string {
	// chart := make(map[string]string)
	//
	//	for _, inputbox := range m.Inputboxes {
	//		chart[inputbox.Name] = inputbox.Read()
	//	}
	//
	// return charttion A
	return nil
}

func FromFile(path string) (*Display, error) {
	display := &Display{
		confirmations: make(map[string]bool),
	}
	file, err := os.Open(path)
	if err != nil {
		return display, err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(display)
	if err != nil {
		return display, err
	}

	return display, nil
}
