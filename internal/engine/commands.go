package engine

import (
	"io/fs"
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/backend/triggerenv"
)

// TODO: I'd like to change this so that it's customisable per game.
// This would entail:
// - Turn Commands into an interface with methods for renderer, assetloader, etc. (but only basic stuff)
// - Create a default implementation
// - Enable extension via struct embedding
type Commands struct {
	renderer    *renderer.Renderer
	assetloader *assetloader.AssetLoader
	scene       *Scene
	inputServer *input.InputServer
	slamboxEnv  *slambox.SlamboxEnvironment
	triggerEnv  *triggerenv.TriggerEnv
}

func (c *Commands) Renderer() *renderer.Renderer {
	return c.renderer
}

func (c *Commands) AssetLoader() *assetloader.AssetLoader {
	return c.assetloader
}

func (c *Commands) Scene() *Scene {
	return c.scene
}

func (c *Commands) InputHandler() *input.InputServer {
	return c.inputServer
}

func (c *Commands) SlamboxEnv() *slambox.SlamboxEnvironment {
	return c.slamboxEnv
}

func (c *Commands) TriggerEnv() *triggerenv.TriggerEnv {
	return c.triggerEnv
}

func NewCommands(options ...Option) *Commands {
	servers := defaultCommands()

	for _, option := range options {
		option(servers)
	}

	return servers
}

func defaultCommands() *Commands {
	return &Commands{
		renderer:    renderer.NewRenderer(480, 270, 4, true, true),
		assetloader: assetloader.NewAssetLoader(&assets.FS),
		inputServer: input.NewInputServer(),
		slamboxEnv:  slambox.NewSlamboxEnvironment(8),
		triggerEnv:  triggerenv.NewTriggerEnv(),
	}
}

type Option func(*Commands)

func WithRenderer(width, height, pixelScale int, fullScreen, hideCursor bool) Option {
	return func(s *Commands) {
		s.renderer = renderer.NewRenderer(width, height, pixelScale, fullScreen, hideCursor)
	}
}

func WithCamera(tileSize int) Option {
	return func(s *Commands) {
		s.slamboxEnv = slambox.NewSlamboxEnvironment(tileSize)
	}
}

func WithAssetLoader(fs fs.FS) Option {
	return func(s *Commands) {
		s.assetloader = assetloader.NewAssetLoader(fs)
	}
}

func WithSlamboxEnvironment(tileSize int) Option {
	return func(s *Commands) {
		s.slamboxEnv = slambox.NewSlamboxEnvironment(tileSize)
	}
}
