package player

import (
	"mask_of_the_tomb/assets"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/libraries/animation"
	"path/filepath"
	"time"
)

// Connects the player and animation libraries
// But note! This also requires the fsm module in order to decide when each animation
// should be played

const (
	idleAnim = iota
	dashInitAnim
	dashLoopAnim
	slamAnim
)

var (
	idleSpritesheetPath     = filepath.Join(assets.PlayerFolder, "player-idle-Sheet.png")
	dashInitSpritesheetPath = filepath.Join(assets.PlayerFolder, "player-init-jump-Sheet.png")
	dashLoopSpritesheetPath = filepath.Join(assets.PlayerFolder, "player-loop-jump-Sheet.png")
	slamSpritesheetPath     = filepath.Join(assets.PlayerFolder, "player-slam-Sheet.png")

	playerAnimationMap = map[int]*animation.Animation{
		idleAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(idleSpritesheetPath)),
			time.Millisecond*140,
			animation.Strip,
			animation.Loop,
			-1,
		),
		dashInitAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(dashInitSpritesheetPath)),
			time.Millisecond*80,
			animation.Strip,
			animation.Once,
			dashLoopAnim,
		),
		dashLoopAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(dashLoopSpritesheetPath)),
			time.Millisecond*80,
			animation.Strip,
			animation.Loop,
			-1,
		),
		slamAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(slamSpritesheetPath)),
			time.Millisecond*80,
			animation.Strip,
			animation.Once,
			idleAnim,
		),
	}
)
