package main

import (
	"fmt"
	"strings"
	"time"
)

// generated with https://mholt.github.io/json-to-go/
type SWAPITechResponse struct {
	Message      string                 `json:"message"`
	Next         any                    `json:"next"`
	Previous     any                    `json:"previous"`
	Results      []SWAPITechResResource `json:"results"`
	TotalPages   int                    `json:"total_pages"`
	TotalRecords int                    `json:"total_records"`
}

type SWAPITechResResource struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
	URL  string `json:"url"`
}

type SWAPITechSingleResourceResponse[T any] struct {
	Message string                     `json:"message"`
	Result  SWAPITechResourceResult[T] `json:"result"`
}
type SWAPITechPerson struct {
	BirthYear string    `json:"birth_year"`
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	EyeColor  string    `json:"eye_color"`
	Gender    string    `json:"gender"`
	HairColor string    `json:"hair_color"`
	Height    string    `json:"height"`
	Homeworld string    `json:"homeworld"`
	Mass      string    `json:"mass"`
	Name      string    `json:"name"`
	SkinColor string    `json:"skin_color"`
	URL       string    `json:"url"`
}
type SWAPITechResourceResult[T any] struct {
	Description string `json:"description"`
	ID          string `json:"_id"`
	Properties  T      `json:"properties"`
	UID         string `json:"uid"`
	V           int    `json:"__v"`
}

type SWAPITechPlanet struct {
	Climate        string    `json:"climate"`
	Created        time.Time `json:"created"`
	Diameter       string    `json:"diameter"`
	Edited         time.Time `json:"edited"`
	Gravity        string    `json:"gravity"`
	Name           string    `json:"name"`
	OrbitalPeriod  string    `json:"orbital_period"`
	Population     string    `json:"population"`
	RotationPeriod string    `json:"rotation_period"`
	SurfaceWater   string    `json:"surface_water"`
	Terrain        string    `json:"terrain"`
	URL            string    `json:"url"`
}

const SWAPITechAPIURL = "https://swapi.tech/api/"

// resourceName is  should be in plural (people / planets)
func SWAPITech_getAll[T any](resourceName string) ([]T, error) {
	url := strings.Join([]string{SWAPITechAPIURL, resourceName, "?page=1&limit=100000"}, "")
	initalResponse, err := getReadAndUnmarshall[SWAPITechResponse](url)

	results := []T{}

	if err != nil {
		return results, err
	}

	for _, v := range initalResponse.Results {
		item, error := getReadAndUnmarshall[SWAPITechSingleResourceResponse[T]](v.URL)
		if error != nil {
			fmt.Println("[swapi-tech.datasource SWAPITech_getAll]", "read or unmarshal error for", resourceName, v, error)
			return results, nil
		}

		results = append(results, item.Result.Properties)
	}

	return results, nil
}
