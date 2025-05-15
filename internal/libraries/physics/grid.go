package physics

type tilemap[T any] struct {
	Tiles    [][]T
	TileSize float64
}

func (g *tilemap[any]) gridPosToWorld(x, y int) (float64, float64) {
	return float64(x) * g.TileSize, float64(y) * g.TileSize
}

func (g *tilemap[any]) worldPosToGrid(x, y float64) (int, int) {
	return int(x / g.TileSize), int(y / g.TileSize)
}
