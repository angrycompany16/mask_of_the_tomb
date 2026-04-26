package bundles

import (
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/align"
	"mask_of_the_tomb/internal/engine/actors/UI/container"
	"mask_of_the_tomb/internal/engine/actors/UI/cursor"
	"mask_of_the_tomb/internal/engine/actors/UI/selectable"
	"mask_of_the_tomb/internal/engine/actors/UI/selectlist"
	"mask_of_the_tomb/internal/engine/actors/UI/slider"
	"mask_of_the_tomb/internal/engine/actors/UI/textbox"
	"mask_of_the_tomb/internal/engine/actors/UI/uigraphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/commands"
)

func MakeOptionsbundle() engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) {
		gw, gh := cmd.Renderer.GetGameSize()
		ps := cmd.Renderer.GetPixelScale()

		rootAlign := scene.SpawnActor("RootAlign", align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
				container.WithRect(maths.NewRect(0, 0, gw*ps, gh*ps)),
			),
			align.WithIsRow(false),
			align.WithSpacing([]float64{1}),
		), cmd)

		scene.SpawnActor("CoolCursor", cursor.NewCursor(
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
		), cmd)

		buttonAlign := scene.AddChild("Align",
			selectlist.NewSelectList(
				align.NewAlign(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					align.WithIsRow(false),
					align.WithSpacing([]float64{1, 1, 1, 1, 1, 1}),
				),
			), rootAlign, cmd)

		masterAlign := scene.AddChild("MasterAlign", align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
			),
		), buttonAlign, cmd)

		scene.AddChild("Master",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Master volume"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
				}),
			), masterAlign, cmd)

		scene.AddChild("MasterSlider",
			slider.NewSlider(
				selectable.NewSelectable(
					textbox.NewTextBox(
						container.NewContainer(
							nodeactor.NewNode(),
						),

						"fonts/JSE_AmigaAMOS.ttf",
					),
					selectable.WithCallback(func(cmd *commands.Commands) {
					}),
				), 0, 100, 10, 50,
			), masterAlign, cmd)

		musicAlign := scene.AddChild("MusicAlign", align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
			),
		), buttonAlign, cmd)

		scene.AddChild("Music",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Music volume"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
				}),
			), musicAlign, cmd)

		scene.AddChild("MusicSlider",
			slider.NewSlider(
				selectable.NewSelectable(
					textbox.NewTextBox(
						container.NewContainer(
							nodeactor.NewNode(),
						),
						"fonts/JSE_AmigaAMOS.ttf",
					),
					selectable.WithCallback(func(cmd *commands.Commands) {
					}),
				), 0, 100, 10, 50,
			), musicAlign, cmd)

		sfxAlign := scene.AddChild("SfxAlign", align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
			),
		), buttonAlign, cmd)

		scene.AddChild("Sfx",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Sfx volume"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
				}),
			), sfxAlign, cmd)

		scene.AddChild("SfxSlider",
			slider.NewSlider(
				selectable.NewSelectable(
					textbox.NewTextBox(
						container.NewContainer(
							nodeactor.NewNode(),
						),
						"fonts/JSE_AmigaAMOS.ttf",
					),
					selectable.WithCallback(func(cmd *commands.Commands) {
					}),
				), 0, 100, 10, 50,
			), sfxAlign, cmd)

		qualityNode := scene.AddChild("Quality",
			align.NewAlign(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				align.WithSpacing([]float64{1, 1}),
			), buttonAlign, cmd)

		scene.AddChild("QualityText",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Quality"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
				}),
			), qualityNode, cmd)

		qualitySelectable := scene.AddChild("QualitySelectable",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText(""),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
				}),
			), qualityNode, cmd)

		qualityOpAlign := scene.AddChild("QualityOption",
			selectlist.NewSelectList(
				align.NewAlign(
					container.NewContainer(
						nodeactor.NewNode(),
						container.WithAutoAlign(container.Fill),
					),
					align.WithIsRow(true),
					align.WithSpacing([]float64{1, 1, 1}),
				),
			), qualitySelectable, cmd)

		scene.AddChild("QualityLow",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Low"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
				}),
			), qualityOpAlign, cmd)

		scene.AddChild("QualityMid",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Mid"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
				}),
			), qualityOpAlign, cmd)

		scene.AddChild("QualityHigh",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("High"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
				}),
			), qualityOpAlign, cmd)

		scene.AddChild("Back",
			selectable.NewSelectable(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Back to main menu"),
				),
				selectable.WithCallback(func(cmd *commands.Commands) {
					// Hmmmm... We don't want this to happen if we spawn
					// this as a gameplay pause screen
					scenemanager, _ := commands.Get[engine.SceneManager](cmd)
					scenemanager.SpawnScene("MainMenu", cmd)
				}),
			), buttonAlign, cmd)
	}
}
