package wfc

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/backend/maths"
	"math"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

type SimpleTileWFC struct {
	gridSize      float64
	modules       []*module
	tilemap       [][]*simpleTile
	Width, Height int
	tileset       *ebiten.Image
	rng           *rand.Rand
}

// simpleTile contains a slice of integers which index into the Modules slice,
// specifying which possible modules it can represent. When a simpleTile is
// collapsed, CollapsedState gets set to the corresponding module index,
// or -1 in the case of a failure
type simpleTile struct {
	possibleModules []bool
	entropy         float64
	collapsedState  int
	isCollapsed     bool
}

// Represents a single tile, with its constraints
type module struct {
	tilesetRegion image.Rectangle
	rules         []*directionalRule // Could in theory use a map?
}

// Represents a constraint where only certain tiles are permitted
// in the direction specified. These are given as integers which
// index into the Modules slice of the WaveFunctionSetup.
type directionalRule struct {
	direction maths.Direction
	permitted []int
}

func (w *SimpleTileWFC) AddModule(module *module) *SimpleTileWFC {
	w.modules = append(w.modules, module)
	return w
}

// Spawns a tilemap where each tile is in a totally unconstrained state
func (w *SimpleTileWFC) InitTiles() {
	w.tilemap = make([][]*simpleTile, w.Height)
	for i := range w.tilemap {
		w.tilemap[i] = make([]*simpleTile, w.Width)
		for j := range w.tilemap[i] {
			w.tilemap[i][j] = newEmptyTile(w.modules)
		}
	}
}

// Collects the tiles into a cohesive image.
func (w *SimpleTileWFC) MakeImage() *ebiten.Image {
	image := ebiten.NewImage(w.Width*int(w.gridSize), w.Height*int(w.gridSize))
	for i := range w.tilemap {
		for j := range w.tilemap[i] {
			tile := w.tilemap[i][j]
			if tile.collapsedState == -1 || !tile.isCollapsed {
				continue
			}
			region := w.modules[tile.collapsedState].tilesetRegion
			x, y := w.gridSize*float64(i), w.gridSize*float64(j)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y)
			image.DrawImage(w.tileset.SubImage(region).(*ebiten.Image), op)
		}
	}
	return image
}

// Starts the WFC algorithm. Takes a starting index which is the point
// where the iteration loop will begin.
func (w *SimpleTileWFC) Collapse(startI, startJ int) {
	w.tilemap[startI][startJ].entropy = 0
	for w.iterate() {
		fmt.Println("Updating WFC")
	}
}

// The iteration process is like this:
//  1. Find a tile to collapse (the lowest entropy one)
//  2. Collapse the tile
//  3. Propagate, i.e. update the possible modules of neighbouring tile
//     depending on the rules of the recently collapsed tile. Also recompute
//     entropy.
//
// Return whether a tile got into a contradictory state or not.
func (w *SimpleTileWFC) iterate() bool {
	if w.isFinished() {
		return false
	}
	i, j := w.findLowestEntropy()
	ok := w.tilemap[i][j].collapse(w.rng)

	if !ok {
		return false
	}

	w.propagateTilemap(i, j)
	return true
}

func (w *SimpleTileWFC) isFinished() bool {
	isFinished := true
	for i := range w.tilemap {
		for j := range w.tilemap[i] {
			isFinished = isFinished && w.tilemap[i][j].isCollapsed
		}
	}
	return isFinished
}

func (w *SimpleTileWFC) findLowestEntropy() (resI, resJ int) {
	maxEntropy := math.Inf(1)
	for i := range w.tilemap {
		for j := range w.tilemap[i] {
			if w.tilemap[i][j].isCollapsed {
				continue
			}

			if w.tilemap[i][j].entropy < maxEntropy {
				maxEntropy = w.tilemap[i][j].entropy
				resI = i
				resJ = j
			}
		}
	}
	return
}

// Propagates the tilemap after collapsing the tile at index (collapsedI,
// collapsedJ)
func (w *SimpleTileWFC) propagateTilemap(collapsedI, collapsedJ int) {
	// Due to only having four-directional movement we only need these
	// directions
	directions := []maths.Vec2{
		maths.NewVec2(1, 0),
		maths.NewVec2(0, 1),
		maths.NewVec2(-1, 0),
		maths.NewVec2(0, -1),
	}
	collapsedTile := w.tilemap[collapsedI][collapsedJ]
	for _, dir := range directions {
		i := collapsedI + int(dir.Y)
		j := collapsedJ + int(dir.X)
		dirEnum := maths.DirFromVector(dir.X, dir.Y)

		outOfBounds := (i < 0 || i >= len(w.tilemap) || j < 0 || j >= len(w.tilemap[0]))
		if outOfBounds {
			continue
		}

		isCollapsed := w.tilemap[i][j].isCollapsed
		if isCollapsed {
			continue
		}

		tile := w.tilemap[i][j]
		for _, rule := range w.modules[collapsedTile.collapsedState].rules {
			if rule.direction != dirEnum {
				continue
			}

			for moduleIndex, possible := range tile.possibleModules {
				if !possible {
					continue
				}
				if !slices.Contains(rule.permitted, moduleIndex) {
					tile.possibleModules[moduleIndex] = false
				}
			}
		}
		tile.computeEntropy()
	}
}

// Compute entropy in simple way: Just take the sum of possible modules.
func (t *simpleTile) computeEntropy() {
	s := 0.0
	for _, possible := range t.possibleModules {
		if possible {
			s += 1.0
		}
	}
	t.entropy = s
}

// Collapses a tile by (deterministically!) choosing the (note necessarily) first possible
// modules from the tile's possible modules. Returns whether the tile
// failed to collapse or not.
func (t *simpleTile) collapse(rng *rand.Rand) bool {
	// Choose random tile from possibles
	possibles := make([]int, 0)
	for i, possible := range t.possibleModules {
		if !possible {
			continue
		}
		possibles = append(possibles, i)
	}

	if len(possibles) == 0 {
		t.collapsedState = -1
		t.isCollapsed = true
		return false
	}

	n := rng.Intn(len(possibles))
	t.collapsedState = possibles[n]
	t.isCollapsed = true
	return true
}

func NewSimpleTiled(gridSize float64, width, height int, tileset *ebiten.Image, seed int64) *SimpleTileWFC {
	wfc := &SimpleTileWFC{
		gridSize: gridSize,
		Width:    width,
		Height:   height,
		modules:  make([]*module, 0),
		tileset:  tileset,
		rng:      rand.New(rand.NewSource(seed)),
	}
	return wfc
}

func NewModule(tilesetRegion image.Rectangle, rules ...*directionalRule) *module {
	return &module{
		tilesetRegion: tilesetRegion,
		rules:         rules,
	}
}

func NewDirectionalRule(direction maths.Direction, permitted ...int) *directionalRule {
	return &directionalRule{
		direction: direction,
		permitted: permitted,
	}
}

// Creates a state where all modules are possible
func newEmptyTile(modules []*module) *simpleTile {
	newState := simpleTile{}
	newState.possibleModules = make([]bool, len(modules))
	for i := range modules {
		newState.possibleModules[i] = true
	}
	newState.entropy = float64(len(modules))
	return &newState
}
