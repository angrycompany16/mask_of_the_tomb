package engine

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/renderer"
)

type ServerArgs struct {
	GameWidth, GameHeight int
	PixelScale            int
}

// This is quickly turning into some sort of ctx..
type Servers struct {
	renderer    *renderer.Renderer
	assetloader *assetloader.AssetLoader
	scene       *Scene
	inputServer *input.InputServer
}

func (s *Servers) Renderer() *renderer.Renderer {
	return s.renderer
}

func (s *Servers) AssetLoader() *assetloader.AssetLoader {
	return s.assetloader
}

func (s *Servers) Scene() *Scene {
	return s.scene
}

func (s *Servers) InputHandler() *input.InputServer {
	return s.inputServer
}

func NewServers(args ServerArgs) *Servers {
	return &Servers{
		renderer:    renderer.NewRenderer(args.GameWidth, args.GameHeight, args.PixelScale),
		assetloader: assetloader.NewAssetLoader(&assets.FS),
		inputServer: input.NewInputServer(),
	}
}
