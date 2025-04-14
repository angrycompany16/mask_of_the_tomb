package assettypes

import (
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/core/assetloader"
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
	LDTKasset := LDTKAsset{
		path: path,
	}

	assetloader.AddAsset(&LDTKasset)
	return &LDTKasset.World
}
