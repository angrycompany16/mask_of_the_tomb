package save

import (
	"encoding/json"
	"fmt"
	"os"
)

// Want: to be able to set an attribute of any entity to be 'saved'.
// That attribute will persist between room changes, and be permanently saved when the game
// is saved.
// Thus the attribute needs to

// Idea: This should allow for arbitrary saving and loading using reflection:
// - To save a struct field, get the name w/ reflection, as well as the value
//   and insert these into a map containing the local save data of the game
// - To load the field, get the value by querying into the map with the
//   struct field as name
// - When saving the game, write the map onto a JSON file
// - When loading the game, decode the JSON file into a map and then load all
//   the struct fields

// Should look something like this
/*
struct Player {
	score int (saveable)
}

func SaveScore() {
	sharedSaveManager.Save(player.score)
}

func LoadScore() {
}

func SaveGame() {
	write everything from the map onto the JSON
}

func LoadGame() {
	// Implemented per-object as we need to choose which fields to load
}
*/

// Improvement idea: Change so that we have a map indexed by the struct name, so we don't
// have to be so careful about duplicates

var (
	GlobalSave = Save{
		GameData: NewGameData(),
		savePath: savePath,
	}
)

type gameData struct {
	CollectedEntityUids map[string]bool `json:"CollectedEntityUids"`
}

func NewGameData() gameData {
	return gameData{
		CollectedEntityUids: make(map[string]bool),
	}
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

// 4 hours of programminG:

/*

type SaveManager struct {
	savePath  string
	localSave map[string]any
}

func (s *SaveManager) SaveLocally(data any) {
	a, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(a, &s.localSave)
}

func (s *SaveManager) LoadLocally(entity any) {
	data, err := json.Marshal(s.localSave)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, entity)
}

type TestEntity struct {
	Data          int       `json:"data"`
	RecursiveData SubEntity `json:"subEntity"`
}

type SubEntity struct {
	SubData   int `json:"subData"`
	LocalData int
}

func main() {
	saveManager := SaveManager{
		localSave: make(map[string]any),
	}
	test := TestEntity{
		Data: 40,
		RecursiveData: SubEntity{
			SubData:   400,
			LocalData: 600,
		},
	}
	saveManager.SaveLocally(&test)

	// pretty.Print(saveManager.localSave)

	for key, value := range saveManager.localSave {
		if reflect.ValueOf(value).Kind() == reflect.Map {
			for _key, _value := range value.(map[string]any) {
				fmt.Printf("    [%s], [%v]\n", _key, _value)
			}
		} else {
			fmt.Printf("[%s], [%v]\n", key, value)
		}
	}

	test.Data = 42
	test.RecursiveData = SubEntity{
		SubData:   401,
		LocalData: 601,
	}

	saveManager.LoadLocally(&test)
	// method(&test, TestEntity{40})
	fmt.Printf("%+v\n", test)
}

func SaveLocally(data any, target map[string]any) {

	entityType := reflect.TypeOf(data).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		tag := field.Tag.Get("saved")
		saved, err := strconv.ParseBool(tag)
		if err != nil || !saved {
			continue
		}
		val := reflect.ValueOf(data).Elem().Field(i)
		if entityType.Field(i).Type.Kind() == reflect.Struct {
			target[field.Name] = make(map[string]any)
			SaveLocally(val.Addr().Interface(), target[field.Name].(map[string]any))
		} else {
			target[field.Name] = val
		}
	}
}


// https://www.youtube.com/watch?v=HPdURyombM0
func LoadLocally(entity any, target map[string]any) {
	entityType := reflect.TypeOf(entity).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		value := entityType.Field(i)
		tag := value.Tag.Get("saved")
		saved, err := strconv.ParseBool(tag)
		if err != nil || !saved {
			continue
		}

		savedValue := target[value.Name]

		oldField := reflect.ValueOf(entity).Elem().Field(i)
		newField := reflect.ValueOf(savedValue)

		if value.Type.Kind() == reflect.Struct {
			fmt.Println(oldField)
			LoadLocally(oldField.Addr().Interface(), savedValue.(map[string]any))
		} else {
			fmt.Println("odl field: ", oldField)
			fmt.Println("new ifeld:", newField)
			oldField.Set(newField)
		}
	}
}

func method(existingEntity interface{}, newEntity interface{}) {
	entityType := reflect.TypeOf(existingEntity).Elem()
	for i := 0; i < entityType.NumField(); i++ {
		value := entityType.Field(i)
		tag := value.Tag.Get("saved")
		saved, err := strconv.ParseBool(tag)
		if err != nil || !saved {
			continue
		}

		fmt.Println("test")
		oldField := reflect.ValueOf(existingEntity).Elem().Field(i)
		newField := reflect.ValueOf(newEntity).FieldByName(value.Name)

		if value.Type.Kind() == reflect.Struct {
			method(oldField.Addr().Interface(), newField.Interface())
		} else {
			fmt.Println(oldField)
			fmt.Println(newField)
			oldField.Set(newField)
		}
	}
}
*/
// I actually don't know if i can do this anymore
// AND I'm *still* ass at reflection
// Time to end it all
