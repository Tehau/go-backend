package main

import (
	"backend/character"
	"github.com/gorilla/mux"
	"net/http"
)

const apiBasePath = "/api/v1"

func main() {
	r := mux.NewRouter()
	character.SetupRoutes(apiBasePath, r)
	http.ListenAndServe(":5000", r)
}
