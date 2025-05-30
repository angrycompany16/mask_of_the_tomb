package assettypes

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"path/filepath"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type LDTKAsset struct {
	path  string
	World ebitenLDTK.World
}

func (a *LDTKAsset) Load() error {
	world, err := ebitenLDTK.LoadWorld(a.path)
	if err != nil {
		return err
	}

	a.World = world
	LDTKpath := filepath.Clean(filepath.Join(a.path, ".."))
	for i := 0; i < len(a.World.Defs.Tilesets); i++ {
		tileset := &a.World.Defs.Tilesets[i]
		tilesetPath := filepath.Join(LDTKpath, tileset.RelPath)
		tilesetImage, _, err := ebitenutil.NewImageFromFile(tilesetPath)
		if err != nil {
			return err
		}
		tileset.Image = tilesetImage
	}
	return nil
}

func NewLDTKAsset(path string) *ebitenLDTK.World {
	asset, exists := assetloader.Exists(path)
	if exists {
		return &asset.(*LDTKAsset).World
	}

	LDTKasset := LDTKAsset{
		path: path,
	}

	assetloader.Load(path, &LDTKasset)
	return &LDTKasset.World
}
