package main

import (
	"fmt"
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

func convertSWAPIPersonToPerson(swapiPerson SWAPITechPerson) (personDTO, error) {
	height, heightConversionError := numericStringOrUnknownToFloatOrNil(swapiPerson.Height)

	id, idConversionError := getResourceIDFromURL(swapiPerson.URL)

	var planetID any
	homeworld, homeworldConversionError := getResourceIDFromURL(swapiPerson.Homeworld)

	mass, massConvError := numericStringOrUnknownToFloatOrNil(swapiPerson.Mass)

	errList := []error{heightConversionError, idConversionError, homeworldConversionError, massConvError}
	if hasError(errList) {
		return personDTO{}, getFirstError(errList)
	}

	if planetIsValid(homeworld) {
		planetID = homeworld
	} else {
		planetID = nil
	}

	return personDTO{
		Id:        id,
		Name:      swapiPerson.Name,
		Height:    height,
		Mass:      mass,
		Created:   swapiPerson.Created,
		Edited:    swapiPerson.Edited,
		Homeworld: planetID,
	}, nil
}

func convertSWAPIDevPersonToPerson(swapiPerson SWAPIDevPerson) (personDTO, error) {
	height, heightConversionError := numericStringOrUnknownToFloatOrNil(swapiPerson.Height)

	id, idConversionError := getResourceIDFromURL(swapiPerson.URL)

	var planetID any
	homeworld, homeworldConversionError := getResourceIDFromURL(swapiPerson.Homeworld)

	mass, massConvError := numericStringOrUnknownToFloatOrNil(swapiPerson.Mass)

	errList := []error{heightConversionError, idConversionError, homeworldConversionError, massConvError}
	if hasError(errList) {
		return personDTO{}, getFirstError(errList)
	}

	if planetIsValid(homeworld) {
		planetID = homeworld
	} else {
		planetID = nil
	}

	return personDTO{
		Id:        id,
		Name:      swapiPerson.Name,
		Height:    height,
		Mass:      mass,
		Created:   swapiPerson.Created,
		Edited:    swapiPerson.Edited,
		Homeworld: planetID,
	}, nil
}

var personCache []personDTO

// Attempts to fetch all people from swapi.dev first, then swapi.tech.
// Them maps them to personDTO
// caches results in personCache
func getAllPeople() ([]personDTO, error) {
	if len(personCache) != 0 {
		return personCache, nil
	}

	SWAPIDevResults, SWAPIDevError := getAllFromSWAPIDev[SWAPIDevPerson]("people")

	if SWAPIDevError != nil {
		fmt.Println("[people.service getAllPeople]", "SWAPIDevError", SWAPIDevError)
		SWAPITechResults, SWAPITechError := SWAPITech_getAll[SWAPITechPerson]("people")
		if SWAPITechError != nil {
			return personCache, SWAPITechError
		}
		people, conversionError := convertMany(SWAPITechResults, convertSWAPIPersonToPerson)
		if conversionError != nil {
			fmt.Println("[people.service getAllPeople]", "swapiTech people conversion error", conversionError)
			return personCache, conversionError
		}
		personCache = people
	} else {
		people, conversionError := convertMany(SWAPIDevResults, convertSWAPIDevPersonToPerson)
		if conversionError != nil {
			fmt.Println("[people.service getAllPeople]", "swapiDev people conversion error", conversionError)
			return personCache, conversionError
		}
		personCache = people
	}

	return personCache, nil
}
