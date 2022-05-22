package main

import (
	"log"
	"net/http"
	"time"
)

// Used to tackle overfetching
type planetDTO struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Diameter   any    `json:"diameter"` // float or nil
	Climate    string `json:"climate"`
	Population any    `json:"population"` // float or nil
}

type swapiPlanetDTO struct {
	Name           string    `json:"name"`
	RotationPeriod string    `json:"rotation_period"`
	OrbitalPeriod  string    `json:"orbital_period"`
	Diameter       string    `json:"diameter"`
	Climate        string    `json:"climate"`
	Gravity        string    `json:"gravity"`
	Terrain        string    `json:"terrain"`
	SurfaceWater   string    `json:"surface_water"`
	Population     string    `json:"population"`
	Residents      []string  `json:"residents"`
	Films          []string  `json:"films"`
	Created        time.Time `json:"created"`
	Edited         time.Time `json:"edited"`
	URL            string    `json:"url"`
}

type swapiPlanetsReponse = swapiMultipleResourcesResponse[swapiPlanetDTO]

func swapiPlanetToPlanet(swapiPlanet swapiPlanetDTO) (planetDTO, error) {
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

func handleGetPlanets() func(http.ResponseWriter, *http.Request) {
	swapiDTOs := getAllFromSwapi("planets", getFromPage[swapiPlanetDTO])
	planets, conversionError := convertMany(swapiDTOs, swapiPlanetToPlanet)
	if conversionError != nil {
		log.Fatalln("Planets conversion error", conversionError)
	}
	json, marshallingError := marshal(planets)
	if marshallingError != nil {
		log.Fatalln("Planets marshalling error", marshallingError)
	}
	return func(writer http.ResponseWriter, req *http.Request) {
		addHeaders(writer)
		writer.Write(json)
	}

}
