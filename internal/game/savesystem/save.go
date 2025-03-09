package save

import (
	"encoding/json"
	"fmt"
	"os"
)

// Improvement idea: Change so that we have a map indexed by the struct name, so we don't
// have to be so careful about duplicates

var (
	GlobalSave = Save{
		GameData: NewGameData(),
		savePath: savePath,
	}
)

type gameData struct {
}

func NewGameData() gameData {
	return gameData{}
}

type Save struct {
	GameData gameData
	savePath string
}

func (s *Save) SaveGame() {
	file, err := os.Create(s.savePath)
	if err != nil {
		fmt.Println("Could not open file ", s.savePath)
		fmt.Println(err)
		return
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(&s.GameData)
	if err != nil {
		fmt.Println("Could not write save data to ", s.savePath)
		fmt.Println(err)
		return
	}
}

func (s *Save) LoadGame() {
	gameData := NewGameData()
	file, err := os.Open(s.savePath)
	if err != nil {
		fmt.Println("Could not open file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&gameData)
	if err != nil {
		fmt.Println("Could not decode JSON")
		fmt.Println(err)
		return
	}
	s.GameData = gameData
}
