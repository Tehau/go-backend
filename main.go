package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type fooHandler struct {
	Message string
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(f.Message))
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar called"))
}

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

func main() {
	jsonFile, err := os.Open("data/rickandmortycharacter.json")

	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	data := CharacterNodes{}
	err = json.NewDecoder(jsonFile).Decode(&data)
	for i := 0; i < len(data.CharacterNodes); i++ {
		fmt.Println("Name: ", data.CharacterNodes[i].Name)
		fmt.Println("Status: ", data.CharacterNodes[i].Status)
	}
	http.Handle("/foo", &fooHandler{Message: "hello World"})
	http.HandleFunc("/bar", barHandler)
	http.ListenAndServe(":5000", nil)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
}
