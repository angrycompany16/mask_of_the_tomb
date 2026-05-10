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

// Let's not make this generic - We'll use simple images instead
type WaveFunctionSetup struct {
	GridSize      float64
	Modules       []*Module
	Tilemap       [][]*Tile
	Width, Height int
	Tileset       *ebiten.Image
	rng           *rand.Rand
}

// Tile contains a slice of integers which index into the Modules slice,
// specifying which possible modules it can represent. When a tile is
// collapsed, CollapsedState gets set to the corresponding module index,
// or -1 in the case of a failure
type Tile struct {
	PossibleModules []bool
	Entropy         float64
	CollapsedState  int
	IsCollapsed     bool
}

// Represents a single tile, with its constraints
type Module struct {
	tilesetRegion image.Rectangle
	rules         []*DirectionalRule // Could in theory use a map?
}

// Represents a constraint where only certain tiles are permitted
// in the direction specified. These are given as integers which
// index into the Modules slice of the WaveFunctionSetup.
type DirectionalRule struct {
	direction maths.Direction
	permitted []int
}

func (w *WaveFunctionSetup) AddModule(module *Module) *WaveFunctionSetup {
	w.Modules = append(w.Modules, module)
	return w
}

// Spawns a tilemap where each tile is in a totally unconstrained state
func (w *WaveFunctionSetup) InitTiles() {
	w.Tilemap = make([][]*Tile, w.Height)
	for i := range w.Tilemap {
		w.Tilemap[i] = make([]*Tile, w.Width)
		for j := range w.Tilemap[i] {
			w.Tilemap[i][j] = NewEmptyTile(w.Modules)
		}
	}
}

// Collects the tiles into a cohesive image.
func (w *WaveFunctionSetup) MakeImage() *ebiten.Image {
	image := ebiten.NewImage(w.Width*int(w.GridSize), w.Height*int(w.GridSize))
	for i := range w.Tilemap {
		for j := range w.Tilemap[i] {
			tile := w.Tilemap[i][j]
			if tile.CollapsedState == -1 || !tile.IsCollapsed {
				continue
			}
			region := w.Modules[tile.CollapsedState].tilesetRegion
			x, y := w.GridSize*float64(i), w.GridSize*float64(j)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y)
			image.DrawImage(w.Tileset.SubImage(region).(*ebiten.Image), op)
		}
	}
	return image
}

// Starts the WFC algorithm. Takes a starting index which is the point
// where the iteration loop will begin.
func (w *WaveFunctionSetup) Collapse(startI, startJ int) {
	w.Tilemap[startI][startJ].Entropy = 0
	for w.Iterate() {
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
func (w *WaveFunctionSetup) Iterate() bool {
	if w.IsFinished() {
		return false
	}
	i, j := w.FindLowestEntropy()
	ok := w.Tilemap[i][j].Collapse(w.rng)

	if !ok {
		return false
	}

	w.PropagateTilemap(i, j)
	return true
}

func (w *WaveFunctionSetup) IsFinished() bool {
	isFinished := true
	for i := range w.Tilemap {
		for j := range w.Tilemap[i] {
			isFinished = isFinished && w.Tilemap[i][j].IsCollapsed
		}
	}
	return isFinished
}

func (w *WaveFunctionSetup) FindLowestEntropy() (resI, resJ int) {
	maxEntropy := math.Inf(1)
	for i := range w.Tilemap {
		for j := range w.Tilemap[i] {
			if w.Tilemap[i][j].IsCollapsed {
				continue
			}

			if w.Tilemap[i][j].Entropy < maxEntropy {
				maxEntropy = w.Tilemap[i][j].Entropy
				resI = i
				resJ = j
			}
		}
	}
	return
}

// Propagates the tilemap after collapsing the tile at index (collapsedI,
// collapsedJ)
func (w *WaveFunctionSetup) PropagateTilemap(collapsedI, collapsedJ int) {
	// Due to only having four-directional movement we only need these
	// directions
	directions := []maths.Vec2{
		maths.NewVec2(1, 0),
		maths.NewVec2(0, 1),
		maths.NewVec2(-1, 0),
		maths.NewVec2(0, -1),
	}
	collapsedTile := w.Tilemap[collapsedI][collapsedJ]
	for _, dir := range directions {
		i := collapsedI + int(dir.Y)
		j := collapsedJ + int(dir.X)
		dirEnum := maths.DirFromVector(dir.X, dir.Y)

		outOfBounds := (i < 0 || i >= len(w.Tilemap) || j < 0 || j >= len(w.Tilemap[0]))
		if outOfBounds {
			continue
		}

		isCollapsed := w.Tilemap[i][j].IsCollapsed
		if isCollapsed {
			continue
		}

		tile := w.Tilemap[i][j]
		for _, rule := range w.Modules[collapsedTile.CollapsedState].rules {
			if rule.direction != dirEnum {
				continue
			}

			for moduleIndex, possible := range tile.PossibleModules {
				if !possible {
					continue
				}
				if !slices.Contains(rule.permitted, moduleIndex) {
					tile.PossibleModules[moduleIndex] = false
				}
			}
		}
		tile.ComputeEntropy()
	}
}

// Compute entropy in simple way: Just take the sum of possible modules.
func (t *Tile) ComputeEntropy() {
	s := 0.0
	for _, possible := range t.PossibleModules {
		if possible {
			s += 1.0
		}
	}
	t.Entropy = s
}

// Collapses a tile by (deterministically!) choosing the (note necessarily) first possible
// modules from the tile's possible modules. Returns whether the tile
// failed to collapse or not.
func (t *Tile) Collapse(rng *rand.Rand) bool {
	// Choose random tile from possibles
	possibles := make([]int, 0)
	for i, possible := range t.PossibleModules {
		if !possible {
			continue
		}
		possibles = append(possibles, i)
	}

	if len(possibles) == 0 {
		t.CollapsedState = -1
		t.IsCollapsed = true
		return false
	}

	n := rng.Intn(len(possibles))
	t.CollapsedState = possibles[n]
	t.IsCollapsed = true
	return true
}

func NewWFC(gridSize float64, width, height int, tileset *ebiten.Image, seed int64) *WaveFunctionSetup {
	wfc := &WaveFunctionSetup{
		GridSize: gridSize,
		Width:    width,
		Height:   height,
		Modules:  make([]*Module, 0),
		Tileset:  tileset,
		rng:      rand.New(rand.NewSource(seed)),
	}
	return wfc
}

func NewModule(tilesetRegion image.Rectangle, rules ...*DirectionalRule) *Module {
	return &Module{
		tilesetRegion: tilesetRegion,
		rules:         rules,
	}
}

func NewDirectionalRule(direction maths.Direction, permitted ...int) *DirectionalRule {
	return &DirectionalRule{
		direction: direction,
		permitted: permitted,
	}
}

// Creates a state where all modules are possible
func NewEmptyTile(modules []*Module) *Tile {
	newState := Tile{}
	newState.PossibleModules = make([]bool, len(modules))
	for i := range modules {
		newState.PossibleModules[i] = true
	}
	newState.Entropy = float64(len(modules))
	return &newState
}
