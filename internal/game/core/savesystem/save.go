package save

import (
	"encoding/json"
	"fmt"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/files"
	"mask_of_the_tomb/internal/game/world/levelmemory"
	"os"
	"path/filepath"
)

type GameData struct {
	WorldStateMemory map[string]levelmemory.LevelMemory
	SpawnRoomName    string
}

func SaveGame(data GameData, profile int) {
	// fmt.Println("Saving game.....")
	// defer fmt.Println("Done!")
	savePath := getSavePath(profile)
	exists := errs.Must(files.Exists(savePath))
	if !exists {
		os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	}
	file := errs.Must(os.Create(savePath))
	defer file.Close()
	errs.MustSingle(json.NewEncoder(file).Encode(&data))
}

func LoadGame(profile int) GameData {
	// fmt.Println("Loading game.....")
	// defer fmt.Println("Done!")
	savePath := getSavePath(profile)
	exists := errs.Must(files.Exists(savePath))
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
