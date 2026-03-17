package slamboxactor

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game/actors/tracker"

	"github.com/hajimehoshi/ebiten/v2"
)

type Slambox struct {
	*tracker.Tracker
	rect         *maths.Rect
	backendIndex int
	slamRequest  maths.Direction
	inChain      bool
	inGroup      bool
	isCenter     bool
	gizmosImage  *ebiten.Image
}

func (s *Slambox) Init(cmd *engine.Commands) {
	s.backendIndex = cmd.SlamboxEnv().AddSlambox(s.rect)
}

func (s *Slambox) Update(cmd *engine.Commands) {
	s.Tracker.Update(cmd)

	x, y := s.Tracker.GetPos()
	s.rect.SetPos(x, y)

	if s.slamRequest == maths.DirNone && !s.inChain && !s.inGroup {
		return
	}

	targetX, targetY := cmd.SlamboxEnv().SlamSlambox(s.backendIndex, s.slamRequest)
	s.Tracker.SetTarget(targetX, targetY)
	s.slamRequest = maths.DirNone
}

func (s *Slambox) DrawGizmo(cmd *engine.Commands) {
	s.gizmosImage.Clear()
	vector64.StrokeRect(s.gizmosImage, 0, 0, s.rect.Width()-1, s.rect.Height()-1, 1, color.RGBA{255, 0, 0, 255}, false)

	cmd.Renderer().Request(opgen.Pos(s.gizmosImage, s.rect.Left(), s.rect.Top()), s.gizmosImage, "Overlay", 0)
}

func (s *Slambox) RequestSlam(dir maths.Direction) {
	s.slamRequest = dir
}

func (s *Slambox) GetSlamRequest() maths.Direction {
	return s.slamRequest
}

func (s *Slambox) GetRect() *maths.Rect {
	return s.rect
}

func (s *Slambox) IsCenter() bool {
	return s.isCenter
}

// TODO: Replace rect with width, height
func NewSlambox(tracker *tracker.Tracker, rect *maths.Rect) *Slambox {
	return &Slambox{
		Tracker:      tracker,
		rect:         rect,
		backendIndex: 0,
		slamRequest:  maths.DirNone,
		gizmosImage:  ebiten.NewImage(int(rect.Width()), int(rect.Height())),
	}
}
