# Scene system

The skeleton of the entire scene system is the game.go file with the Game struct.

## Scenes

A scene is defined by the following:

```go
type Scene interface {
	Init()
	Update(sceneStack *scene.SceneStack) (*scene.SceneTransition, bool)
	Draw()
	GetName() string
}
```

In other words, scenes have init, update, draw, and a name. Scenes are added to a stack, and there is functionality for popping the first scene, removing scenes by name, and more.

Switching scenes is done by returning a `SceneTransition` from the update method as well as `true`

Note: It is (apparently) convention to use camelCase for whatever GetName() returns.