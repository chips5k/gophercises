package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

var data = []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to urlshort!")
	})

	mappings := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	var o []map[string]string
	err := yaml.Unmarshal(data, &o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal yaml: %v", err)
	}

	for _, m := range o {
		mappings[m["path"]] = m["url"]
	}
	
	registerMappings(mappings, mux)
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func registerMappings(mappings map[string]string, mux *http.ServeMux) {
	for k, v := range mappings {
		url := v
		mux.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, url, 308)
		})
	}
}
