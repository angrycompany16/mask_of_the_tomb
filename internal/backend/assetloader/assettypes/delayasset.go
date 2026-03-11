package assettypes

import (
	"io/fs"
	"time"
)

type delayAsset struct { // Debug asset which takes a long time to load
	timeout time.Duration
}

func (a *delayAsset) Load(fs fs.FS) (any, error) {
	time.Sleep(a.timeout)
	return 1, nil
}

func NewDelayAsset(timeout time.Duration) *delayAsset {
	return &delayAsset{
		timeout: timeout,
	}
}
