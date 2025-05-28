package node

type selectable interface {
	SetSelected()
	SetDeselected()
}

type ConfirmInfo struct {
	IsConfirmed bool
	SliderVal   float64
}

type OverWriteInfo struct {
	SliderVal float64
}
