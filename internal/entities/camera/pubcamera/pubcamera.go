package pubcamera

const (
	CameraEntityName = "Camera"
)

type CameraAdvertiser struct {
	PosX, PosY float64
}

func (c *CameraAdvertiser) Read() any {
	return *c
}
