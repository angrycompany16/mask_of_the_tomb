package assetloader

import (
	"mask_of_the_tomb/internal/errs"
	"path/filepath"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
)

type LDTKAsset struct {
	path string
	// TODO: Make into a pointer
	World ebitenLDTK.World
	// done  chan int
}

func (a *LDTKAsset) load() {
	a.World = errs.Must(ebitenLDTK.LoadWorld(a.path))

	LDTKpath := filepath.Clean(filepath.Join(a.path, ".."))
	for i := 0; i < len(a.World.Defs.Tilesets); i++ {
		tileset := &a.World.Defs.Tilesets[i]
		tilesetPath := filepath.Join(LDTKpath, tileset.RelPath)
		tileset.Image = errs.MustNewImageFromFile(tilesetPath)
	}
}

func NewLDTKAsset(path string) *LDTKAsset {
	return &LDTKAsset{
		path: path,
	}
}
