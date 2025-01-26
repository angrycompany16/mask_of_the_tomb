package game

import (
	"fmt"
	"image"
	"log"
	"mask_of_the_tomb/ebitenLDTK"
	. "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/files"
	"mask_of_the_tomb/player"
	. "mask_of_the_tomb/utils"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpaceLayerName   = "Playerspace"
	spawnPosEntityName     = "SpawnPosition"
	doorEntityName         = "Door"
	doorOtherSideFieldName = "OtherSide"
)

type World struct {
	worldLDTK       ebitenLDTK.LDTKWorld
	activeLevel     *ebitenLDTK.LDTKLevel
	currentTileSize float64
	tiles           [][]int
}

func (w *World) Init() {
	w.worldLDTK = *files.LazyLDTK(files.LDTKMapPath)

	w.activeLevel = &w.worldLDTK.Levels[0]
	w.tiles = w.worldLDTK.MakeBitmapFromLayer(w.activeLevel, playerSpaceLayerName)

	// One folder back to access LDTK folder
	LDTKpath := path.Clean(path.Join(files.LDTKMapPath, ".."))
	for i := 0; i < len(w.worldLDTK.Defs.Tilesets); i++ {
		tileset := &w.worldLDTK.Defs.Tilesets[i]
		tilesetPath := path.Join(LDTKpath, tileset.RelPath)
		tileset.Image = files.LazyImage(tilesetPath)
	}

	changeActiveLevel(w, 0)
}

func (w *World) Update() {
	// Anything...
}

func (w *World) Draw(surf *ebiten.Image) {
	for i := len(w.activeLevel.LayerInstances) - 1; i >= 0; i-- {
		layerInstance := w.activeLevel.LayerInstances[i]

		layer, err := w.worldLDTK.GetLayerByUid(layerInstance.LayerDefUid)
		if err != nil {
			log.Fatal(err)
		}

		if layer.Type != ebitenLDTK.LayerTypeTiles {
			continue
		}

		tileset, err := w.worldLDTK.GetTilesetByUid(layer.TilesetUid)
		if err != nil {
			log.Fatal(err)
		}

		tileSize := tileset.TileGridSize
		for _, tile := range layerInstance.GridTiles {
			scaleX, scaleY := 1.0, 1.0
			switch tile.TileOrientation {
			case ebitenLDTK.OrientationFlipX:
				scaleX = -1
			case ebitenLDTK.OrientationFlipY:
				scaleY = -1
			case ebitenLDTK.OrientationFlipXY:
				scaleX, scaleY = -1, -1
			}
			DrawAtScaled(tileset.Image.SubImage(
				image.Rect(
					tile.Src[0],
					tile.Src[1],
					tile.Src[0]+tileSize,
					tile.Src[1]+tileSize,
				),
			).(*ebiten.Image), surf, F64(tile.Px[0]), F64(tile.Px[1]), scaleX, scaleY, 0.5, 0.5)
		}
	}
}

func (w *World) GetSpawnPoint() (float64, float64) {
	for _, layerInstance := range w.activeLevel.LayerInstances {
		for _, entityInstance := range layerInstance.EntityInstances {
			entity, err := w.worldLDTK.GetEntityByUid(entityInstance.Uid)
			HandleLazy(err)
			if entity.Name != spawnPosEntityName {
				continue
			}
			return F64(entityInstance.Px[0]), F64(entityInstance.Px[1])
		}
	}
	return 0, 0
}

func (w *World) getCollision(moveDir player.MoveDirection, x, y float64) (float64, float64) {
	gridX, gridY := w.worldToGrid(x, y)
	switch moveDir {
	case player.DirUp:
		for i := gridY; i >= 0; i-- {
			if w.tiles[i][gridX] == 1 {
				newX, newY := w.gridToWorld(gridX, i+1)
				return newX, newY
			}
		}
		return x, y
	case player.DirDown:
		for i := gridY; i < len(w.tiles); i++ {
			if w.tiles[i][gridX] == 1 {
				newX, newY := w.gridToWorld(gridX, i-1)
				return newX, newY
			}
		}
		return x, y
	case player.DirLeft:
		for i := gridX; i >= 0; i-- {
			if w.tiles[gridY][i] == 1 {
				newX, newY := w.gridToWorld(i+1, gridY)
				return newX, newY
			}
		}
		return x, y
	case player.DirRight:
		for i := gridX; i < len(w.tiles[0]); i++ {
			if w.tiles[gridY][i] == 1 {
				newX, newY := w.gridToWorld(i-1, gridY)
				return newX, newY
			}
		}
		return x, y
	default:
		return x, y
	}
}

func (w *World) TryDoorOverlap(x, y float64) (bool, ebitenLDTK.LDTKEntityInstance) {
	for _, layerInstance := range w.activeLevel.LayerInstances {
		for _, entityInstance := range layerInstance.EntityInstances {
			entity, err := w.worldLDTK.GetEntityByUid(entityInstance.Uid)
			HandleLazy(err)
			if entity.Name != doorEntityName {
				continue
			}
			if F64(entityInstance.Px[0]) == x && F64(entityInstance.Px[1]) == y {
				return true, entityInstance
			}
		}
	}
	return false, ebitenLDTK.LDTKEntityInstance{}
}

func (w *World) ExitByDoor(doorEntity ebitenLDTK.LDTKEntityInstance) (int, int) {
	entity, err := w.worldLDTK.GetEntityByUid(doorEntity.Uid)
	HandleLazy(err)
	if entity.Name != doorEntityName {
		return 0, 0 // Should honestly be error handled
	}
	for _, fieldInstance := range doorEntity.FieldInstances {
		if fieldInstance.Name != doorOtherSideFieldName {
			continue
		}

		nextLevel, err := w.worldLDTK.GetLevelByIid(fieldInstance.EntityRefValue.LevelIid)
		HandleLazy(err)
		// Change level
		changeActiveLevel(w, nextLevel.Uid)

		otherSideEntityRef, err := w.activeLevel.GetEntityInstanceByIid(fieldInstance.EntityRefValue.EntityIid)
		HandleLazy(err)
		return otherSideEntityRef.Px[0], otherSideEntityRef.Px[1]
	}
	return 0, 0
}

func (w *World) gridToWorld(x, y int) (float64, float64) {
	return F64(x) * w.currentTileSize, F64(y) * w.currentTileSize
}

func (w *World) worldToGrid(x, y float64) (int, int) {
	return int(x / w.currentTileSize), int(y / w.currentTileSize)
}

func changeActiveLevel[T string | int](world *World, id T) error {
	var newLevel ebitenLDTK.LDTKLevel
	var err error

	switch v := any(id).(type) {
	case string:
		newLevel, err = world.worldLDTK.GetLevelByName(v)
		if err != nil {
			fmt.Println("Error when switching levels (id string)")
			return err
		}
	case int:
		newLevel, err = world.worldLDTK.GetLevelByUid(v)
		if err != nil {
			fmt.Println("Error when switching levels (id int)")
			return err
		}
	}

	world.activeLevel = &newLevel

	playerspace, err := world.activeLevel.GetLayerInstanceByName(playerSpaceLayerName)
	if err != nil {
		world.currentTileSize = 1
		return nil
	}

	layer, err := world.worldLDTK.GetLayerByUid(playerspace.LayerDefUid)
	HandleLazy(err)
	world.currentTileSize = F64(layer.GridSize)
	return nil
}
