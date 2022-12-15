package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func characterReader() CharacterNodes {

	jsonFile, err := os.Open("data/rickandmortycharacter.json")

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	data := CharacterNodes{}
	err = json.NewDecoder(jsonFile).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	return data
}

type fooHandler struct {
	Message string
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(f.Message))
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar called"))
}

func charactersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := characterReader()
		_characterJson, err := json.MarshalIndent(data.CharacterNodes, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(_characterJson)

	}
}

func characterHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "character/")

	charID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	character, itemId := findCharacterByID(charID)
	if character == nil {
		fmt.Println("Error at ID : ", itemId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		_characterJson, err := json.MarshalIndent(character, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(_characterJson)
	case http.MethodPut:
		var updatedCharacter Character
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedCharacter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedCharacter.Id != charID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		character = &updatedCharacter
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func findCharacterByID(id int) (*Character, int) {
	data := characterReader()
	if len(data.CharacterNodes) >= id {
		return &data.CharacterNodes[id-1], id
	}
	return nil, 0
}

func main() {

	http.Handle("/foo", &fooHandler{Message: "hello World"})
	http.HandleFunc("/bar", barHandler)
	http.HandleFunc("/character", charactersHandler)
	http.HandleFunc("/character/", characterHandler)

	http.ListenAndServe(":5000", nil)
	// defer the closing of our jsonFile so that we can parse it later on
}
