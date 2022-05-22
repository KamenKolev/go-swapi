package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Used to tackle overfetching
func peopleController() func(http.ResponseWriter, *http.Request) {
	people, error := getAllPeople()
	if error != nil {
		log.Fatalln("Could not fetch people", error)
	}

	json, marshallingError := json.Marshal(people)
	if marshallingError != nil {
		log.Fatalln("people marshalling error", marshallingError)
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		addHeaders(writer)
		writer.Write(json)
	}
}
