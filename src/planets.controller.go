package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// handles GET requests to /, returns all items
// caches all items in memory and sends them right away
func planetsController() func(http.ResponseWriter, *http.Request) {
	planets, error := getAllPlanets()
	if error != nil {
		log.Fatalln("[planets.controller]", "Could not fetch planets")
	}

	json, marshallingError := json.Marshal(planets)
	if marshallingError != nil {
		log.Fatalln("[planets.controller]", "marshalling error", marshallingError)
	}
	return func(writer http.ResponseWriter, req *http.Request) {
		addHeaders(writer)
		writer.Write(json)
	}
}
