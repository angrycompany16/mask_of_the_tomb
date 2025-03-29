package editor

import (
	"errors"
	"fmt"
	ui "mask_of_the_tomb/internal/game/UI"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	ErrTerminated = errors.New("Terminated")
)

type Editor struct {
	editorUI *ui.UI
}

func (e *Editor) Init() {

}

func (e *Editor) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyControl) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		fmt.Println("Save the asset file")
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) || inpututil.IsKeyJustPressed(ebiten.KeyO) {
		fmt.Println("Open file manager to look for file")
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) || inpututil.IsKeyJustPressed(ebiten.KeyC) {
		return ErrTerminated
	}
	// How to make save/load system for assets?
	return nil
}

func (e *Editor) Draw() {
	// Draw menu and active asset
}

func NewEditor() *Editor {
	return &Editor{
		editorUI: ui.NewUI(),
	}
}
