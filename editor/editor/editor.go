package editor

import (
	"errors"
	"fmt"
	"mask_of_the_tomb/editor/fileio"
	ui "mask_of_the_tomb/internal/game/UI"
	"mask_of_the_tomb/internal/game/core/assetloader"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type editorState int

const (
	Loading editorState = iota
	Preview
	OpeningFile
)

var (
	ErrTerminated      = errors.New("Terminated")
	uiPath             = filepath.Join("assets", "menus", "editor", "mainUI.yaml")
	loadingScreenPath  = filepath.Join("assets", "menus", "editor", "loadingscreen.yaml")
	openFileScreenPath = filepath.Join("assets", "menus", "editor", "openfile.yaml")
	loadFinishedChan   = make(chan int)
)

type Editor struct {
	editorUI *ui.UI
	state    editorState
}

func (e *Editor) Init() {
	e.editorUI.LoadPreamble(loadingScreenPath)
	e.editorUI.Load(uiPath, openFileScreenPath)

	go assetloader.LoadAll(loadFinishedChan)
}

func (e *Editor) Update() error {
	e.editorUI.Update()
	submits := e.editorUI.GetSubmits()
	switch e.state {
	case Loading:
		select {
		case <-loadFinishedChan:
			fmt.Println("Finished loading")
			e.state = Preview
			e.editorUI.SwitchActiveMenu("mainUI")
		default:
		}
	case Preview:
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyS) {
			fmt.Println("Save the asset file")
		}
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyO) {
			e.state = OpeningFile
			e.editorUI.SwitchActiveMenu("openfile")
			fmt.Println("Open file manager to look for file")
		}
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
			return ErrTerminated
		}
	case OpeningFile:
		if submits["path"] != "" {
			fmt.Println("Searching for file", submits["path"])
			files := make([]string, 0)
			fileio.FindFiles(submits["path"], files)
			// TODO: Dynamically create selectables
			e.editorUI.SwitchActiveMenu("mainUI")
			e.state = Preview
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			e.editorUI.SwitchActiveMenu("mainUI")
			e.state = Preview
		}
	}
	// How to make save/load system for assets?
	return nil
}

func (e *Editor) Draw() {
	e.editorUI.Draw()
	// Draw menu and active asset
}

func NewEditor() *Editor {
	return &Editor{
		editorUI: ui.NewUI(),
		state:    Loading,
	}
}
