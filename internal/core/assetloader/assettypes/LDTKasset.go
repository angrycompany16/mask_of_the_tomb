package assettypes

import (
	"fmt"
	"mask_of_the_tomb/internal/core/assetloader"
	"path/filepath"
	"strings"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// TODO: Add funcitonality for loading LDTK from bytes (go embed)
type LDTKAsset struct {
	path     string
	World    ebitenLDTK.World
	Tilesets map[string]*ebiten.Image
}

func (a *LDTKAsset) Load() error {
	world, err := ebitenLDTK.LoadWorld(a.path)

	if err != nil {
		fmt.Println("error when loading world")
		return err
	}

	a.Tilesets = make(map[string]*ebiten.Image, len(a.World.Defs.Tilesets))
	a.World = world
	LDTKpath := filepath.Clean(filepath.Join(a.path, ".."))
	for i := 0; i < len(a.World.Defs.Tilesets); i++ {
		tileset := &a.World.Defs.Tilesets[i]
		tilesetPath := filepath.Join(LDTKpath, tileset.RelPath)

		if !strings.HasSuffix(tilesetPath, ".png") {
			fmt.Println("Skipping non-png tileset")
			continue
		}

		tilesetImage, _, err := ebitenutil.NewImageFromFile(tilesetPath)
		if err != nil {
			return err
		}
		a.Tilesets[tileset.Name] = tilesetImage
	}
	return nil
}

func GetLDTKAsset(name string) (*LDTKAsset, error) {
	ldtkAsset, err := assetloader.GetAsset(name)
	return ldtkAsset.(*LDTKAsset), err
}

func NewLDTKAsset(path string) *LDTKAsset {
	return &LDTKAsset{
		path: path,
	}
}
