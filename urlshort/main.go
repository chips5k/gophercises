package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)

	mappings := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	registerMap(mappings, mux)

	http.ListenAndServe(":8080", mux)

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to urlshort!")
} 

func registerMap(mappings map[string]string, mux *http.ServeMux) {

	for k, v := range mappings {
		mux.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Redirect to <b>%s</b>", v)
		})
	}

}

// func handleYaml() {

// }
