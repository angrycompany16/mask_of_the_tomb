package node

type selectable interface {
	SetSelected(suppressSound bool)
	SetDeselected()
}

type ConfirmInfo struct {
	IsConfirmed bool
	SliderVal   float64
}

type OverWriteInfo struct {
	SliderVal float64
}
