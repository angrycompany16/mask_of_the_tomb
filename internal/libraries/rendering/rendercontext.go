package rendering

// A problem: This object will be used in lots of different places, and so it will
// have to be imported by many different things, however this is the opposite of what
// we want...
type RenderContext struct {
	camX, camY float64
}
