package editor

import (
	"errors"
	"fmt"
	"mask_of_the_tomb/editor/fileio"
	ui "mask_of_the_tomb/internal/game/UI"
	"mask_of_the_tomb/internal/game/core/assetloader"
	"mask_of_the_tomb/internal/game/core/assetloader/delayasset"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/core/rendering/camera"
	"mask_of_the_tomb/internal/game/physics/particles"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TODO (big):
// Create a system for the UI. Right now it's just random pieces of code mingled
// together and king of hardcoded, but there's a glimpse of structure below all that.
// Use events and all that to make the system work and behave very nicely.
// Remember that just because it's oop it doesn't have to be bad code

// Although we will experiment with the actor model engine

const (
	defaultSpawnX, defaultSpawnY = 100, 100
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
	editorUI             *ui.UI
	state                editorState
	activeParticleSystem *particles.ParticleSystem
}

func (e *Editor) Init() {
	e.editorUI.LoadPreamble(loadingScreenPath)
	e.editorUI.Load(uiPath, openFileScreenPath)

	delayAsset := delayasset.NewDelayAsset(time.Second)
	assetloader.AddAsset(&delayAsset)

	go assetloader.LoadAll(loadFinishedChan)
}

// TODO: Allow the editor to change parameters using sliders and such
// Need to:
// - Display all the parameters
// - Write to file
// - For now, let's just use input boxes
// - We need some kind of validation system though
func (e *Editor) Update() error {
	e.editorUI.Update()
	switch e.state {
	case Loading:
		select {
		case <-loadFinishedChan:
			fmt.Println("Finished loading")
			e.state = Preview
			e.editorUI.SwitchActiveDisplay("mainUI")

			camera.Init(
				rendering.GameWidth,
				rendering.GameHeight,
				rendering.GameWidth/2,
				rendering.GameHeight/2,
			)
		default:
		}
	case Preview:
		if e.activeParticleSystem != nil {
			e.activeParticleSystem.Update()
		}

		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyS) {
			fmt.Println("Save the asset file")
		}
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyO) {
			e.state = OpeningFile
			e.editorUI.SwitchActiveDisplay("openfile")
			fmt.Println("Open file manager to look for file")
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			if e.activeParticleSystem != nil {
				e.activeParticleSystem.Play()
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
			return ErrTerminated
		}
	case OpeningFile:
		// TODO: Reconfigure with events
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			_type, asset := fileio.OpenAsset(e.editorUI.GetFileSearchValue())
			switch _type {
			case "ParticleSystem":
				ps := asset.(*particles.ParticleSystem)
				e.activeParticleSystem = ps
				e.activeParticleSystem.Play()
				e.editorUI.ResetFileSearch()
				ps.PosX = defaultSpawnX
				ps.PosY = defaultSpawnY
			}
			e.editorUI.SwitchActiveDisplay("mainUI")
			e.state = Preview
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			e.editorUI.SwitchActiveDisplay("mainUI")
			e.state = Preview
		}
	}
	// How to make save/load system for assets?
	return nil
}

func (e *Editor) Draw() {
	e.editorUI.Draw()
	// Draw menu and active asset
	if e.state == Preview && e.activeParticleSystem != nil {
		e.activeParticleSystem.Draw()
	}
}

func NewEditor() *Editor {
	return &Editor{
		editorUI: ui.NewUI(),
		state:    Loading,
	}
}
