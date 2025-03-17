package testentity

import (
	"fmt"
	"mask_of_the_tomb/internal/engine/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TestEntity struct {
	*entities.Entity
	foo int
}

func (t *TestEntity) Update() {

}

func NewTestEntity(data int) *TestEntity {
	return &TestEntity{
		foo: data,
	}
}

type SpecialEntity struct {
	*entities.Entity
	foo int
}

// Hell yeah mr. special entity
func (s *SpecialEntity) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		fmt.Println("I am", s.Entity.GetUniqueID())
		fmt.Println("My children are", s.Entity.GetChildren())
		fmt.Println("My parent is", s.Entity.GetParent().GetUniqueID())
		fmt.Println(s.foo)
	}
}

func NewSpecialEntity(data int) *SpecialEntity {
	return &SpecialEntity{
		foo: data,
	}
}
