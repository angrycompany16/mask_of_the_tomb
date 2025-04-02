package fileio

import (
	"fmt"
	"io/fs"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/game/core/rendering"
	"mask_of_the_tomb/internal/game/physics/particles"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func FindFiles(name string, results []string) []string {
	err := filepath.WalkDir("assets", func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if !dirEntry.IsDir() && strings.Contains(path, name) {
			results = append(results, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return results
}

func OpenAsset(path string) (string, any) {
	contents := errs.Must(os.ReadFile(path))

	result := make(map[string]any)
	err := yaml.Unmarshal(contents, result)
	if err != nil {
		fmt.Println("Failed to unmarshal", path)
		return "", particles.ParticleSystem{}
	}
	// fmt.Println("type:", result["Type"])
	_type := result["Type"].(string)
	switch _type {
	case "ParticleSystem":
		// TODO: Need to add the rendering layer into the file config
		particleSystem := errs.Must(particles.FromFile(path, rendering.RenderLayers.Playerspace))
		return _type, particleSystem
	}
	return _type, nil
}
