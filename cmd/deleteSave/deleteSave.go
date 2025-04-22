package main

import (
	"flag"
	"fmt"
	save "mask_of_the_tomb/internal/game/core/savesystem"
)

var saveProfile int

func main() {
	flag.IntVar(&saveProfile, "saveprofile", 1, "Save profile to delete (99 for dev profile)")

	fmt.Println("-- Deleting save", saveProfile)

	save.SaveGame(save.GameData{}, saveProfile)
}
