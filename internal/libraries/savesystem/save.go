package save

import (
	"encoding/json"
	"fmt"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/fileio"
	"mask_of_the_tomb/internal/core/resources"
	"os"
	"path/filepath"
)

type SaveData struct {
	SpawnRoomName string
	Settings      resources.SettingsConfig
}

// TODO: Consider: should saving to the disk also be done in a separate thread?
func SaveGame(data SaveData, profile int) {
	savePath := getSavePath(profile)
	exists := errs.Must(fileio.Exists(savePath))
	if !exists {
		os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	}
	file := errs.Must(os.Create(savePath))
	defer file.Close()
	errs.MustSingle(json.NewEncoder(file).Encode(&data))
}

func getSavePath(profile int) string {
	return filepath.Join("save", fmt.Sprintf("save%d.json", profile))
}

func NewSave() SaveData {
	return SaveData{
		Settings: resources.SettingsConfig{
			MasterVolume: 75,
			SoundVolume:  100,
			MusicVolume:  100,
		},
	}
}
