package autotile

import (
	"image"
	"mask_of_the_tomb/internal/core/ebitenrenderutil"
	"mask_of_the_tomb/internal/core/maths"

	"github.com/hajimehoshi/ebiten/v2"
)

type tileKind int

const (
	EMPTY tileKind = iota
	FREE
	WALL
	SPIKE
)

type tileType int

const (
	TOP_LEFT tileType = iota
	TOP
	TOP_RIGHT
	LEFT
	CENTER
	RIGHT
	BOTTOM_LEFT
	BOTTOM
	BOTTOM_RIGHT

	TOP_LEFT_INNER
	TOP_RIGHT_INNER
	BOTTOM_LEFT_INNER
	BOTTOM_RIGHT_INNER

	SPIKE_TOP
	SPIKE_BOTTOM
	SPIKE_LEFT
	SPIKE_RIGHT
)

type tilePresence int

type tileNeighbourData [8]tileKind
type tilemapSrcData map[tileType]image.Rectangle
type tileRule struct {
	_type tileType
	def   [8]tileKind
}
type RectList struct {
	List []*maths.Rect
	Kind tileKind
}

// Enable conversion from rects to tile data
func CreateSprite(
	srcTilemap *ebiten.Image,
	dst *ebiten.Image,
	tilemapSrcData tilemapSrcData,
	ruleset []tileRule,

	tileSize float64,
	rect *maths.Rect,
	kind tileKind,
	// Not good
	// but we will just roll with it
	neighbourLists ...RectList,
) {
	drawTile := func(i, j int) {
		localX, localY := float64(i)*tileSize, float64(j)*tileSize
		worldX := localX + rect.Left() + tileSize/2
		worldY := localY + rect.Top() + tileSize/2

		// Painful and needs a refactor probably
		var ul, um, ur, ml, mr, bl, bm, br bool
		ul = rect.Contains(worldX-tileSize, worldY-tileSize)
		um = rect.Contains(worldX, worldY-tileSize)
		ur = rect.Contains(worldX+tileSize, worldY-tileSize)
		ml = rect.Contains(worldX-tileSize, worldY)
		mr = rect.Contains(worldX+tileSize, worldY)
		bl = rect.Contains(worldX-tileSize, worldY+tileSize)
		bm = rect.Contains(worldX, worldY+tileSize)
		br = rect.Contains(worldX+tileSize, worldY+tileSize)

		tileData := tileNeighbourData{
			fromBool(kind, EMPTY, ul),
			fromBool(kind, EMPTY, um),
			fromBool(kind, EMPTY, ur),
			fromBool(kind, EMPTY, ml),
			fromBool(kind, EMPTY, mr),
			fromBool(kind, EMPTY, bl),
			fromBool(kind, EMPTY, bm),
			fromBool(kind, EMPTY, br),
		}

		for _, neighbourList := range neighbourLists {
			var ul, um, ur, ml, mr, bl, bm, br bool
			for _, wall := range neighbourList.List {
				ul = ul || wall.Contains(worldX-tileSize, worldY-tileSize)
				um = um || wall.Contains(worldX, worldY-tileSize)
				ur = ur || wall.Contains(worldX+tileSize, worldY-tileSize)
				ml = ml || wall.Contains(worldX-tileSize, worldY)
				mr = mr || wall.Contains(worldX+tileSize, worldY)
				bl = bl || wall.Contains(worldX-tileSize, worldY+tileSize)
				bm = bm || wall.Contains(worldX, worldY+tileSize)
				br = br || wall.Contains(worldX+tileSize, worldY+tileSize)
			}

			tileData = tileNeighbourData{
				fromBool(neighbourList.Kind, tileData[0], ul),
				fromBool(neighbourList.Kind, tileData[1], um),
				fromBool(neighbourList.Kind, tileData[2], ur),
				fromBool(neighbourList.Kind, tileData[3], ml),
				fromBool(neighbourList.Kind, tileData[4], mr),
				fromBool(neighbourList.Kind, tileData[5], bl),
				fromBool(neighbourList.Kind, tileData[6], bm),
				fromBool(neighbourList.Kind, tileData[7], br),
			}
		}

		tilemapRect := getTile(tileData, ruleset, tilemapSrcData)
		tile := srcTilemap.SubImage(tilemapRect).(*ebiten.Image)
		ebitenrenderutil.DrawAt(tile, dst, localX, localY)
	}

	for i := 0; i < int(rect.Width()/tileSize); i++ {
		drawTile(i, 0)
		drawTile(i, int(rect.Height()/tileSize-1))
	}

	for i := 0; i < int(rect.Height()/tileSize); i++ {
		drawTile(0, i)
		drawTile(int(rect.Width()/tileSize-1), i)
	}

	for i := 1; i < int(rect.Width()/tileSize-1); i++ {
		for j := 1; j < int(rect.Height()/tileSize-1); j++ {
			localX, localY := float64(i)*tileSize, float64(j)*tileSize
			ebitenrenderutil.DrawAt(srcTilemap.SubImage(tilemapSrcData[CENTER]).(*ebiten.Image), dst, localX, localY)
		}
	}
}

func fromBool[T any](_true T, _false T, in bool) T {
	if in {
		return _true
	} else {
		return _false
	}
}

func getTile(data tileNeighbourData, rules []tileRule, tileRectData tilemapSrcData) image.Rectangle {
	for _, rule := range rules {
		if evalRule(data, rule) {
			return tileRectData[rule._type]
		}
	}
	return image.Rect(0, 0, 0, 0)
}

func evalRule(data tileNeighbourData, rule tileRule) bool {
	for i := 0; i < 8; i++ {
		if rule.def[i] == FREE {
			continue
		} else if rule.def[i] != data[i] {
			return false
		}
	}
	return true
}

func GetDefaultTileRectData(tileRegionX, tileRegionY, tileSize float64) tilemapSrcData {
	tileRectData := make(tilemapSrcData)

	tileRectData[BOTTOM_RIGHT_INNER] = image.Rect(
		int(tileRegionX+tileSize*3),
		int(tileRegionY),
		int(tileRegionX+tileSize*4),
		int(tileRegionY+tileSize),
	)

	tileRectData[BOTTOM_LEFT_INNER] = image.Rect(
		int(tileRegionX+tileSize*4),
		int(tileRegionY),
		int(tileRegionX+tileSize*5),
		int(tileRegionY+tileSize),
	)

	tileRectData[TOP_RIGHT_INNER] = image.Rect(
		int(tileRegionX+tileSize*3),
		int(tileRegionY+tileSize),
		int(tileRegionX+tileSize*4),
		int(tileRegionY+tileSize*2),
	)

	tileRectData[TOP_LEFT_INNER] = image.Rect(
		int(tileRegionX+tileSize*4),
		int(tileRegionY+tileSize),
		int(tileRegionX+tileSize*5),
		int(tileRegionY+tileSize*2),
	)

	tileRectData[TOP_LEFT] = image.Rect(
		int(tileRegionX),
		int(tileRegionY),
		int(tileRegionX+tileSize),
		int(tileRegionY+tileSize),
	)

	tileRectData[TOP] = image.Rect(
		int(tileRegionX+tileSize),
		int(tileRegionY),
		int(tileRegionX+2*tileSize),
		int(tileRegionY+tileSize),
	)

	tileRectData[TOP_RIGHT] = image.Rect(
		int(tileRegionX+2*tileSize),
		int(tileRegionY),
		int(tileRegionX+3*tileSize),
		int(tileRegionY+tileSize),
	)

	tileRectData[LEFT] = image.Rect(
		int(tileRegionX),
		int(tileRegionY+tileSize),
		int(tileRegionX+tileSize),
		int(tileRegionY+2*tileSize),
	)

	tileRectData[CENTER] = image.Rect(
		int(tileRegionX+tileSize),
		int(tileRegionY+tileSize),
		int(tileRegionX+2*tileSize),
		int(tileRegionY+2*tileSize),
	)

	tileRectData[RIGHT] = image.Rect(
		int(tileRegionX+2*tileSize),
		int(tileRegionY+tileSize),
		int(tileRegionX+3*tileSize),
		int(tileRegionY+2*tileSize),
	)

	tileRectData[BOTTOM_LEFT] = image.Rect(
		int(tileRegionX),
		int(tileRegionY+2*tileSize),
		int(tileRegionX+tileSize),
		int(tileRegionY+3*tileSize),
	)

	tileRectData[BOTTOM] = image.Rect(
		int(tileRegionX+tileSize),
		int(tileRegionY+2*tileSize),
		int(tileRegionX+2*tileSize),
		int(tileRegionY+3*tileSize),
	)

	tileRectData[BOTTOM_RIGHT] = image.Rect(
		int(tileRegionX+2*tileSize),
		int(tileRegionY+2*tileSize),
		int(tileRegionX+3*tileSize),
		int(tileRegionY+3*tileSize),
	)

	tileRectData[SPIKE_TOP] = image.Rect(
		int(tileRegionX+tileSize*5),
		int(tileRegionY),
		int(tileRegionX+tileSize*6),
		int(tileRegionY+tileSize),
	)

	tileRectData[SPIKE_LEFT] = image.Rect(
		int(tileRegionX+tileSize*6),
		int(tileRegionY),
		int(tileRegionX+tileSize*7),
		int(tileRegionY+tileSize),
	)

	tileRectData[SPIKE_RIGHT] = image.Rect(
		int(tileRegionX+tileSize*5),
		int(tileRegionY+tileSize),
		int(tileRegionX+tileSize*6),
		int(tileRegionY+tileSize*2),
	)

	tileRectData[SPIKE_BOTTOM] = image.Rect(
		int(tileRegionX+tileSize*6),
		int(tileRegionY+tileSize),
		int(tileRegionX+tileSize*7),
		int(tileRegionY+tileSize*2),
	)

	return tileRectData
}

func GetDefaultTileRuleset() []tileRule {
	ruleset := []tileRule{
		tileRule{_type: TOP_LEFT_INNER, def: [8]tileKind{
			EMPTY, WALL, FREE,
			WALL /**/, FREE,
			FREE, FREE, WALL},
		},
		tileRule{_type: TOP_RIGHT_INNER, def: [8]tileKind{
			FREE, WALL, EMPTY,
			FREE /**/, WALL,
			WALL, FREE, FREE},
		},
		tileRule{_type: BOTTOM_LEFT_INNER, def: [8]tileKind{
			FREE, FREE, WALL,
			WALL /**/, FREE,
			EMPTY, WALL, FREE},
		},
		tileRule{_type: BOTTOM_RIGHT_INNER, def: [8]tileKind{
			WALL, FREE, FREE,
			FREE /**/, WALL,
			FREE, WALL, EMPTY},
		},
		tileRule{_type: TOP_LEFT, def: [8]tileKind{
			EMPTY, EMPTY, FREE,
			EMPTY /**/, FREE,
			FREE, FREE, FREE},
		},
		tileRule{_type: TOP_RIGHT, def: [8]tileKind{
			FREE, EMPTY, EMPTY,
			FREE /**/, EMPTY,
			FREE, FREE, FREE},
		},
		tileRule{_type: BOTTOM_LEFT, def: [8]tileKind{
			FREE, FREE, FREE,
			EMPTY /**/, FREE,
			EMPTY, EMPTY, FREE},
		},
		tileRule{_type: BOTTOM_RIGHT, def: [8]tileKind{
			FREE, FREE, FREE,
			FREE /**/, EMPTY,
			FREE, EMPTY, EMPTY},
		},
		tileRule{_type: LEFT, def: [8]tileKind{
			FREE, FREE, FREE,
			EMPTY /**/, WALL,
			FREE, FREE, FREE},
		},
		tileRule{_type: RIGHT, def: [8]tileKind{
			FREE, FREE, FREE,
			WALL /**/, EMPTY,
			FREE, FREE, FREE},
		},
		tileRule{_type: TOP, def: [8]tileKind{
			FREE, EMPTY, FREE,
			FREE /**/, FREE,
			FREE, WALL, FREE},
		},
		tileRule{_type: BOTTOM, def: [8]tileKind{
			FREE, WALL, FREE,
			FREE /**/, FREE,
			FREE, EMPTY, FREE},
		},
		tileRule{_type: CENTER, def: [8]tileKind{
			WALL, WALL, WALL,
			WALL /**/, WALL,
			WALL, WALL, WALL},
		},
	}

	return ruleset
}

func GetDefaultSpikeRules() []tileRule {
	ruleset := []tileRule{
		tileRule{_type: SPIKE_TOP, def: [8]tileKind{
			FREE, EMPTY, FREE,
			FREE /* */, FREE,
			FREE, WALL, FREE},
		},
		tileRule{_type: SPIKE_RIGHT, def: [8]tileKind{
			FREE, FREE, FREE,
			EMPTY /* */, WALL,
			FREE, FREE, FREE},
		},
		tileRule{_type: SPIKE_BOTTOM, def: [8]tileKind{
			FREE, WALL, FREE,
			FREE /* */, FREE,
			FREE, EMPTY, FREE},
		},
		tileRule{_type: SPIKE_LEFT, def: [8]tileKind{
			FREE, FREE, FREE,
			WALL /* */, EMPTY,
			FREE, FREE, FREE},
		},
	}
	return ruleset
}
