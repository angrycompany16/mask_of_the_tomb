package main

import save "mask_of_the_tomb/internal/game/core/savesystem"

func main() {
	save.GlobalSave.GameData = save.NewGameData()
	save.GlobalSave.SaveGame()
}
