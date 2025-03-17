package main

import save "mask_of_the_tomb/internal/libraries/savesystem"

func main() {
	save.GlobalSave.GameData = save.NewGameData()
	save.GlobalSave.SaveGame()
}
