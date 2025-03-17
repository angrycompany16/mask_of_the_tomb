package entities

import "fmt"

type Entity struct {
	id       string
	children map[string]*Entity
	parent   *Entity
}

type entity struct {
	*Entity
	self updater
}

type updater interface {
	Update()
}

type preUpdater interface {
	PreUpdate()
}

type postUpdater interface {
	PostUpdate()
}

type drawer interface {
	Draw()
}

func (e *Entity) SetParent(parent *Entity) {
	e.parent = parent
	parent.children[e.id] = e
}

func (e *Entity) AddChild(child *Entity) {
	e.children[child.id] = child
	child.parent = e
}

func (_Entity *Entity) GetUniqueID() string {
	entityList := _entityManager.entities[_Entity.id]
	for i, _entity := range entityList {
		if _Entity == _entity.Entity {
			return fmt.Sprintf("%s%d", _entity.id, i)
		}
	}
	return _Entity.id
}

func (_Entity *Entity) GetID() string {
	return _Entity.id
}

func (_Entity *Entity) GetChildren() map[string]*Entity {
	return _Entity.children
}

func (_Entity *Entity) GetParent() *Entity {
	return _Entity.parent
}
