package graphic

import (
	"fmt"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/camera"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/transform2D"
)

type Graphic struct {
	*transform2D.Transform2D
	camera *camera.Camera
}

// A pretty big safety issue discovered: cmd.Scene() can be nil
// in OnTreeAdd...

// func (g *Graphic) OnTreeAdd(node *engine.Node, cmd *engine.Commands) {
// 	g.Transform2D.OnTreeAdd(node, cmd)
// 	camNode, ok := engine.GetNodeByType[*camera.Camera](cmd.Scene())
// 	if !ok {
// 		return
// 		// fmt.Println("Død og jøde")
// 	}
// 	fmt.Println(camNode)
// 	camActor, ok := engine.GetActor[*camera.Camera](camNode.GetValue())
// 	g.camera = camActor
// }

// This doesn't get called in Spawn()
func (g *Graphic) Init(cmd *engine.Commands) {
	g.Transform2D.Init(cmd)
	camNode, ok := engine.GetNodeByType[*camera.Camera](cmd.Scene())
	if !ok {
		fmt.Println("No camera was found! Instantiating default camera")
		camNode = cmd.Scene().SpawnActor("Camera", camera.NewCamera(
			transform2D.NewTransform2D(
				nodeactor.NewNode(),
			),
		), cmd)
	}
	camActor, ok := engine.GetActor[*camera.Camera](camNode.GetValue())
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
