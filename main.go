package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
		_characterJson, err := json.Marshal(data.CharacterNodes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(_characterJson)

	}
}

func main() {
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

	//for i := 0; i < len(data.CharacterNodes); i++ {
	//	fmt.Printf("Name: %s\n", data.CharacterNodes[i].Name)
	//	fmt.Printf("Status: %s\n", data.CharacterNodes[i].Status)
	//}

	http.Handle("/foo", &fooHandler{Message: "hello World"})
	http.HandleFunc("/bar", barHandler)
	http.HandleFunc("/character", charactersHandler)
	http.ListenAndServe(":5000", nil)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
}
