package main

import "mask_of_the_tomb/save"

func main() {
	save.GlobalSave.GameData = save.NewGameData()
	save.GlobalSave.SaveGame()
}
