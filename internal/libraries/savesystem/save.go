package save

import (
	"encoding/json"
	"fmt"
	"mask_of_the_tomb/internal/core/errs"
	"mask_of_the_tomb/internal/core/fileio"
	"os"
	"path/filepath"
)

type GameData struct {
	SpawnRoomName string
}

func SaveGame(data GameData, profile int) {
	savePath := getSavePath(profile)
	exists := errs.Must(fileio.Exists(savePath))
	if !exists {
		os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	}
	file := errs.Must(os.Create(savePath))
	defer file.Close()
	errs.MustSingle(json.NewEncoder(file).Encode(&data))
}

func LoadGame(profile int) GameData {
	savePath := getSavePath(profile)
	exists := errs.Must(fileio.Exists(savePath))
	if !exists {
		SaveGame(GameData{}, profile)
		return GameData{}
	}

	gameData := GameData{}
	file := errs.Must(os.Open(savePath))
	defer file.Close()

	errs.MustSingle(json.NewDecoder(file).Decode(&gameData))
	return gameData
}

func getSavePath(profile int) string {
	return filepath.Join("save", fmt.Sprintf("save%d.json", profile))
}
