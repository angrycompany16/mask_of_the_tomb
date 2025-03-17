package app

import (
	"mask_of_the_tomb/examples/tree/testentity"
	"mask_of_the_tomb/internal/engine/entities"
	"mask_of_the_tomb/internal/libraries/errs"
	"mask_of_the_tomb/internal/libraries/rendering"
	"mask_of_the_tomb/internal/libraries/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type App struct{}

func (a *App) Init() {
	var prev *entities.Entity
	for i := range 40 {
		utils.UNUSED(i)
		testEntity := testentity.NewTestEntity(400)
		entity := entities.RegisterEntity(testEntity, "Test entity")
		testEntity.Entity = entity
		if prev != nil {
			entity.SetParent(prev)
		}

		prev = entity
	}
	specialEntity := testentity.NewSpecialEntity(401)
	entity := entities.RegisterEntity(specialEntity, "Special entity")
	specialEntity.Entity = entity
	specialEntity.SetParent(prev)
}

func (a *App) Update() error {
	entities.PreUpdate()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errs.ErrTerminated
	}

	entities.Update()
	entities.PostUpdate()

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	rendering.RenderLayers.Draw(screen)
}

func (a *App) Layout(outsideHeight, outsideWidth int) (int, int) {
	return rendering.GameWidth * rendering.PixelScale, rendering.GameHeight * rendering.PixelScale
}

func MakeApp() *App {
	return &App{}
}
