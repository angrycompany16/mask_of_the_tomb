package ebitenLDTK

type LayerType string

const (
	LayerTypeTiles    = "Tiles"
	LayerTypeEntities = "Entities"
	LayerTypeIntGrid  = "IntGrid"
)

type Layer struct {
	Name           string    `json:"__identifier"`
	Type           LayerType `json:"__type"`
	GridSize       float64   `json:"__gridSize"`
	TilesetUid     int       `json:"__tilesetDefUid"`
	TilesetRelPath string    `json:"__tilesetRelPath"`
	PxOffsetX      int       `json:"pxOffsetX"`
	PxOffsetY      int       `json:"pxOffsetY"`
	GridTiles      []Tile    `json:"gridTiles"`
	Entities       []Entity  `json:"entityInstances"`
	AutoLayerTiles []Tile    `json:"autoLayerTiles"`
}

type Tile struct {
	Px              []float64   `json:"px"`
	Src             []float64   `json:"src"`
	TileOrientation Orientation `json:"f"`
	T               int         `json:"t"`
	D               []int       `json:"d"`
	A               float64     `json:"a"`
}

type Orientation int

const (
	OrientationNone = iota
	OrientationFlipX
	OrientationFlipY
	OrientationFlipXY
)
