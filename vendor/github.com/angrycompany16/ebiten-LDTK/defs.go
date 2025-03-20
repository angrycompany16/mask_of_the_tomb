package ebitenLDTK

import (
	"fmt"
)

type Defs struct {
	Tilesets      []Tileset      `json:"tilesets"`
	Enums         []Enum         `json:"enums"`
	ExternalEnums []ExternalEnum `json:"externalEnums"`
	LevelFields   []LevelField   `json:"levelFields"`
}

type Enum struct {
	// TBA
}

type ExternalEnum struct {
	// TBA
}

type LevelField struct {
	// TBA
}

func (d *Defs) GetTilesetByUid(uid int) (Tileset, error) {
	for _, tileset := range d.Tilesets {
		if tileset.Uid == uid {
			return tileset, nil
		}
	}
	return Tileset{}, fmt.Errorf("tileset with uid [%d] was not found", uid)
}
