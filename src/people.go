package main

import (
	"log"
	"net/http"
	"time"
)

// Used to tackle overfetching
type personDTO struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Height    any       `json:"height"` // float or nil
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	Homeworld any       `json:"homeworld"` // the ID only. Could be nil if the planet is unknown. In this case, this is planet 28
	Mass      any       `json:"mass"`      // float or nil
}

type swapiPersonDTO struct {
	Name      string    `json:"name"`
	Height    string    `json:"height"`
	Mass      string    `json:"mass"`
	HairColor string    `json:"hair_color"`
	SkinColor string    `json:"skin_color"`
	EyeColor  string    `json:"eye_color"`
	BirthYear string    `json:"birth_year"`
	Gender    string    `json:"gender"`
	Homeworld string    `json:"homeworld"`
	Films     []string  `json:"films"`
	Species   []string  `json:"species"`
	Vehicles  []string  `json:"vehicles"`
	Starships []string  `json:"starships"`
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	URL       string    `json:"url"`
}

type swapiPeopleReponse = swapiMultipleResourcesResponse[swapiPersonDTO]

func swapiPersonToPerson(swapiPerson swapiPersonDTO) (personDTO, error) {
	height, heightConversionError := numericStringOrUnknownToFloatOrNil(swapiPerson.Height)
	if heightConversionError != nil {
		return personDTO{}, heightConversionError
	}
	id, idConversionError := getResourceIDFromURL(swapiPerson.URL)
	if idConversionError != nil {
		return personDTO{}, idConversionError
	}

	var homeworld any
	homeworld, homeworldConversionError := getResourceIDFromURL(swapiPerson.Homeworld)
	if homeworldConversionError != nil {
		return personDTO{}, homeworldConversionError
	}
	mass, massConvError := numericStringOrUnknownToFloatOrNil(swapiPerson.Mass)
	if massConvError != nil {
		return personDTO{}, massConvError
	}

	// Planet 28 is "unknown" and has zero other useful info. Better return null instead
	if homeworld == 28 {
		homeworld = nil
	}

	return personDTO{
		Id:        id,
		Name:      swapiPerson.Name,
		Height:    height,
		Mass:      mass,
		Created:   swapiPerson.Created,
		Edited:    swapiPerson.Edited,
		Homeworld: homeworld,
	}, nil
}

func handleGetPeople() func(http.ResponseWriter, *http.Request) {
	swapiDTOs := getAllFromSwapi[swapiPersonDTO]("people", getFromPage[swapiPersonDTO])
	people, conversionError := convertMany[swapiPersonDTO, personDTO](swapiDTOs, swapiPersonToPerson)
	if conversionError != nil {
		log.Fatalln("people conversion error", conversionError)
	}
	json, marshallingError := marshal[personDTO](people)
	if marshallingError != nil {
		log.Fatalln("people marshalling error", marshallingError)
	}

	return func(writer http.ResponseWriter, req *http.Request) {
		addHeaders(writer)
		writer.Write(json)
	}
}
