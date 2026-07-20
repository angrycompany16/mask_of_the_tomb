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
	"mask_of_the_tomb/internal/engine/actors/UI/textbox"
	"mask_of_the_tomb/internal/engine/actors/UI/uigraphic"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
	"mask_of_the_tomb/internal/engine/actors/sound"
	"mask_of_the_tomb/internal/engine/commands"
)

func MakeMainMenuBundle() engine.Bundle {
	return func(cmd *commands.Commands, scene *engine.Scene) {
		gw, gh := cmd.Renderer.GetGameSize()
		ps := cmd.Renderer.GetPixelScale()

		rootAlign := scene.SpawnActor("RootAlign", align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
				container.WithRect(maths.NewRect(0, 0, gw*ps, gh*ps)),
			),
			align.WithIsRow(false),
			align.WithSpacing([]float64{2, 3}),
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

		titleActor := textbox.NewTextBox(
			container.NewContainer(
				nodeactor.NewNode(),
			),
			"fonts/JSE_AmigaAMOS.ttf",
			textbox.WithText("Meletus' tomb"),
		)

		rootAlign.AddChild(titleActor, "Title", engine.MakeOnTreeAdd(titleActor, cmd))

		selectList := selectlist.NewSelectList(
			align.NewAlign(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				align.WithIsRow(false),
				align.WithSpacing([]float64{1, 1, 1}),
			),
		)

		selectListNode := rootAlign.AddChild(selectList, "selectList", engine.MakeOnTreeAdd(selectList, cmd))

		playButton := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Play video game"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
				scenemanager, _ := commands.Get[engine.SceneManager](cmd)
				scenemanager.SpawnScene("d5ae6780-1030-11f0-996f-efbed2df7e2d", cmd)
			}),
		)

		selectListNode.AddChild(playButton, "playButton", engine.MakeOnTreeAdd(playButton, cmd))

		optionsButton := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Options"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
				scenemanager, _ := commands.Get[engine.SceneManager](cmd)
				scenemanager.SpawnScene("OptionsMenu", cmd)
			}),
		)

		selectListNode.AddChild(optionsButton, "optionsButton", engine.MakeOnTreeAdd(optionsButton, cmd))

		quitButton := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Don't play video game"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
				cmd.GameInfo.Exit = true
			}),
		)

		selectListNode.AddChild(quitButton, "quitButton", engine.MakeOnTreeAdd(quitButton, cmd))

		selectSound := sound.NewSoundPlayer(
			nodeactor.NewNode(),
			sound.WithSoundData("sfx/select3.ogg", false, "select"),
			sound.WithStartTriggers(selectList.OnSelectEv),
			//			sound.WithStartTriggers(playButton.OnHoverStart, optionsButton.OnHoverStart, quitButton.OnHoverStart, selectList.OnSelect),
		)

		rootAlign.AddChild(selectSound, "selectSound", engine.MakeOnTreeAdd(selectSound, cmd))
	}
}
