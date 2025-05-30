package ui

import (
	"mask_of_the_tomb/internal/core/rendering"
	"mask_of_the_tomb/internal/libraries/node"
	"os"

	"gopkg.in/yaml.v3"
)

type Layer struct {
	Name          string                      `yaml:"Name"`
	Root          node.Root                   `yaml:"Root"`
	confirmations map[string]node.ConfirmInfo // TODO: Maybe include some extra info with the confirmations map
}

func (d *Layer) Update() {
	d.confirmations = make(map[string]node.ConfirmInfo)
	d.Root.Update(d.confirmations)
}

func (d *Layer) Draw() {
	d.Root.Draw(0, 0, rendering.GAME_WIDTH*rendering.PIXEL_SCALE, rendering.GAME_HEIGHT*rendering.PIXEL_SCALE)
}

func (d *Layer) GetConfirmed() map[string]node.ConfirmInfo {
	return d.confirmations
}

// func (d *Layer) SetValues(values map[string]node.OverWriteInfo) {
// 	d.Root.
// }

func (d *Layer) GetSubmitted() map[string]string {
	return nil
}

func FromFile(path string) (*Layer, error) {
	layer := &Layer{
		confirmations: make(map[string]node.ConfirmInfo),
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
