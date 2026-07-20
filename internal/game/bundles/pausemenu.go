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
	"mask_of_the_tomb/internal/engine/actors/sound"
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

		rootAlign := align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
				container.WithRect(maths.NewRect(0, 0, gw*ps, gh*ps)),
			),
			align.WithIsRow(false),
			align.WithSpacing([]float64{2, 3}),
		)

		rootAlignNode := rootContainer.AddChild(rootAlign, "rootAlign", engine.MakeOnTreeAdd(rootAlign, cmd))

		cursor := cursor.NewCursor(
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
		)

		rootContainer.AddChild(cursor, "cursor", engine.MakeOnTreeAdd(cursor, cmd))

		title := textbox.NewTextBox(
			container.NewContainer(
				nodeactor.NewNode(),
			),
			"fonts/JSE_AmigaAMOS.ttf",
			textbox.WithText("Pause"),
		)

		rootAlignNode.AddChild(title, "title", engine.MakeOnTreeAdd(title, cmd))

		buttonAlign := selectlist.NewSelectList(
			align.NewAlign(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				align.WithIsRow(false),
				align.WithSpacing([]float64{1, 1, 1}),
			),
		)

		buttonAlignNode := rootAlignNode.AddChild(buttonAlign, "buttonAlign", engine.MakeOnTreeAdd(buttonAlign, cmd))

		playButton := selectable.NewSelectable(
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
		)

		buttonAlignNode.AddChild(playButton, "playButton", engine.MakeOnTreeAdd(playButton, cmd))

		optionsButton := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Options"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {

			}),
		)

		buttonAlignNode.AddChild(optionsButton, "optionsButton", engine.MakeOnTreeAdd(optionsButton, cmd))

		quitButton := selectable.NewSelectable(
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
		)

		buttonAlignNode.AddChild(quitButton, "quitButton", engine.MakeOnTreeAdd(quitButton, cmd))

		selectSound := sound.NewSoundPlayer(
			nodeactor.NewNode(),
			sound.WithSoundData("sfx/select3.ogg", false, "select"),
			sound.WithStartTriggers(buttonAlign.OnSelectEv),
		)

		rootAlignNode.AddChild(selectSound, "selectSound", engine.MakeOnTreeAdd(selectSound, cmd))

		return rootContainer
	}
}
