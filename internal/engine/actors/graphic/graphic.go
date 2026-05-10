package graphic

import (
	"fmt"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/camera"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
	"mask_of_the_tomb/internal/engine/commands"
)

// TODO: Remove. This is most likely not needed at all

// Although in that case a question still remains: How to implement camera
// shake without Graphic?
type Graphic struct {
	*transform2D.Transform2D
	camera *camera.Camera
}

// This doesn't get called in Spawn()
func (g *Graphic) Init(cmd *commands.Commands) {
	g.Transform2D.Init(cmd)
	scene, ok := commands.Get[engine.Scene](cmd)
	if !ok {
		panic("Scene is missing from commands")
	}
	camNode, ok := engine.GetNodeByType[*camera.Camera](scene)
	if !ok {
		fmt.Println("No camera was found! Instantiating default camera")

		camNode = scene.SpawnActor("Camera", camera.NewCamera(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
		), cmd)
	}
	camActor, ok := engine.As[*camera.Camera](camNode.GetValue())
	g.camera = camActor
}

func (g *Graphic) GetCamera() *camera.Camera {
	return g.camera
}

func NewGraphic(transform2d *transform2D.Transform2D) *Graphic {
	return &Graphic{
		Transform2D: transform2d,
	}
}
