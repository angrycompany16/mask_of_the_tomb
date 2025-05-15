package assetloader

var (
	_assetLoader = assetLoader{
		assetPool: make(map[string]Asset),
	}
)

// Very simple asset loader interface
type Asset interface {
	Load()
}

type assetLoader struct {
	assetPool map[string]Asset
}

func Exists(hash string) (Asset, bool) {
	asset, ok := _assetLoader.assetPool[hash]
	return asset, ok
}

func Load(hash string, asset Asset) {
	_assetLoader.assetPool[hash] = asset
}

func LoadAll(doneChan chan<- int) {
	for _, asset := range _assetLoader.assetPool {
		asset.Load()
	}
	doneChan <- 1
}
