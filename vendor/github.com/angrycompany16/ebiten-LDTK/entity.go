package ebitenLDTK

import "fmt"

type Entity struct {
	Name   string     `json:"__identifier"`
	Tile   TileVisual `json:"__tile"`
	Iid    string     `json:"iid"`
	Width  float64    `json:"width"`
	Height float64    `json:"height"`
	DefUid int        `json:"defUid"`
	Px     []float64  `json:"px"`
	Fields []Field    `json:"fieldInstances"`
}

type RenderMode string

const (
	RenderModeTile      = "Tile"
	RenderModeRectangle = "Rectangle"
)

type TileVisual struct {
	TilesetUid int     `json:"tilesetUid"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	W          float64 `json:"w"`
	H          float64 `json:"h"`
}

func (e *Entity) GetFieldByName(name string) (Field, error) {
	for _, field := range e.Fields {
		if field.Name == name {
			return field, nil
		}
	}
	return Field{}, fmt.Errorf("field with name [%s] was not found", name)
}
