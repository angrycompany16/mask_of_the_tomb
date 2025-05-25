package ui

import (
	"mask_of_the_tomb/internal/core/rendering"
	"os"

	"gopkg.in/yaml.v3"
)

type Display struct {
	Name          string `yaml:"Name"`
	Root          Root   `yaml:"Root"`
	confirmations map[string]bool
}

func (d *Display) Update() {
	d.Root.Update(d.confirmations)
}

func (d *Display) Draw() {
	d.Root.Draw(0, 0, rendering.GAME_WIDTH*rendering.PIXEL_SCALE, rendering.GAME_HEIGHT*rendering.PIXEL_SCALE)
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
