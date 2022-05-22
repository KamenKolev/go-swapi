package main

import (
	"net/http"
)

// The API has downtime. Luckily, there's another version of it for which an adapter can be written
const mainDomain = "https://swapi.dev/api/"

// const secondaryDomain = "https://swapi.tech/api/"

var apiDomain string // one of the above
// func switchUsedDomain() {
// 	if apiDomain == mainDomain {
// 		apiDomain = secondaryDomain
// 	} else {
// 		apiDomain = mainDomain
// 	}
// }

func main() {
	apiDomain = mainDomain

	http.HandleFunc("/people", handleGetPeople())
	http.HandleFunc("/planets", handleGetPlanets())
	http.ListenAndServe(":8080", nil)
}
