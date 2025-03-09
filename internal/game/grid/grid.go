package grid

type Tilemap[T any] struct {
	Tiles    [][]T
	TileSize float64
}

func (g *Tilemap[any]) GridPosToWorld(x, y int) (float64, float64) {
	return float64(x) * g.TileSize, float64(y) * g.TileSize
}

func (g *Tilemap[any]) WorldPosToGrid(x, y float64) (int, int) {
	return int(x / g.TileSize), int(y / g.TileSize)
}

func NewGrid[T any](tileSize float64, tiles [][]T) Tilemap[T] {
	return Tilemap[T]{
		Tiles:    tiles,
		TileSize: tileSize,
	}
}
