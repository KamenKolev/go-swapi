package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Used to tackle overfetching
type planetDTO struct {
	Climate    string `json:"climate"`
	Diameter   any    `json:"diameter"` // float or nil
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Population any    `json:"population"` // float or nil
}

func planetsController() func(http.ResponseWriter, *http.Request) {
	planets, error := getAllPlanets()
	if error != nil {
		log.Fatalln("Could not fetch planets")
	}

	json, marshallingError := json.Marshal(planets)
	if marshallingError != nil {
		log.Fatalln("Planets marshalling error", marshallingError)
	}
	return func(writer http.ResponseWriter, req *http.Request) {
		addHeaders(writer)
		writer.Write(json)
	}
}
