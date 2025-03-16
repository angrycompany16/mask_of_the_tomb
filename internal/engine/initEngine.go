package engine

import "mask_of_the_tomb/internal/engine/events"

func Init() {
	events.InitEventManager()
}
