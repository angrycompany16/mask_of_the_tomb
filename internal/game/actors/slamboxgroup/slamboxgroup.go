package slamboxgroup

import (
	"mask_of_the_tomb/internal/backend/events"
	"mask_of_the_tomb/internal/backend/maths"
	"mask_of_the_tomb/internal/backend/slambox"
	"mask_of_the_tomb/internal/engine"
	"mask_of_the_tomb/internal/engine/commands"
	"mask_of_the_tomb/internal/game/actors/tracker"
	"mask_of_the_tomb/internal/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type SlamboxGroup struct {
	*tracker.Tracker
	rect         *maths.Rect
	subrects     []*maths.Rect
	backendGroup *slambox.SlamboxGroup
	BackendIndex int
	hitRectIndex int

	slamRequest  maths.Direction
	gizmosImage  *ebiten.Image
	OnMoveFinish *events.EventBus
	hasParticles bool
}

func (s *SlamboxGroup) Init(cmd *commands.Commands) {
	s.Tracker.Init(cmd)
	slamboxenv, ok := commands.Get[slambox.SlamboxEnvironment](cmd)
	if !ok {
		panic("Missing slambox env (Slambox)")
	}

	s.BackendIndex = slamboxenv.AddSlamboxGroup(
		slambox.NewSlamboxGroup(
			append(s.subrects, s.rect),
			len(s.subrects),
		),
	)

	s.backendGroup = slamboxenv.GetSlamboxGroups()[s.BackendIndex]

	x, y := s.Tracker.GetPos()
	s.backendGroup.SetPos(x, y)
	s.Transform2D.SetPos(x, y)
}

func (s *SlamboxGroup) Update(cmd *commands.Commands) {
	s.Tracker.Update(cmd)

	slamboxenv, _ := commands.Get[slambox.SlamboxEnvironment](cmd)
	scene, _ := commands.Get[engine.Scene](cmd)

	x, y := s.Tracker.GetPos()
	s.backendGroup.SetPos(x, y)
	s.Transform2D.SetPos(x, y)

	if data, raised := s.OnMoveFinish.Poll(); raised && s.hasParticles {
		dir := data["dir"].(maths.Direction)
		if s.hitRectIndex == len(s.subrects) {
			scene.SpawnBundle(cmd, MakeSlamboxParticlesBundle(s.rect, dir))
		} else {
			scene.SpawnBundle(cmd, MakeSlamboxParticlesBundle(s.subrects[s.hitRectIndex], dir))
		}
	}

	if s.slamRequest == maths.DirNone {
		return
	}

	slammedRects, i := slamboxenv.SlamSlamboxGroup(s.BackendIndex, s.slamRequest)
	s.hitRectIndex = i
	s.Tracker.SetTarget(slammedRects[s.backendGroup.CenterIndex].X, slammedRects[s.backendGroup.CenterIndex].Y)
	s.slamRequest = maths.DirNone
}

//func (s *Slambox) DrawGizmo(cmd *commands.Commands) {
//	s.Tracker.DrawGizmo(cmd)
//	s.gizmosImage.Clear()
//	vector64.StrokeRect(s.gizmosImage, 0, 0, s.rect.Width-1, s.rect.Height-1, 1, color.RGBA{255, 0, 0, 255}, false)
//
//	camX, camY := s.GetCamera().WorldToCam(s.rect.Left(), s.rect.Top(), false)
//
//	cmd.Renderer.Request(opgen.Pos(s.gizmosImage, camX, camY), s.gizmosImage, renderer.RenderTarget{
//		renderer.SCREEN,
//		"Overlay",
//	}, 0)
//}

func (s *SlamboxGroup) RequestSlam(dir maths.Direction) {
	s.slamRequest = dir
}

func (s *SlamboxGroup) GetSlamRequest() maths.Direction {
	return s.slamRequest
}

func defaultSlambox(tracker *tracker.Tracker) *SlamboxGroup {
	x, y := tracker.GetPos()
	return &SlamboxGroup{
		Tracker:      tracker,
		rect:         maths.NewRect(x, y, 8, 8),
		subrects:     make([]*maths.Rect, 0),
		slamRequest:  maths.DirNone,
		gizmosImage:  ebiten.NewImage(8, 8),
		OnMoveFinish: events.NewBusFrom(tracker.OnMoveFinishEv),
		hasParticles: true,
	}
}

// TODO: The better approach would be to make the particles
// entirely possible to specify, so that playerparticles
// and slambox particles became the same thing lowkey
// Or we could extract the particles to some other place
// That may honestly be the better move, but it might not
// matter that much.
func NewSlamboxGroup(tracker *tracker.Tracker, options ...utils.Option[SlamboxGroup]) *SlamboxGroup {
	slambox := defaultSlambox(tracker)

	for _, option := range options {
		option(slambox)
	}

	return slambox
}

func WithRects(mainRect *maths.Rect, subrects []*maths.Rect) utils.Option[SlamboxGroup] {
	return func(s *SlamboxGroup) {
		s.rect = mainRect
		s.subrects = subrects
	}
}

func WithHasParticles(hasParticles bool) utils.Option[SlamboxGroup] {
	return func(s *SlamboxGroup) {
		s.hasParticles = hasParticles
	}
}
