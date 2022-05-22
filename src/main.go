package main

import (
	"net/http"
)

// The API has downtime. Luckily, there's another copy of it
const mainDomain = "https://swapi.dev/"
const secondaryDomain = "https://swapi.tech/"

var usedDomain string // on of the above
func switchUsedDomain() {
	if usedDomain == mainDomain {
		usedDomain = secondaryDomain
	} else {
		usedDomain = mainDomain
	}
}

func main() {
	http.HandleFunc("/people", handleGetPeople())
	http.HandleFunc("/planets", handleGetPlanets())
	http.ListenAndServe(":8080", nil)
}
