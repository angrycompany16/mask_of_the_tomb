package servers

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/engine/servers/assetloader"
	"mask_of_the_tomb/internal/engine/servers/globals"
	"mask_of_the_tomb/internal/engine/servers/renderer"
)

type ServerArgs struct {
	GameWidth, GameHeight int
	PixelScale            int
}

type Servers struct {
	renderer    *renderer.Renderer
	globals     *globals.Globals
	assetloader *assetloader.AssetLoader
}

func (s *Servers) Renderer() *renderer.Renderer {
	return s.renderer
}

func (s *Servers) Globals() *globals.Globals {
	return s.globals
}

func (s *Servers) AssetLoader() *assetloader.AssetLoader {
	return s.assetloader
}

func NewServers(args ServerArgs) *Servers {
	return &Servers{
		renderer:    renderer.NewRenderer(args.GameWidth, args.GameHeight, args.PixelScale),
		globals:     globals.NewGlobals(),
		assetloader: assetloader.NewAssetLoader(&assets.FS),
	}
}
