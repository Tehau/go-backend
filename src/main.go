package main

import (
	"backend/character"
	"net/http"
)

const apiBasePath = "/api/v1"

func main() {
	character.SetupRoutes(apiBasePath)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		return
	}
}
