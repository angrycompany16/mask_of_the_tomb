package slamboxactor

import (
	"image/color"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/opgen"
	"mask_of_the_tomb/internal/backend/vector64"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/game/actors/tracker"
	"mask_of_the_tomb/internal/utils"

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
	s.Tracker.Init(cmd)
	s.backendIndex = cmd.SlamboxEnv().AddSlambox(s.rect)
}

func (s *Slambox) Update(cmd *engine.Commands) {
	s.Tracker.Update(cmd)

	x, y := s.Tracker.GetPos()
	s.rect.SetPos(x, y)
	gw, gh := cmd.Renderer().GetGameSize()
	s.Transform2D.SetPos(x-gw/2, y-gh/2)

	if s.slamRequest == maths.DirNone && !s.inChain && !s.inGroup {
		return
	}

	targetX, targetY := cmd.SlamboxEnv().SlamSlambox(s.backendIndex, s.slamRequest)
	s.Tracker.SetTarget(targetX, targetY)
	s.slamRequest = maths.DirNone
}

func (s *Slambox) DrawGizmo(cmd *engine.Commands) {
	s.Tracker.DrawGizmo(cmd)
	s.gizmosImage.Clear()
	vector64.StrokeRect(s.gizmosImage, 0, 0, s.rect.Width()-1, s.rect.Height()-1, 1, color.RGBA{255, 0, 0, 255}, false)

	camX, camY := s.GetCamera().WorldToCamCustomOffset(s.rect.Left(), s.rect.Top(), 0, 0, false)

	cmd.Renderer().Request(opgen.Pos(s.gizmosImage, camX, camY), s.gizmosImage, "Overlay", 0)
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

func defaultSlambox(tracker *tracker.Tracker) *Slambox {
	x, y := tracker.GetPos()
	return &Slambox{
		Tracker:     tracker,
		rect:        maths.NewRect(x, y, 8, 8),
		slamRequest: maths.DirNone,
		gizmosImage: ebiten.NewImage(8, 8),
	}
}

// TODO: Replace rect with width, height
func NewSlambox(tracker *tracker.Tracker, options ...utils.Option[Slambox]) *Slambox {
	slambox := defaultSlambox(tracker)

	for _, option := range options {
		option(slambox)
	}

	return slambox
}

func WithPos(x, y float64) utils.Option[Slambox] {
	return func(s *Slambox) {
		s.rect.SetPos(x, y)
		s.Tracker.SetPos(x, y)
	}
}

func WithSize(width, height float64) utils.Option[Slambox] {
	return func(s *Slambox) {
		s.rect.SetSize(width, height)
		s.gizmosImage = ebiten.NewImage(int(width), int(height))
	}
}
