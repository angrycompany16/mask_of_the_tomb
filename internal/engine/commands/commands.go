package commands

import (
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/input"
	"mask_of_the_tomb/internal/backend/renderer"
	"reflect"
)

type Commands struct {
	Renderer     *renderer.Renderer
	AssetLoader  *assetloader.AssetLoader
	InputHandler *input.InputHandler

	Globals map[reflect.Type]interface{}
}

func Set[T any](cmd *Commands, value *T) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	cmd.Globals[t] = value
}

func Get[T any](cmd *Commands) (*T, bool) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	value, ok := cmd.Globals[t]
	return value.(*T), ok
}

func NewCommands(renderer *renderer.Renderer, assetloader *assetloader.AssetLoader, inputhandler *input.InputHandler) *Commands {
	return &Commands{
		Renderer:     renderer,
		AssetLoader:  assetloader,
		InputHandler: inputhandler,
		Globals:      make(map[reflect.Type]interface{}),
	}
}
