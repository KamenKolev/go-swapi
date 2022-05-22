package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func addHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
}

// TODO would this also work for swapi tech urls?
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
	if s == "" {
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

func getReadAndUnmarshall[T any](url string) (T, error) {
	resp, err := http.Get(url)
	var result T

	if resp.StatusCode != 200 {
		return result, err
	}

	// wow this is hacky. Totally nil
	if err != nil {
		return result, err
	}

	body, readingError := ioutil.ReadAll(resp.Body)
	if readingError != nil {
		return result, err
	}

	unmarshallingError := json.Unmarshal(body, &result)
	if unmarshallingError != nil {
		return result, err
	}

	return result, nil
}
