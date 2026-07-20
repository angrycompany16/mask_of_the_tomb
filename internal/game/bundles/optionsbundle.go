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
	"mask_of_the_tomb/internal/engine/actors/sound"
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

		buttonAlign := selectlist.NewSelectList(
			align.NewAlign(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				align.WithIsRow(false),
				align.WithSpacing([]float64{1, 1, 1, 1, 1, 1}),
			),
		)

		buttonAlignNode := rootAlign.AddChild(buttonAlign, "buttonAlign", engine.MakeOnTreeAdd(buttonAlign, cmd))

		masterAlign := align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
			),
		)

		masterAlignNode := buttonAlignNode.AddChild(masterAlign, "masterAlign", engine.MakeOnTreeAdd(masterAlign, cmd))

		masterVolume := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Master volume"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
			}),
		)

		masterAlignNode.AddChild(masterVolume, "masterVolume", engine.MakeOnTreeAdd(masterVolume, cmd))

		masterSlider := slider.NewSlider(
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
		)

		masterAlignNode.AddChild(masterSlider, "masterSlider", engine.MakeOnTreeAdd(masterSlider, cmd))

		musicAlign := align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
			),
		)

		musicAlignNode := buttonAlignNode.AddChild(musicAlign, "musicAlign", engine.MakeOnTreeAdd(musicAlign, cmd))

		musicVolume := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Music volume"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
			}),
		)

		musicAlignNode.AddChild(musicVolume, "musicVolume", engine.MakeOnTreeAdd(musicVolume, cmd))

		musicSlider := slider.NewSlider(
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
		)

		musicAlignNode.AddChild(musicSlider, "musicSlider", engine.MakeOnTreeAdd(musicSlider, cmd))

		sfxAlign := align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
			),
		)

		sfxAlignNode := buttonAlignNode.AddChild(sfxAlign, "sfxAlign", engine.MakeOnTreeAdd(sfxAlign, cmd))

		sfxVolume := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Sfx volume"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
			}),
		)

		sfxAlignNode.AddChild(sfxVolume, "sfxAlign", engine.MakeOnTreeAdd(sfxVolume, cmd))

		sfxSlider := slider.NewSlider(
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
		)

		sfxAlignNode.AddChild(sfxSlider, "sfxSlider", engine.MakeOnTreeAdd(sfxSlider, cmd))

		qualityAlign := align.NewAlign(
			container.NewContainer(
				nodeactor.NewNode(),
			),
			align.WithSpacing([]float64{1, 1}),
		)

		qualityAlignNode := buttonAlignNode.AddChild(qualityAlign, "qualityAlign", engine.MakeOnTreeAdd(qualityAlign, cmd))

		qualityText := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Quality"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
			}),
		)

		qualityAlignNode.AddChild(qualityText, "qualityText", engine.MakeOnTreeAdd(qualityText, cmd))

		qualitySelectable := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText(""),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
			}),
		)

		qualitySelectableNode := qualityAlignNode.AddChild(qualitySelectable, "qualitySelectable", engine.MakeOnTreeAdd(qualitySelectable, cmd))

		qualityOption := selectlist.NewSelectList(
			align.NewAlign(
				container.NewContainer(
					nodeactor.NewNode(),
					container.WithAutoAlign(container.Fill),
				),
				align.WithIsRow(true),
				align.WithSpacing([]float64{1, 1, 1}),
			),
		)

		qualityOptionNode := qualitySelectableNode.AddChild(qualityOption, "qualityOption", engine.MakeOnTreeAdd(qualityOption, cmd))

		qualityLow := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Low"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
			}),
		)

		qualityOptionNode.AddChild(qualityLow, "qualityLow", engine.MakeOnTreeAdd(qualityLow, cmd))

		qualityMid := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("Mid"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
			}),
		)

		qualityOptionNode.AddChild(qualityMid, "qualityMid", engine.MakeOnTreeAdd(qualityMid, cmd))

		qualityHigh := selectable.NewSelectable(
			textbox.NewTextBox(
				container.NewContainer(
					nodeactor.NewNode(),
				),
				"fonts/JSE_AmigaAMOS.ttf",
				textbox.WithText("High"),
			),
			selectable.WithCallback(func(cmd *commands.Commands) {
			}),
		)

		qualityOptionNode.AddChild(qualityHigh, "qualityHigh", engine.MakeOnTreeAdd(qualityHigh, cmd))

		backButton := selectable.NewSelectable(
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
		)

		buttonAlignNode.AddChild(backButton, "backButton", engine.MakeOnTreeAdd(backButton, cmd))

		selectSound := sound.NewSoundPlayer(
			nodeactor.NewNode(),
			sound.WithSoundData("sfx/select3.ogg", false, "select"),
			sound.WithStartTriggers(buttonAlign.OnSelectEv, qualityOption.OnSelectEv, masterSlider.OnChangeEv, musicSlider.OnChangeEv, sfxSlider.OnChangeEv),
		)

		rootAlign.AddChild(selectSound, "selectSound", engine.MakeOnTreeAdd(selectSound, cmd))

	}
}
