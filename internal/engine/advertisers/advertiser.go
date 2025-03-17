package advertisers

import "fmt"

// A struct that should allow any entity to expose data for other entities to *only* read from
// Read will return whatever data the advertiser currently contains
var _advertiserManager = AdvertiserManager{
	advertisers: make(map[string]advertiser),
}

// TODO: Try to convert to generic thing?
type advertiser interface {
	Read() any
}

type AdvertiserManager struct {
	advertisers map[string]advertiser
}

func GetAdvertiser(id string) advertiser {
	advertiser, ok := _advertiserManager.advertisers[id]
	if !ok {
		panic(fmt.Sprintf("Could not find advertiser with id %s", id))
	}

	return advertiser
}

func RegisterAdvertiser(_advertiser advertiser, id string) {
	_advertiserManager.advertisers[id] = _advertiser
}
