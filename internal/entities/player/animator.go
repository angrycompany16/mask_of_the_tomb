package player

import (
	"mask_of_the_tomb/internal/entities/animation"
	"mask_of_the_tomb/internal/libraries/errs"
)

const (
	idleAnim = iota
	dashInitAnim
	dashLoopAnim
	slamAnim
)

var (
	playerAnimationMap = map[int]*animation.Animation{
		idleAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(idleSpritesheetPath)),
			0.14,
			animation.Strip,
			animation.Loop,
			-1,
		),
		dashInitAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(dashInitSpritesheetPath)),
			0.08,
			animation.Strip,
			animation.Once,
			dashLoopAnim,
		),
		dashLoopAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(dashLoopSpritesheetPath)),
			0.08,
			animation.Strip,
			animation.Loop,
			-1,
		),
		slamAnim: animation.NewAnimation(
			animation.NewSpritesheetAuto(errs.MustNewImageFromFile(slamSpritesheetPath)),
			0.08,
			animation.Strip,
			animation.Once,
			idleAnim,
		),
	}
)
