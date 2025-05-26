package ui

import (
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/libraries/node"
	"os"

	"gopkg.in/yaml.v3"
)

type Layer struct {
	Name          string    `yaml:"Name"`
	Root          node.Root `yaml:"Root"`
	confirmations map[string]bool
}

func (d *Layer) Update() {
	d.Root.Update(d.confirmations)
}

func (d *Layer) Draw() {
	d.Root.Draw(0, 0, rendering.GAME_WIDTH*rendering.PIXEL_SCALE, rendering.GAME_HEIGHT*rendering.PIXEL_SCALE)
}

func (d *Layer) GetConfirmed() map[string]bool {
	return d.confirmations
}

func (d *Layer) GetSubmitted() map[string]string {
	// chart := make(map[string]string)
	//
	//	for _, inputbox := range m.Inputboxes {
	//		chart[inputbox.Name] = inputbox.Read()
	//	}
	//
	// return charttion A
	return nil
}

func FromFile(path string) (*Layer, error) {
	layer := &Layer{
		confirmations: make(map[string]bool),
	}
	file, err := os.Open(path)
	if err != nil {
		return layer, err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(layer)
	if err != nil {
		return layer, err
	}

	return layer, nil
}
