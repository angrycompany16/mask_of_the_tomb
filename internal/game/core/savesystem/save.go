package save

import (
	"encoding/json"
	"mask_of_the_tomb/internal/errs"
	"mask_of_the_tomb/internal/files"
	"mask_of_the_tomb/internal/game/world/levelmemory"
	"os"
	"path/filepath"
)

var (
	savePath = filepath.Join("save", "savedata.json")
)

type GameData struct {
	WorldStateMemory map[string]levelmemory.LevelMemory
	SpawnRoomName    string
}

func SaveGame(data GameData) {
	// fmt.Println("Saving game.....")
	// defer fmt.Println("Done!")
	exists := errs.Must(files.Exists(savePath))
	if !exists {
		os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	}
	file := errs.Must(os.Create(savePath))
	defer file.Close()
	errs.MustSingle(json.NewEncoder(file).Encode(&data))
}

func LoadGame() GameData {
	// fmt.Println("Loading game.....")
	// defer fmt.Println("Done!")
	exists := errs.Must(files.Exists(savePath))
	if !exists {
		SaveGame(GameData{})
		return GameData{}
	}

	gameData := GameData{}
	file := errs.Must(os.Open(savePath))
	defer file.Close()

	errs.MustSingle(json.NewDecoder(file).Decode(&gameData))
	return gameData
}
