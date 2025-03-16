package entities

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	entityManagerSingleton = entityManager{
		entities: make(map[string]entity, 0),
	}
)

type entity interface {
	Update()
}

type preUpdater interface {
	PreUpdate()
}

type postUpdater interface {
	PostUpdate()
}

type entityManager struct {
	entities map[string]entity
}

// TODO: This can probably be optimized a lot given that we are looping over all entities three times
func UpdateEntities() {
	for _, entity := range entityManagerSingleton.entities {
		if preUpdater, ok := entity.(preUpdater); ok {
			preUpdater.PreUpdate()
		}
	}

	for _, entity := range entityManagerSingleton.entities {
		entity.Update()
	}

	for _, entity := range entityManagerSingleton.entities {
		if postUpdater, ok := entity.(postUpdater); ok {
			postUpdater.PostUpdate()
		}
	}
}

func RegisterEntity(_entity entity, id string) {
	if id == "" {
		id = uuid.NewString()
		fmt.Println("No ID was found for entity", _entity)
		fmt.Println("Using randomly generated ID", id)
	}

	entityManagerSingleton.entities[id] = _entity
}
