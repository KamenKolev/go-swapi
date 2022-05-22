package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func addHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
}

// Would return 2 for a URL such as "https://swapi.dev/api/planets/2/"
func getResourceIDFromURL(url string) (int, error) {
	urlSplit := strings.Split(url, "/")
	id := urlSplit[len(urlSplit)-2]
	intID, stringConversionError := strconv.Atoi(id)
	return intID, stringConversionError
}

func numericStringOrUnknownToFloatOrNil(s string) (any, error) {
	if s == "unknown" {
		return nil, nil
	}

	// The API uses commas to signify thousands. They don't play well with ParseFloat
	s = strings.ReplaceAll(s, ",", "")

	number, convError := strconv.ParseFloat(s, 64)
	if convError != nil {
		return nil, convError
	}

	return number, nil
}

type swapiMultipleResourcesResponse[T any] struct {
	Count    int `json:"count"`
	Next     any `json:"next"`     // nil | string
	Previous any `json:"previous"` // nil | string
	Results  []T `json:"results"`
}

func getAllFromSwapi[T any](resourceName string, requestFn func(page int, resourceName string) (swapiMultipleResourcesResponse[T], error)) []T {
	firstRes, initialGetError := requestFn(1, resourceName)
	if initialGetError != nil {
		fmt.Println("initial get error for", resourceName)
		fmt.Println(initialGetError)
	}

	results := make([]T, firstRes.Count)
	copy(results, firstRes.Results)

	pages := int(math.Ceil(float64(firstRes.Count) / 10))

	for page := 2; page <= pages; page++ {
		res, error := requestFn(page, resourceName)
		if error != nil {
			fmt.Println(error)
		}
		for i, v := range res.Results {
			results[i+page*10-10] = v
		}
	}

	return results
}

func getFromPage[T any](page int, resourceName string) (swapiMultipleResourcesResponse[T], error) {
	url := strings.Join([]string{apiDomain, resourceName, "?page=", strconv.Itoa(page)}, "")
	resp, err := http.Get(url)

	if err != nil {
		// TODO infinite retry could totally backfire
		fmt.Println("Failed getFromPage for", resourceName, "from", url)

		// Toggles requests between the two domains
		// switchUsedDomain()

		return getFromPage[T](page, resourceName)
	} else {
		body, readingError := ioutil.ReadAll(resp.Body)
		var unmarshalled swapiMultipleResourcesResponse[T]

		if readingError != nil {
			fmt.Println("readingError error thrown")
			fmt.Println(readingError)
			return unmarshalled, readingError
		}

		unmarshallingError := json.Unmarshal(body, &unmarshalled)

		if unmarshallingError != nil {
			fmt.Println("unmarshalling error thrown")
			fmt.Println(unmarshallingError)
			return unmarshalled, unmarshallingError
		}

		return unmarshalled, nil
	}
}

func convertMany[I any, O any](inputs []I, converter func(I) (O, error)) ([]O, error) {
	results := make([]O, len(inputs))
	for i, v := range inputs {
		output, err := converter(v)
		if err != nil {
			fmt.Println("Person Conversion error thrown", err, v)
			return results, err
		}
		results[i] = output
	}

	return results, nil
}

func marshal[T any](data []T) ([]byte, error) {
	resultsJSON, marshallingError := json.Marshal(data)

	if marshallingError != nil {
		return nil, nil
	} else {
		return resultsJSON, nil
	}
}
