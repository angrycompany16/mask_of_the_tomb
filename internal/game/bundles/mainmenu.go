package bundles

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/colors"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/UI/align"
	"mask_of_the_tomb/internal/engine/actors/UI/button"
	"mask_of_the_tomb/internal/engine/actors/UI/buttonalign"
	"mask_of_the_tomb/internal/engine/actors/UI/container"
	"mask_of_the_tomb/internal/engine/actors/UI/textbox"
	"mask_of_the_tomb/internal/engine/actors/nodeactor"
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

		scene.AddChild("Title", textbox.NewTextBox(
			container.NewContainer(
				nodeactor.NewNode(),
			),
			"fonts/JSE_AmigaAMOS.ttf",
			textbox.WithText("Meletus' tomb"),
		), rootAlign, cmd)

		buttonAlign := scene.AddChild("Align",
			buttonalign.NewButtonAlign(
				align.NewAlign(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					align.WithIsRow(false),
					align.WithSpacing([]float64{1, 1, 1}),
				),
			), rootAlign, cmd)

		scene.AddChild("Text1",
			button.NewButton(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Play video game"),
				),
				colors.ColorPair{
					BrightColor: color.RGBA{255, 255, 255, 255},
					DarkColor:   color.RGBA{100, 100, 100, 255},
				},
				colors.ColorPair{
					BrightColor: color.RGBA{255, 0, 0, 255},
					DarkColor:   color.RGBA{150, 50, 50, 255},
				},
			), buttonAlign, cmd)

		scene.AddChild("Text2",
			button.NewButton(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Options"),
				),
				colors.ColorPair{
					BrightColor: color.RGBA{255, 255, 255, 255},
					DarkColor:   color.RGBA{100, 100, 100, 255},
				},
				colors.ColorPair{
					BrightColor: color.RGBA{255, 0, 0, 255},
					DarkColor:   color.RGBA{150, 50, 50, 255},
				},
			), buttonAlign, cmd)

		scene.AddChild("Text3",
			button.NewButton(
				textbox.NewTextBox(
					container.NewContainer(
						nodeactor.NewNode(),
					),
					"fonts/JSE_AmigaAMOS.ttf",
					textbox.WithText("Don't play video game"),
				),
				colors.ColorPair{
					BrightColor: color.RGBA{255, 255, 255, 255},
					DarkColor:   color.RGBA{100, 100, 100, 255},
				},
				colors.ColorPair{
					BrightColor: color.RGBA{255, 0, 0, 255},
					DarkColor:   color.RGBA{150, 50, 50, 255},
				},
			), buttonAlign, cmd)
	}
}
