package assettypes

import (
	"mask_of_the_tomb/internal/core/assetloader"
	"mask_of_the_tomb/internal/core/errs"
	"path/filepath"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

type LDTKAsset struct {
	path  string
	World ebitenLDTK.World
}

func (a *LDTKAsset) Load() {
	a.World = errs.Must(ebitenLDTK.LoadWorld(a.path))

	LDTKpath := filepath.Clean(filepath.Join(a.path, ".."))
	for i := 0; i < len(a.World.Defs.Tilesets); i++ {
		tileset := &a.World.Defs.Tilesets[i]
		tilesetPath := filepath.Join(LDTKpath, tileset.RelPath)
		tileset.Image = errs.MustNewImageFromFile(tilesetPath)
	}
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
