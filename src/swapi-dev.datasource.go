package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// generated with https://mholt.github.io/json-to-go/
type SWAPIDevMultiResourceResponse[T any] struct {
	Count    int `json:"count"`
	Next     any `json:"next"`     // nil | string
	Previous any `json:"previous"` // nil | string
	Results  []T `json:"results"`
}
type SWAPIDevPeopleResponse = SWAPIDevMultiResourceResponse[SWAPIDevPerson]
type SWAPIDevPlanetsResponse = SWAPIDevMultiResourceResponse[SWAPIDevPlanet]

type SWAPIDevPerson struct {
	BirthYear string    `json:"birth_year"`
	Created   time.Time `json:"created"`
	Edited    time.Time `json:"edited"`
	EyeColor  string    `json:"eye_color"`
	Films     []string  `json:"films"`
	Gender    string    `json:"gender"`
	HairColor string    `json:"hair_color"`
	Height    string    `json:"height"`
	Homeworld string    `json:"homeworld"`
	Mass      string    `json:"mass"`
	Name      string    `json:"name"`
	SkinColor string    `json:"skin_color"`
	Species   []string  `json:"species"`
	Starships []string  `json:"starships"`
	URL       string    `json:"url"`
	Vehicles  []string  `json:"vehicles"`
}

type SWAPIDevPlanet struct {
	Climate        string    `json:"climate"`
	Created        time.Time `json:"created"`
	Diameter       string    `json:"diameter"`
	Edited         time.Time `json:"edited"`
	Films          []string  `json:"films"`
	Gravity        string    `json:"gravity"`
	Name           string    `json:"name"`
	OrbitalPeriod  string    `json:"orbital_period"`
	Population     string    `json:"population"`
	Residents      []string  `json:"residents"`
	RotationPeriod string    `json:"rotation_period"`
	SurfaceWater   string    `json:"surface_water"`
	Terrain        string    `json:"terrain"`
	URL            string    `json:"url"`
}

const SWAPIDevAPIURL = "https://swapi.dev/api/"

// resourceName can be "people", "planets" ...
func getAllFromSWAPIDev[T any](resourceName string) ([]T, error) {
	requestFn := getFromPage[T]
	firstRes, initialGetError := requestFn(1, resourceName)
	results := []T{}

	if initialGetError != nil {
		fmt.Println("[swapi-dev.datasource getAllFromSWAPIDev]", "initial get error for", resourceName, initialGetError)
		return results, initialGetError
	}

	results = append(results, firstRes.Results...)

	pages := int(math.Ceil(float64(firstRes.Count) / 10))

	for page := 2; page <= pages; page++ {
		res, error := requestFn(page, resourceName)
		if error != nil {
			fmt.Println("[swapi-dev.datasource getAllFromSWAPIDev]", "request failed for page", page, "resource", resourceName)
			return results, error
		}

		results = append(results, res.Results...)
	}

	return results, nil
}

func getFromPage[T any](page int, resourceName string) (SWAPIDevMultiResourceResponse[T], error) {
	url := strings.Join([]string{SWAPIDevAPIURL, resourceName, "?page=", strconv.Itoa(page)}, "")
	resp, err := http.Get(url)
	var results SWAPIDevMultiResourceResponse[T]

	if err != nil {
		fmt.Println("[swapi-dev.datasource getFromPage]", "Failed for", resourceName, "from", url)

		return results, err
	} else {
		body, readingError := ioutil.ReadAll(resp.Body)

		if readingError != nil {
			fmt.Println("[swapi-dev.datasource getFromPage]", "readingError error thrown", readingError)
			return results, readingError
		}

		unmarshallingError := json.Unmarshal(body, &results)

		if unmarshallingError != nil {
			fmt.Println("[swapi-dev.datasource getFromPage]", "unmarshalling error thrown", unmarshallingError)
			return results, unmarshallingError
		}

		return results, nil
	}
}
