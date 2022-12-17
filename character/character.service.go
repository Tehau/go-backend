package character

import (
	"backend/cors"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const characterBasePath = "character"

func SetupRoutes(apiBasePath string) {
	characterListHandler := http.HandlerFunc(charactersHandler)
	characterItemHandler := http.HandlerFunc(characterHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, characterBasePath), cors.Middleware(characterListHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, characterBasePath), cors.Middleware(characterItemHandler))
}

func charactersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		characterList := getCharacterList()
		_characterJson, err := json.MarshalIndent(characterList, "", "    ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(_characterJson)

	}
}

func characterHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, characterBasePath+"/")
	charID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	character := getCharacter(charID)

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
		_, err = addOrUpdateCharacter(*character)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeCharacter(charID)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
