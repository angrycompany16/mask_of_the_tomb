package globaldata

import (
	"encoding/json"
	"fmt"
	"mask_of_the_tomb/internal/backend/fileio"
	"os"
	"path/filepath"
)

const SAVEDIR = "save"

// A problem -- There are actually two different levels of persistence existing here.
// On one hand, Config should be loaded upon opening the game, and does not have a saveprofile.
// On the other hand, Inventory and UnlockedDoors should be loaded upon pressing play and *do* have a saveprofile.
// So we actually need to split it into three I guess?
type PersistentData struct {
	Config  Config
	Profile Profile
}

type Profile struct {
	Inventory     Inventory
	UnlockedDoors map[string]bool
}

type Config struct {
	MasterVolume float64
	MusicVolume  float64
	SfxVolume    float64
}

type Inventory struct {
	Keys map[string]*Key
}

type Key struct {
	OpensIid string
	Used     bool
}

func (p *PersistentData) makeSavePath(profile int, isConfig bool) string {
	if isConfig {
		return filepath.Join(SAVEDIR, fmt.Sprintf("config.json"))
	} else {
		return filepath.Join(SAVEDIR, fmt.Sprintf("save%d.json", profile))
	}
}

func (p *PersistentData) ensureSavePathExists(profile int, isConfig bool) error {
	savepath := p.makeSavePath(profile, isConfig)

	saveFileExists, err := fileio.Exists(savepath)
	if err != nil {
		return err
	}

	if !saveFileExists {
		err2 := os.MkdirAll(filepath.Dir(SAVEDIR), os.ModePerm)
		if err2 != nil {
			return err2
		}

		_, err2 = os.Create(savepath)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func (p *PersistentData) Save(profile int, isConfig bool) {
	savepath := p.makeSavePath(profile, isConfig)
	err := p.ensureSavePathExists(profile, isConfig)
	if err != nil {
		fmt.Errorf("SAVE: Problem with save directory:", err)
		return
	}

	file, err := os.Create(savepath)

	if err != nil {
		fmt.Errorf("SAVE: Error when creating save file:", err)
		return
	}

	defer file.Close()
	var content any
	if isConfig {
		content = &p.Config
	} else {
		content = &p.Profile
	}

	err = json.NewEncoder(file).Encode(content)

	if err != nil {
		fmt.Errorf("SAVE: Error when saving:", err)
		return
	}
}

func (p *PersistentData) Load(profile int, isConfig bool) {
	savepath := p.makeSavePath(profile, isConfig)
	err := p.ensureSavePathExists(profile, isConfig)
	if err != nil {
		fmt.Errorf("LOAD: Problem with save directory:", err)
		return
	}

	file, err := os.Open(savepath)
	if err != nil {
		fmt.Errorf("LOAD: Error when opening file:", err)
		return
	}
	defer file.Close()

	var target any
	if isConfig {
		target = &p.Config
	} else {
		target = &p.Profile
	}

	err = json.NewDecoder(file).Decode(target)
	if err != nil {
		fmt.Errorf("LOAD: Error when decoding save data:", err)
		return
	}

	fmt.Println("Finished loading:", target)
}

func (p *PersistentData) UnlockDoor(entityIid string) {
	p.Profile.UnlockedDoors[entityIid] = true
}

func (i *Inventory) HasKey(doorIid string) (bool, string) {
	for keyIid, key := range i.Keys {
		if key.OpensIid == doorIid {
			return true, keyIid
		}
	}
	return false, ""
}

func newConfig() Config {
	return Config{
		MasterVolume: 0.5,
		MusicVolume:  0.5,
		SfxVolume:    0.5,
	}
}

func newProfile() Profile {
	return Profile{
		Inventory: Inventory{
			Keys: make(map[string]*Key),
		},
		UnlockedDoors: make(map[string]bool),
	}
}

func newPersistentData() PersistentData {
	return PersistentData{
		Config:  newConfig(),
		Profile: newProfile(),
	}
}
