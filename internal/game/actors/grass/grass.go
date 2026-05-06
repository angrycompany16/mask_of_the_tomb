package grass

import (
	"fmt"
	"image"
	"mask_of_the_tomb/internal/backend/assetloader"
	"mask_of_the_tomb/internal/backend/assetloader/assettypes"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/renderer"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/actors/graphic"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/player"
	"mask_of_the_tomb/internal/utils"

	"math"
	"math/rand/v2"

	ebitenLDTK "github.com/angrycompany16/ebiten-LDTK"
	"github.com/aquilax/go-perlin"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	bladesPertile = 4.0
	imagePadding  = 16.0
)

type Grass struct {
	*graphic.Graphic
	image            *ebiten.Image
	blades           []*grassBlade
	perlin           *perlin.Perlin
	player           *player.Player
	grassTilemapPath string
	direction        maths.Direction
	grassTilemap     *assetloader.AssetRef[ebiten.Image]
	target           renderer.RenderTarget
	drawOrder        int
}

func (g *Grass) OnTreeAdd(node *engine.Node, cmd *commands.Commands) {
	g.Graphic.OnTreeAdd(node, cmd)
	g.grassTilemap = assetloader.StageAsset[ebiten.Image](
		cmd.AssetLoader,
		g.grassTilemapPath,
		assettypes.NewImageAsset(g.grassTilemapPath),
	)
}

func (g *Grass) Init(cmd *commands.Commands) {
	g.Graphic.Init(cmd)
	for _, blade := range g.blades {
		blade.init(g.grassTilemap.Value())
	}
}

// TODO: Optimize with some sort of batched update logic.
func (g *Grass) Update(cmd *commands.Commands) {
	g.Graphic.Update(cmd)

	// TODO: Don't do this. IT's dumb
	if g.player == nil {
		scene, _ := commands.Get[engine.Scene](cmd)
		playerNode, ok := scene.GetNodeByName("Player")
		if !ok {
			fmt.Println("Death")
			return
		}
		player, ok := engine.As[*player.Player](playerNode.GetValue())
		if !ok {
			fmt.Println("Grusom død")
			return
		}
		g.player = player
	}

	g.image.Clear()

	gPosX, gPosY := g.Transform2D.GetPos(false)
	gAngle := g.Transform2D.GetAngle(false)
	gScaleX, gScaleY := g.Transform2D.GetScale(false)
	camX, camY := g.GetCamera().WorldToCam(gPosX, gPosY, true)

	playerX, playerY := g.player.GetCenterPos()
	dirX, dirY := g.player.GetMovedir()
	playerVelX, playerVelY := dirX*g.player.MoveSpeed, dirY*g.player.MoveSpeed

	grassSpacePlayerX := playerX - gPosX + imagePadding
	grassSpacePlayerY := playerY - gPosY + imagePadding

	for _, blade := range g.blades {
		blade.Update(grassSpacePlayerX, grassSpacePlayerY, playerVelX, playerVelY, g.perlin, cmd.GameInfo.GetTime())
		bladeOp := opgen.PosRot(blade.sprite, blade.posX, blade.posY, blade.angle+blade.angleOffset+maths.DirToRadians(g.direction), 0.5, 1.0)
		g.image.DrawImage(blade.sprite, bladeOp)
	}

	cmd.Renderer.Request(opgen.PosRotScale(
		g.image,
		camX-imagePadding, camY-imagePadding,
		gAngle,
		gScaleX, gScaleY,
		0, 0,
	), g.image, g.target, g.drawOrder)
}

func NewGrass(
	graphic *graphic.Graphic,
	entity *ebitenLDTK.Entity,
	gridSize float64,
	grassTilemapPath string,
	grassWindSeed int64,
	target renderer.RenderTarget,
	drawOrder int,
) *Grass {
	directionField := utils.Must(entity.GetFieldByName("Direction"))
	direction := maths.DirFromString(ebitenLDTK.As[ebitenLDTK.Enum](directionField).Value)

	newGrass := Grass{
		Graphic:          graphic,
		image:            ebiten.NewImage(int(entity.Width+imagePadding*2), int(entity.Height+imagePadding*2)),
		perlin:           perlin.NewPerlin(2.0, 0.05, 4, grassWindSeed),
		blades:           make([]*grassBlade, 0),
		grassTilemapPath: grassTilemapPath,
		target:           target,
		drawOrder:        drawOrder,
		direction:        direction,
	}

	dirX, dirY := maths.VecFromDir(direction)
	interpX, interpY := math.Abs(dirX), math.Abs(dirY)

	t := 0.0
	length := 0.0
	switch direction {
	case maths.DirUp:
		t = imagePadding + entity.Height
		length = entity.Width
	case maths.DirDown:
		t = imagePadding
		length = entity.Width
	case maths.DirRight:
		t = imagePadding
		length = entity.Height
	case maths.DirLeft:
		t = imagePadding + entity.Width
		length = entity.Height
	}

	for i := 0; i < int(bladesPertile*length/gridSize); i++ {
		s := float64(i)/bladesPertile*gridSize + imagePadding
		newGrass.blades = append(newGrass.blades, newGrassBlade(s*(1-interpX)+t*interpX, s*(1-interpY)+t*interpY, direction))
	}

	return &newGrass
}

const GRASS_FREQUENCY = 0.4
const DAMPING_RATIO = 0.3
const WIND_SPEED = 0.45
const WIND_BIAS = 0.3
const WIND_STRENGTH = 0.03
const SPRITE_WIDTH, SPRITE_HEIGHT = 4, 16

type grassBlade struct {
	sprite      *ebiten.Image
	direction   maths.Direction
	angle       float64
	angularVel  float64
	angleOffset float64
	posX, posY  float64
}

func (g *grassBlade) Update(playerX, playerY, playerVelX, playerVelY float64, perlin *perlin.Perlin, time float64) {
	// Interaction forces
	// TODO: Add interactions from slamboxes and other entities

	relX, relY := 0.0, 0.0
	LR := 0.0
	windS := 1.0
	switch g.direction {
	case maths.DirUp:
		relX = playerVelX
		relY = playerVelY
		LR = g.posX - playerX
	case maths.DirDown:
		windS = -1
		relX = -playerVelX
		relY = playerVelY
		LR = playerX - g.posX
	case maths.DirRight:
		windS = 0
		relX = playerVelY
		relY = playerVelX
		LR = g.posY - playerY
	case maths.DirLeft:
		windS = 0
		relX = -playerVelY
		relY = playerVelX
		LR = playerY - g.posY
	}

	playerDist := maths.Length(playerX-g.posX, playerY-g.posY)
	velocityEffect := 1 - maths.SmoothStep(10, 15, playerDist)

	velocityStrengthX := 0.0
	if relX < 0 {
		velocityStrengthX = maths.SmoothStep(-50, 0, relX) - 1
	} else {
		velocityStrengthX = maths.SmoothStep(0, 50, relX)
	}

	velocityStrengthY := maths.SmoothStep(0, 40, math.Abs(relY))
	velocityStrengthY *= math.Copysign(1, LR)
	// velocityStrengthY *= s
	velocityStrength := velocityStrengthX + velocityStrengthY

	g.angularVel += velocityEffect * velocityStrength * 3

	// Wind forces
	// TODO: Calculate once per cell instead of once per blade
	noiseFactor := perlin.Noise3D(g.posX, g.posY, time*WIND_SPEED) + WIND_BIAS*windS
	g.angularVel += noiseFactor * WIND_STRENGTH

	g.angularVel += (-2*DAMPING_RATIO*GRASS_FREQUENCY*g.angularVel - math.Pow(GRASS_FREQUENCY, 2)*g.angle) * 0.2
	g.angle += g.angularVel * 0.2
}

func (g *grassBlade) init(grassTilemap *ebiten.Image) {
	tileIndex := rand.IntN(8)
	grassSprite := grassTilemap.SubImage(image.Rect(tileIndex*SPRITE_WIDTH, 0, (tileIndex+1)*SPRITE_WIDTH, SPRITE_HEIGHT)).(*ebiten.Image)
	g.sprite = grassSprite
}

func newGrassBlade(posX, posY float64, direction maths.Direction) *grassBlade {
	return &grassBlade{
		posX: posX, posY: posY, direction: direction,
	}
}
