package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// handles GET requests to /, returns all items
// caches all items in memory and sends them right away
func peopleController() func(http.ResponseWriter, *http.Request) {
	people, error := getAllPeople()
	if error != nil {
		log.Fatalln("[people.controller]", "Could not fetch people", error)
	}

	json, marshallingError := json.Marshal(people)
	if marshallingError != nil {
		log.Fatalln("[people.controller]", "people marshalling error", marshallingError)
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		addHeaders(writer)
		writer.Write(json)
	}
}
