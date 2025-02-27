package world

import (
	"mask_of_the_tomb/ebitenLDTK"
	ebitenrenderutil "mask_of_the_tomb/ebitenRenderUtil"
	"mask_of_the_tomb/files"
	"mask_of_the_tomb/game/animation"
	"mask_of_the_tomb/rendering"
)

type BreakableBlock struct {
	anim             *animation.Animation
	posX, posY       float64
	offsetX, offsetY float64
}

// Need:
// - Animation play/pause functionality
// - Dynamic level hitboxes

func (b *BreakableBlock) Update() {
	b.anim.Update()
	// If collision occurs
	// Start break timer
	// Update animation timer
}

func (b *BreakableBlock) Draw() {
	// fmt.Println("laskmdc")
	ebitenrenderutil.DrawAt(b.anim.GetSprite(), rendering.RenderLayers.Playerspace, b.posX+b.offsetX, b.posY+b.offsetY)
}

func NewBreakableBlock(entityInstance *ebitenLDTK.Entity) *BreakableBlock {
	return &BreakableBlock{
		anim: animation.NewAnimation(
			animation.NewSpritesheetAuto(files.LazyImage(crumbleSpritePath)),
			0.1,
			animation.Strip,
			animation.Once,
		),
		posX:    entityInstance.Px[0],
		posY:    entityInstance.Px[1],
		offsetX: -2,
		offsetY: -2,
	}
}
