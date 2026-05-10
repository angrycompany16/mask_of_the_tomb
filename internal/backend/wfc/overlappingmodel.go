package wfc

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/backend/maths"
	"math"

	"github.com/aquilax/go-perlin"
	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/hajimehoshi/ebiten/v2"
)

type D8Symmetry int

const (
	ID D8Symmetry = iota
	ROT90
	ROT180
	ROT270
	FLIP_X
	FLIP_Y
	FLIP_DIAG_TL_BR
	FLIP_DIAG_BL_TR
)

type OverlappingModelWFC struct {
	src                       *image.RGBA
	windowSize                int
	outputWidth, outputHeight int
	includeReflections        bool
	includeRotations          bool
	noiseGen                  *perlin.Perlin
	outputImage               *ebiten.Image
	entropyHeap               *binaryheap.Heap
	tiles                     []*Tile
	outputTilemap             [][]*Cell
}

// Represents a possible tile for the (essentially) simple tile algorithm
type Tile struct {
	*image.RGBA
	frequency int
	// total_freq         int
	neighbourRuleUp    []bool
	neighbourRuleDown  []bool
	neighbourRuleLeft  []bool
	neighbourRuleRight []bool
}

type EntropyCoord struct {
	x, y    int
	entropy *float64
}

// Represents a tile in the output image
type Cell struct {
	tileIndex               int
	entropy                 float64
	possibleIndices         []bool
	possibleTilesFreqSum    float64
	possibleTilesFreqLogSum float64
	isCollapsed             bool
}

func (o *OverlappingModelWFC) Generate() {
	o.preprocess()
}

// Chain of execution for algorithm
func (o *OverlappingModelWFC) preprocess() {
	// 1. Get the full list of tiles from the image
	for i := range o.src.Bounds().Dx() {
		for j := range o.src.Bounds().Dy() {
			r := image.Rect(i, j, i+o.windowSize, j+o.windowSize)

			if !r.In(o.src.Bounds()) {
				newImg := wrappedSubImage(o.src, r)
				o.tiles = append(o.tiles, newTile(newImg))
			} else {
				newImg := o.src.SubImage(r).(*image.RGBA)
				o.tiles = append(o.tiles, newTile(newImg))
			}
		}
	}

	// 2. Remove duplicates
	o.dedupTileAtlas(true)

	// 3. Add symmetries
	fmt.Println(len(o.tiles))
	for i := range o.tiles {
		sym := o.computeSymmetries(i)
		if o.includeRotations {
			o.tiles = append(o.tiles, sym[0], sym[1], sym[2])
		}

		if o.includeReflections {
			o.tiles = append(o.tiles, sym[3], sym[4], sym[5], sym[6])
		}
	}

	// 4. Remove duplicates, again
	o.dedupTileAtlas(false)

	// 5. Create a list detailing which tiles can be next to which other
	for i := range o.tiles {
		o.tiles[i].neighbourRuleUp = make([]bool, len(o.tiles))
		o.tiles[i].neighbourRuleDown = make([]bool, len(o.tiles))
		o.tiles[i].neighbourRuleLeft = make([]bool, len(o.tiles))
		o.tiles[i].neighbourRuleRight = make([]bool, len(o.tiles))
	}

	for i := range o.tiles {
		for j := range o.tiles[i].neighbourRuleUp {
			if o.tiles[i].neighbourRuleUp[j] ||
				o.tiles[i].neighbourRuleDown[j] ||
				o.tiles[i].neighbourRuleLeft[j] ||
				o.tiles[i].neighbourRuleRight[j] {
				continue
			}

			if o.tiles[i].overlaps(o.tiles[j], maths.DirUp) {
				o.tiles[i].neighbourRuleUp[j] = true
				o.tiles[j].neighbourRuleDown[i] = true
			}
			if o.tiles[i].overlaps(o.tiles[j], maths.DirDown) {
				o.tiles[i].neighbourRuleDown[j] = true
				o.tiles[j].neighbourRuleUp[i] = true
			}
			if o.tiles[i].overlaps(o.tiles[j], maths.DirLeft) {
				o.tiles[i].neighbourRuleLeft[j] = true
				o.tiles[j].neighbourRuleRight[i] = true
			}
			if o.tiles[i].overlaps(o.tiles[j], maths.DirRight) {
				o.tiles[i].neighbourRuleRight[j] = true
				o.tiles[j].neighbourRuleLeft[i] = true
			}
		}
	}
}

func (o *OverlappingModelWFC) wfcCore() {
	// So we are ready to begin our process
	// 1. Initialize the grid with proper dimensions and such
	o.outputTilemap = make([][]*Cell, o.outputHeight)
	for y := range o.outputHeight {
		o.outputTilemap[y] = make([]*Cell, o.outputWidth)
		for x := range o.outputWidth {
			newCell := newEmptyCell(x, y, o.tiles, o.noiseGen)
			o.outputTilemap[y][x] = newCell
			o.entropyHeap.Push(EntropyCoord{
				x:       x,
				y:       y,
				entropy: &newCell.entropy,
			})
		}
	}

	// 2. Start the iteration
	// The loop:
	// Pop the (first?) element off of the entropy heap, geat the corresp.
	// cell from the map
	// Collapse that cell (based on frequency)
	// Propagate the changes
	// Keep running count over how many cells are left to collapse
	// When count reaches zero, terminate
}

func (o *OverlappingModelWFC) postprocess() {
	// For each cell, render the upper leftmost pixel to the screen
}

// (tileA is viewed as center)
func (t *Tile) overlaps(tileB *Tile, dir maths.Direction) bool {
	bA, bB := t.Bounds(), tileB.Bounds()

	rectA := image.Rect(bA.Min.X, bA.Min.Y, bA.Max.X, bA.Max.Y)
	rectB := image.Rect(bB.Min.X, bB.Min.Y, bB.Max.X, bB.Max.Y)

	switch dir {
	case maths.DirUp:
		rectA.Max.Y -= 1
		rectB.Min.Y += 1
	case maths.DirDown:
		rectA.Min.Y += 1
		rectB.Max.Y -= 1
	case maths.DirLeft:
		rectA.Max.X -= 1
		rectB.Min.X += 1
	case maths.DirRight:
		rectA.Min.X += 1
		rectB.Max.X -= 1
	}

	return equal(
		t.SubImage(rectA).(*image.RGBA),
		tileB.SubImage(rectB).(*image.RGBA),
	)
}

func (c *Cell) RemovePossibles(tiles []*Tile, indices ...int) {
	for i := range indices {
		c.possibleIndices[i] = false
		c.possibleTilesFreqSum -= float64(tiles[i].frequency)
		c.possibleTilesFreqLogSum -= float64(tiles[i].frequency) * math.Log(float64(tiles[i].frequency))
	}
}

func (c *Cell) calculateEntropy() float64 {
	return math.Log2(c.possibleTilesFreqSum) - c.possibleTilesFreqLogSum/c.possibleTilesFreqSum
}

func (o *OverlappingModelWFC) computeSymmetries(tileindex int) [7]*Tile {
	tile := o.tiles[tileindex]

	rot90 := image.NewRGBA(image.Rect(0, 0, o.windowSize, o.windowSize))
	rot180 := image.NewRGBA(image.Rect(0, 0, o.windowSize, o.windowSize))
	rot270 := image.NewRGBA(image.Rect(0, 0, o.windowSize, o.windowSize))
	flipX := image.NewRGBA(image.Rect(0, 0, o.windowSize, o.windowSize))
	flipY := image.NewRGBA(image.Rect(0, 0, o.windowSize, o.windowSize))
	diag1 := image.NewRGBA(image.Rect(0, 0, o.windowSize, o.windowSize))
	diag2 := image.NewRGBA(image.Rect(0, 0, o.windowSize, o.windowSize))

	bounds := tile.Bounds()

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			col := tile.At(bounds.Min.X+x, bounds.Min.Y+y)

			rot90.Set(bounds.Dy()-1-y, x, col)
			rot180.Set(bounds.Dx()-1-x, bounds.Dy()-1-y, col)
			rot270.Set(y, bounds.Dx()-1-x, col)
			flipX.Set(bounds.Dx()-1-x, y, col)
			flipY.Set(x, bounds.Dy()-1-y, col)
			diag1.Set(y, x, col)
			diag2.Set(bounds.Dy()-1-y, bounds.Dx()-1-x, col)
		}
	}

	return [7]*Tile{
		newTile(rot90),
		newTile(rot180),
		newTile(rot270),
		newTile(flipX),
		newTile(flipY),
		newTile(diag1),
		newTile(diag2),
	}
}

func (o *OverlappingModelWFC) dedupTileAtlas(countFrequency bool) {
	duplicateCount := 0
	for i := range o.tiles {
		if i >= len(o.tiles)-duplicateCount {
			break
		}

		for j := range o.tiles {
			if j == i {
				continue
			}

			if j >= len(o.tiles)-duplicateCount {
				break
			}

			if equal(o.tiles[i].RGBA, o.tiles[j].RGBA) {
				// fmt.Println(i, j)
				if countFrequency {
					o.tiles[i].frequency += 1
				}
				duplicateCount += 1
				o.tiles[j] = o.tiles[len(o.tiles)-duplicateCount]
			}
		}
	}
	o.tiles = o.tiles[:len(o.tiles)-duplicateCount]
	fmt.Printf("Removed %d duplicates\n", duplicateCount)
}

func equal(img1, img2 *image.RGBA) bool {
	b1, b2 := img1.Bounds(), img2.Bounds()
	if b1.Dx() != b2.Dx() || b1.Dy() != b2.Dy() {
		fmt.Println("Different image sizes...")
		return false
	}

	for x := 0; x < b1.Dx(); x++ {
		for y := 0; y < b1.Dy(); y++ {
			p1 := img1.At(b1.Min.X+x, b1.Min.Y+y)
			p2 := img2.At(b2.Min.X+x, b2.Min.Y+y)
			if p1 != p2 {
				return false
			}
		}
	}

	return true
}

func (o *OverlappingModelWFC) DrawTileAtlas() *ebiten.Image {
	w := o.src.Bounds().Dx()
	h := o.src.Bounds().Dy()

	img := ebiten.NewImage(w*o.windowSize, h*o.windowSize)

	for i, tile := range o.tiles {
		y := (i % h) * o.windowSize
		x := (i / w) * o.windowSize
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		ebitenTile := ebiten.NewImageFromImage(tile)
		img.DrawImage(ebitenTile, op)
	}
	return img
}

// Purpose-built for wavefcn algorithm
func wrappedSubImage(src *image.RGBA, rect image.Rectangle) *image.RGBA {
	newImg := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))

	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			x := i % src.Bounds().Dx()
			y := j % src.Bounds().Dy()
			col := src.At(x, y)
			newImg.Set(i-rect.Min.X, j-rect.Min.Y, col)
		}
	}

	return newImg
}

func newEmptyCell(x, y int, tiles []*Tile, noiseGen *perlin.Perlin) *Cell {
	newCell := &Cell{
		tileIndex:       -1,
		entropy:         0.001 * noiseGen.Noise2D(float64(x), float64(y)),
		possibleIndices: make([]bool, len(tiles)),
		isCollapsed:     false,
	}
	for i := range tiles {
		newCell.possibleIndices[i] = true
		newCell.possibleTilesFreqSum += float64(tiles[i].frequency)
		newCell.possibleTilesFreqLogSum += float64(tiles[i].frequency) * math.Log2(float64(tiles[i].frequency))
	}
	return newCell
}

func newTile(image *image.RGBA) *Tile {
	return &Tile{
		RGBA:      image,
		frequency: 1,
	}
}

func NewOverlappingModel(
	src *image.RGBA,
	windowSize, outputWidth, outputHeight int,
	includeRotations, includeReflections bool,
	seed int64,
) *OverlappingModelWFC {
	wfc := &OverlappingModelWFC{
		src:                src,
		windowSize:         windowSize,
		outputWidth:        outputWidth,
		outputHeight:       outputHeight,
		includeRotations:   includeRotations,
		includeReflections: includeReflections,
		noiseGen:           perlin.NewPerlin(2, 2, 4, seed),
		entropyHeap:        binaryheap.NewWithIntComparator(),
		tiles:              make([]*Tile, 0),
		outputTilemap:      make([][]*Cell, 0),
	}
	return wfc
}
