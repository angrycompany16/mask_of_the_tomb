package assettypes

import (
	"fmt"
	"image"
	"io/fs"
	"path/filepath"
	"strings"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/hajimehoshi/ebiten/v2"
)

type LDTKData struct {
	World    *ebitenLDTK.World
	Tilesets map[string]*ebiten.Image
}

func NewLDTKData(world *ebitenLDTK.World, tilesets map[string]*ebiten.Image) *LDTKData {
	return &LDTKData{
		World:    world,
		Tilesets: tilesets,
	}
}

type LDTKAsset struct {
	srcPath string
}

// This needs to change - Everything is fucked
func (a *LDTKAsset) Load(fs fs.FS) (any, error) {
	f, err := fs.Open(a.srcPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	world, err := ebitenLDTK.LoadFromStream(f)

	if err != nil {
		return nil, err
	}

	tilesets := make(map[string]*ebiten.Image, len(world.Defs.Tilesets))
	LDTKpath := filepath.Clean(filepath.Join(a.srcPath, ".."))
	for i := 0; i < len(world.Defs.Tilesets); i++ {
		tileset := &world.Defs.Tilesets[i]
		tilesetPath := filepath.Join(LDTKpath, tileset.RelPath)
		// fmt.Println(tilesetPath)

		if !strings.HasSuffix(tilesetPath, ".png") {
			fmt.Println("Tileset loading: Skipping non-png file")
			continue
		}

		f, err := fs.Open(tilesetPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return nil, err
		}

		tilesets[tileset.Name] = ebiten.NewImageFromImage(img)
	}
	return NewLDTKData(&world, tilesets), nil
}

func NewLDTKAsset(srcPath string) *LDTKAsset {
	return &LDTKAsset{
		srcPath: srcPath,
	}
}
