package scenes

import (
	"fmt"
)

type RootSceneBehaviour struct {
}

// Loads basic UI, sets up all assets
func (r *RootSceneBehaviour) Init() {
	fmt.Println("Root scene intialized")
}

func (r *RootSceneBehaviour) Update() {
	fmt.Println("Root scene updating")
}

func (r *RootSceneBehaviour) Draw() {
	fmt.Println("Drawing root scene (nothing)")
}
