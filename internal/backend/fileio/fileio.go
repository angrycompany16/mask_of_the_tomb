package fileio

import (
	"errors"
	"fmt"
	"io/fs"
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

func UnmarshalStruct(path string, out any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(out)
	if err != nil {
		return err
	}
	return nil
}

func UnmarshalMap(path string) (map[string]any, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := make(map[string]any)
	err = yaml.Unmarshal(contents, result)
	if err != nil {
		fmt.Println("Failed to unmarshal", path)
		return nil, err
	}

	return result, nil
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
