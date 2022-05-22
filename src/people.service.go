package main

import (
	"fmt"
	"log"
	"time"
)

type personDTO struct {
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	Height    any       `json:"height"`    // float or nil
	Homeworld any       `json:"homeworld"` // the ID only. Could be nil if the planet is unknown. In this case, this is planet 28
	Id        int       `json:"id"`
	Mass      any       `json:"mass"` // float or nil
	Name      string    `json:"name"`
}

// TODO utilize constraints to dedupe
// type APIPersonDTO interface {
//   SWAPIDevPerson | SWAPITechPerson
// }

// TODO dedupe with swapi dev function
func convertSWAPIPersonToPerson(swapiPerson SWAPITechPerson) (personDTO, error) {
	fmt.Println("CONVERTING", swapiPerson)
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

// TODO dedupe
func convertSWAPIDevPersonToPerson(swapiPerson SWAPIDevPerson) (personDTO, error) {
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

	// TODO check if this is also the case in swapi.tech
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

// Attempts to fetch all people from swapi.dev first, then swapi.tech.
// Them maps them to personDTO
func getAllPeople() ([]personDTO, error) {
	var results []personDTO
	SWAPIDevResults, SWAPIDevError := getAllFromSWAPIDev[SWAPIDevPerson]("people")

	if SWAPIDevError != nil {
		fmt.Println("SWAPIDevError", SWAPIDevError)
		SWAPITechResults, SWAPITechError := SWAPITech_getAll[SWAPITechPerson]("people")
		if SWAPITechError != nil {
			return results, SWAPITechError
		}
		fmt.Println("SWAPITechResults", SWAPITechResults)
		people, conversionError := convertMany(SWAPITechResults, convertSWAPIPersonToPerson)
		if conversionError != nil {
			log.Fatalln("people conversion error", conversionError)
		}
		return people, nil
	} else {
		people, conversionError := convertMany(SWAPIDevResults, convertSWAPIDevPersonToPerson)
		if conversionError != nil {
			log.Fatalln("people conversion error", conversionError)
		}
		return people, nil
	}
}
