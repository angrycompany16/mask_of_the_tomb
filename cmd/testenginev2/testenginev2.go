package main

import (
	"fmt"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	treeUI      debugui.DebugUI
	inspectorUI debugui.DebugUI
	scene       *engine.Scene
}

func (g *Game) Update() error {
	if _, err := g.treeUI.Update(
		g.scene.MakeDrawFunc(),
	); err != nil {
		return err
	}

	FPS := ebiten.ActualFPS()
	TPS := ebiten.ActualTPS()
	fmt.Printf("FPS: %f, TPS: %f\n", FPS, TPS)

	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.treeUI.Draw(screen)
	g.inspectorUI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	// Create the test scene and the root node
	// This looks better. Separation between node and behaviour is clearer.
	testScene := engine.NewScene("test (A scene!!!)", actors.NewNode())
	root := testScene.GetRoot()
	child1ID := testScene.AddChild("Child 1", actors.NewTransform2D(100, 100, 0, 1, 1), root)
	child1, _ := testScene.GetNodeByID(child1ID)
	testScene.AddChild("Child 2", actors.NewTransform2D(100, 100, 0, 1, 1), child1)

	game := &Game{
		scene: testScene,
	}

	ebiten.SetWindowSize(960, 540)
	testScene.Init()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
