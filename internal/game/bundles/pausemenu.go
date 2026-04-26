package bundles

import (
	"fmt"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/backend/triggerenv"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/align"
	"mask_of_the_tomb/internal/engine/actors/UI/container"
	"mask_of_the_tomb/internal/engine/actors/UI/cursor"
	"mask_of_the_tomb/internal/engine/actors/UI/selectable"
	"mask_of_the_tomb/internal/engine/actors/UI/selectlist"
	"mask_of_the_tomb/internal/engine/actors/UI/textbox"
	"mask_of_the_tomb/internal/engine/actors/UI/uigraphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
)

func MakePauseMenuBundle() engine.BundleV2 {
	return func(cmd *commands.Commands, scene *engine.Scene) *engine.Node {
		gw, gh := cmd.Renderer.GetGameSize()
		ps := cmd.Renderer.GetPixelScale()

		rootContainer := scene.SpawnActor("PauseScreenBundle", container.NewContainer(
			nodeactor.NewNode(),
			container.WithRect(maths.NewRect(0, 0, gw*ps, gh*ps)),
		), cmd)

		rootAlign := scene.AddChild("Align", align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
				container.WithRect(maths.NewRect(0, 0, gw*ps, gh*ps)),
			),
			align.WithIsRow(false),
			align.WithSpacing([]float64{2, 3}),
		), rootContainer, cmd)

		scene.AddChild("CoolCursor", cursor.NewCursor(
			uigraphic.NewUIGraphic(
				container.NewContainer(
					nodeactor.NewNode(),
					container.WithRect(maths.NewRect(0, 0, 15, 15)),
				),
				"sprites/icons/circle-15x15.png",
				3,
				renderer.RenderTarget{
					Type: renderer.SCREEN,
					Name: "ScreenUI",
				},
			),
		), rootContainer, cmd)

		scene.AddChild("Title", textbox.NewTextBox(
			container.NewContainer(
				nodeactor.NewNode(),
			),
			"fonts/JSE_AmigaAMOS.ttf",
			textbox.WithText("Pause"),
		), rootAlign, cmd)

		buttonAlign := scene.AddChild("Align",
			selectlist.NewSelectList(
				align.NewAlign(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					align.WithIsRow(false),
					align.WithSpacing([]float64{1, 1, 1}),
				),
			), rootAlign, cmd)

		scene.AddChild("Text1",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Play"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
					playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
					playerControls.Active = true

					scene, _ := commands.Get[engine.Scene](cmd)
					pauseScreenRoot, ok := scene.GetNodeByName("PauseScreenBundle")
					if !ok {
						fmt.Println("Bad pasuemaneuwspawner update callback")
						return
					}
					scene.Delete(pauseScreenRoot)
				}),
			), buttonAlign, cmd)

		scene.AddChild("Text2",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Options"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {

				}),
			), buttonAlign, cmd)

		scene.AddChild("Text3",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Quit"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
					slamboxenv, _ := commands.Get[slambox.SlamboxEnvironment](cmd)
					slamboxenv.Reset()
					triggerenv, _ := commands.Get[triggerenv.TriggerEnv](cmd)
					triggerenv.Reset()

					playerControls := cmd.InputHandler.InputSchemes["PlayerControls"]
					playerControls.Active = true
					scenemanager, _ := commands.Get[engine.SceneManager](cmd)
					scenemanager.SpawnScene("MainMenu", cmd)
				}),
			), buttonAlign, cmd)

		return rootContainer
	}
}
