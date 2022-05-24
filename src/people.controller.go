package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handles GET requests to /, returns all items
func peopleController(writer http.ResponseWriter, req *http.Request) {
	people, error := getAllPeople()
	if error != nil {
		fmt.Println("[people.controller]", "Could not fetch people", error)
	}

	json, marshallingError := json.Marshal(people)
	if marshallingError != nil {
		fmt.Println("[people.controller]", "people marshalling error", marshallingError)
	}

	if error != nil || marshallingError != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	addHeaders(writer)
	writer.Write(json)
}
