package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// handles GET requests to /, returns all items
func planetsController(writer http.ResponseWriter, req *http.Request) {
	planets, error := getAllPlanets()
	if error != nil {
		log.Fatalln("[planets.controller]", "Could not fetch planets")
	}

	json, marshallingError := json.Marshal(planets)
	if marshallingError != nil {
		log.Fatalln("[planets.controller]", "marshalling error", marshallingError)
	}

	if error != nil || marshallingError != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	addHeaders(writer)
	writer.Write(json)
}
