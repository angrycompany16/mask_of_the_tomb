package assettypes

import "time"

type delayAsset struct { // Debug asset which takes a long time to load
	timeout time.Duration
}

func (a *delayAsset) Load() error {
	time.Sleep(a.timeout)
	return nil
}

func NewDelayAsset(timeout time.Duration) *delayAsset {
	return &delayAsset{
		timeout: timeout,
	}
}
