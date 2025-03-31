package fileio

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
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

	for _, file := range results {
		fmt.Println(file)
	}
	return results
}
