package entities

var (
	_entityManager = entityManager{
		entities: make(map[string][]*entity, 0),
	}
)

type entityManager struct {
	entities map[string][]*entity
}

// TODO: This can probably be optimized a lot given that we are looping over all entities three times
func PreUpdate() {
	for _, entityList := range _entityManager.entities {
		for _, entity := range entityList {
			if preUpdater, ok := entity.self.(preUpdater); ok {
				preUpdater.PreUpdate()
			}
		}
	}
}

func Update() {
	for _, entity := range _entityManager.entities {
		for _, entity := range entity {
			entity.self.Update()

			if drawable, ok := entity.self.(drawer); ok {
				drawable.Draw()
			}
		}
	}
}

func PostUpdate() {
	for _, entity := range _entityManager.entities {
		for _, entity := range entity {
			if postUpdater, ok := entity.self.(postUpdater); ok {
				postUpdater.PostUpdate()
			}
		}
	}
}

func RegisterEntity(_updater updater, id string) *Entity {
	_Entity := Entity{
		id:       id,
		parent:   nil,
		children: make(map[string]*Entity),
	}
	_entity := entity{
		Entity: &_Entity,
		self:   _updater,
	}

	if _, ok := _entityManager.entities[id]; ok {
		_entityManager.entities[id] = append(_entityManager.entities[id], &_entity)
	} else {
		_entityManager.entities[id] = []*entity{&_entity}
	}

	return &_Entity
}
