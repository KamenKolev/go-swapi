package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/people", peopleController)
	http.HandleFunc("/planets", planetsController)
	http.ListenAndServe(":8080", nil)
}
