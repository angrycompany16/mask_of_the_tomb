package scenes

import (
	"fmt"
	"mask_of_the_tomb/internal/core/resources"
	"mask_of_the_tomb/internal/core/scene"
	"time"
)

type RootSceneBehaviour struct {
	initTime time.Time
}

// Sets up some basic global data
func (r *RootSceneBehaviour) Init() {
	r.initTime = time.Now()
}

func (r *RootSceneBehaviour) Update() {
	resources.Time = time.Since(r.initTime).Seconds()
}

func (r *RootSceneBehaviour) Draw() {
	fmt.Println("Drawing root scene (nothing)")
}

func (r *RootSceneBehaviour) Transition() (scene.SceneTransition, bool) {
	return scene.SceneTransition{}, false
}
