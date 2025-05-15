package loadparticlesystem

import (
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/libraries/physics/particles"

	"gopkg.in/yaml.v3"
)

// The idea for this function is that it loads whatever is at path and then converts
// it into the correct type and returns that type
func OpenAsset(YAMLdata map[string]any) (string, any) {
	_type := YAMLdata["Type"].(string)
	switch _type {
	case "ParticleSystem":
		rawData := errs.Must(yaml.Marshal(YAMLdata))
		particleSystem := &particles.ParticleSystem{}
		err := yaml.Unmarshal(rawData, &particleSystem)
		if err != nil {
			return _type, particleSystem
		}

		return _type, particleSystem
	case "UIasset":
		// ...
	}
	return _type, nil
}
