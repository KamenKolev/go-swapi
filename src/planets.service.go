package main

import (
	"fmt"
	"log"
)

// swapiPlanet works for both SWAPIDevPlanet and SWAPITechPlanet
func convertSWAPIDevPlanetToPlanet(swapiPlanet SWAPIDevPlanet) (planetDTO, error) {
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

	return planetDTO{
		Id:         id,
		Name:       swapiPlanet.Name,
		Diameter:   diameter,
		Climate:    swapiPlanet.Climate,
		Population: population,
	}, nil
}

// TODO dedupe with convertSWAPIDevPlanetToPlanet
func convertSWAPITechPlanetToPlanet(swapiPlanet SWAPITechPlanet) (planetDTO, error) {
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

	return planetDTO{
		Id:         id,
		Name:       swapiPlanet.Name,
		Diameter:   diameter,
		Climate:    swapiPlanet.Climate,
		Population: population,
	}, nil
}

func getAllPlanets() ([]planetDTO, error) {
	var results []planetDTO

	SWAPIDevResults, SWAPIDevError := getAllFromSWAPIDev[SWAPIDevPlanet]("planets")

	if SWAPIDevError != nil {
		fmt.Println("SWAPIDevError", SWAPIDevError)
		SWAPITechResults, SWAPITechError := SWAPITech_getAll[SWAPITechPlanet]("planets")
		if SWAPITechError != nil {
			return results, SWAPITechError
		}
		planets, conversionError := convertMany(SWAPITechResults, convertSWAPITechPlanetToPlanet)
		if conversionError != nil {
			log.Fatalln("planets conversion error", conversionError)
		}
		return planets, nil
	} else {
		planets, conversionError := convertMany(SWAPIDevResults, convertSWAPIDevPlanetToPlanet)
		if conversionError != nil {
			log.Fatalln("planets conversion error", conversionError)
		}
		return planets, nil
	}
}
