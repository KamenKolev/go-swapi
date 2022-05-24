package main

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/exp/slices"
)

type planetDTO struct {
	Climate    any    `json:"climate"`
	Diameter   any    `json:"diameter"` // float or nil
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Population any    `json:"population"` // float or nil
}

// swapiPlanet works for both SWAPIDevPlanet and SWAPITechPlanet
func convertSWAPIDevPlanetToPlanet(swapiPlanet SWAPIDevPlanet) (planetDTO, error) {
	var result planetDTO
	if swapiPlanet.Name == "unknown" {
		return result, errors.New("unknown planet name")
	}

	diameter, diameterConvError := numericStringOrUnknownToFloatOrNil(swapiPlanet.Diameter)

	population, populationConvError := numericStringOrUnknownToFloatOrNil(swapiPlanet.Population)

	id, idConversionError := getResourceIDFromURL(swapiPlanet.URL)

	var climate any
	if swapiPlanet.Climate == "unknown" {
		climate = nil
	} else {
		climate = swapiPlanet.Climate
	}

	errList := []error{diameterConvError, populationConvError, idConversionError}
	if hasError(errList) {
		return result, getFirstError(errList)
	}

	result = planetDTO{
		Id:         id,
		Name:       swapiPlanet.Name,
		Diameter:   diameter,
		Climate:    climate,
		Population: population,
	}

	return result, nil
}

func convertSWAPITechPlanetToPlanet(swapiPlanet SWAPITechPlanet) (planetDTO, error) {
	var result planetDTO
	diameter, diameterConvError := numericStringOrUnknownToFloatOrNil(swapiPlanet.Diameter)
	if diameterConvError != nil {
		return planetDTO{}, diameterConvError
	}
	population, populationConvError := numericStringOrUnknownToFloatOrNil(swapiPlanet.Population)
	if populationConvError != nil {
		return planetDTO{}, populationConvError
	}
	id, idError := getResourceIDFromURL(swapiPlanet.URL)
	if idError != nil {
		return planetDTO{}, idError
	}

	if swapiPlanet.Name == "unknown" {
		return result, errors.New("unknown planet name")
	}

	result = planetDTO{
		Id:         id,
		Name:       swapiPlanet.Name,
		Diameter:   diameter,
		Climate:    swapiPlanet.Climate,
		Population: population,
	}

	return result, nil

}

var planetCache []planetDTO

// caches results in planetCache
func getAllPlanets() ([]planetDTO, error) {
	if len(planetCache) != 0 {
		return planetCache, nil
	}

	SWAPIDevResults, SWAPIDevError := getAllFromSWAPIDev[SWAPIDevPlanet]("planets")

	if SWAPIDevError != nil {
		fmt.Println("[planets.service getAllPlanets]", "SWAPIDevError", SWAPIDevError)
		SWAPITechResults, SWAPITechError := SWAPITech_getAll[SWAPITechPlanet]("planets")
		if SWAPITechError != nil {
			return planetCache, SWAPITechError
		}

		planets, conversionError := convertMany(SWAPITechResults, convertSWAPITechPlanetToPlanet)
		if conversionError != nil {
			log.Fatalln("[planets.service getAllPlanets]", "SWAPITech planets conversion error", conversionError)
		}

		planetCache = planets
	} else {
		planets, conversionError := convertMany(SWAPIDevResults, convertSWAPIDevPlanetToPlanet)
		if conversionError != nil {
			log.Fatalln("[planets.service getAllPlanets]", "SWAPIDev planets conversion error", conversionError)
		}

		planetCache = planets
	}

	return planetCache, nil
}

func planetIsValid(planetID int) bool {
	planets, _ := getAllPlanets()

	planetIndex := slices.IndexFunc(planets, func(planet planetDTO) bool {
		return planet.Id == planetID
	})

	if planetIndex == -1 {
		return false
	}

	planet := planets[planetIndex]

	return planet.Name != "unknown"
}
