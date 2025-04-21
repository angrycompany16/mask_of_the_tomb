package files

import (
	"errors"
	"io/fs"
	"os"
)

// TODO: Add some stuff from fileio.go
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
