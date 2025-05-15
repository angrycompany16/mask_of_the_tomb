package assettypes

import "time"

type delayAsset struct { // Debug asset which takes a long time to load
	timeout time.Duration
}

func (a *delayAsset) Load() {
	time.Sleep(a.timeout)
}

func NewDelayAsset(timeout time.Duration) delayAsset {
	return delayAsset{
		timeout: timeout,
	}
}
