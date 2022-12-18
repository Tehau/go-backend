package character

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

var characterMap = struct {
	sync.RWMutex
	m map[int]Character
}{m: make(map[int]Character)}

type CharacterNodes struct {
	CharacterNodes []Character `json:"results"`
}

type Character struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Gender  string `json:"gender"`
	Origin  struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"origin"`
	Location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"location"`
	Image   string    `json:"image"`
	Episode []string  `json:"episode"`
	Url     string    `json:"url"`
	Created time.Time `json:"created"`
}

func init() {
	log.Println("Loading data...")
	charMap, err := loadCharacterMap()
	characterMap.m = charMap
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d characters loaded...\n", len(characterMap.m))
}

func loadCharacterMap() (map[int]Character, error) {
	filePath := "data/rickandmortycharacter.json"

	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] doest not exist", filePath)
	}

	jsonFile, err := os.Open(filePath)

	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalln(err)
	}
	data := CharacterNodes{}
	// JSON File in Data Global Variable
	err = json.NewDecoder(jsonFile).Decode(&data)
	if err != nil {
		log.Fatalln(err)
	}

	charMap := make(map[int]Character, len(data.CharacterNodes))
	for i := 0; i < len(data.CharacterNodes); i++ {
		charMap[data.CharacterNodes[i].Id] = data.CharacterNodes[i]
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	return charMap, nil
}

func getCharacter(charID int) *Character {
	characterMap.RLock()
	defer characterMap.RUnlock()
	if char, ok := characterMap.m[charID]; ok {
		return &char
	}
	return nil
}

func removeCharacter(charID int) {
	characterMap.RLock()
	defer characterMap.RUnlock()
	delete(characterMap.m, charID)
}

func getCharacterList() []Character {
	characterMap.RLock()
	characterList := make([]Character, 0, len(characterMap.m))
	for _, value := range characterMap.m {
		characterList = append(characterList, value)
	}
	sort.SliceStable(characterList, func(i, j int) bool {
		return characterList[i].Id < characterList[j].Id
	})
	defer characterMap.RUnlock()
	return characterList
}

func getCharacterIds() []int {
	characterMap.RLock()
	var characterIds []int
	for key := range characterMap.m {
		characterIds = append(characterIds, key)
	}
	characterMap.RUnlock()
	sort.Ints(characterIds)
	return characterIds
}

func getNextCharacterID() int {
	characterIDs := getCharacterIds()
	return characterIDs[len(characterIDs)-1] + 1
}

func addOrUpdateCharacter(character Character) (int, error) {
	addOrUpdateID := -1
	if character.Id > 0 {
		oldCharacter := getCharacter(character.Id)
		if oldCharacter == nil {
			return 0, fmt.Errorf("character id [%d] doesn't exist", character.Id)
		}
		addOrUpdateID = character.Id
	} else {
		addOrUpdateID = getNextCharacterID()
		character.Id = addOrUpdateID
	}
	characterMap.Lock()
	characterMap.m[addOrUpdateID] = character
	characterMap.Unlock()
	return addOrUpdateID, nil

}
