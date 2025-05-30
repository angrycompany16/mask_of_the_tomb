package save

import (
	"encoding/json"
	"errors"
	"mask_of_the_tomb/internal/core/assetloader"
	"os"
)

type SaveAsset struct {
	profile  int
	SaveData SaveData
}

func (a *SaveAsset) Load() error {
	savePath := getSavePath(a.profile)

	_, err := os.Stat(savePath)
	if errors.Is(err, os.ErrNotExist) {
		SaveGame(NewSave(), a.profile)
		return nil
	} else if err != nil {
		return err
	}

	saveData := SaveData{}
	file, err := os.Open(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&saveData)
	if err != nil {
		return err
	}
	a.SaveData = saveData
	return nil
}

func GetSaveAsset(name string) (SaveData, error) {
	saveAsset, err := assetloader.GetAsset(name)
	return saveAsset.(*SaveAsset).SaveData, err
}

func MakeSaveAsset(profile int) *SaveAsset {
	return &SaveAsset{
		profile: profile,
	}
}
