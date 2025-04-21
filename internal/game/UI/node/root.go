package node

import "fmt"

// NOTE: Embedded structs can also be used for custom
// unmarshaling

type Root struct {
	NodeData `yaml:",inline"`
}

func (r *Root) Update(confirmations map[string]bool) {
	r.UpdateChildren(confirmations)
}

func (r *Root) Draw(offsetX, offsetY float64, parentWidth, parentHeight float64) {
	w, h := inheritSize(r.Width, r.Height, parentWidth, parentHeight)
	r.DrawChildren(offsetX+r.PosX, offsetY+r.PosY, w, h)
}

func (r *Root) Reset() {
	fmt.Println("Reset")
	r.ResetChildren()
}
