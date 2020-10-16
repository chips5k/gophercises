package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to urlshort!")
	})

	mappings := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	yf := flag.String("YAML", "", "Where to load YAML data")
	jf := flag.String("JSON", "", "Where to load JSON data")
	flag.Parse()

	if *yf != "" {
		o, err := loadYAML(*yf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load YAML: %v\n", err)
			os.Exit(1)
		}

		for _, m := range o {
			mappings[m["path"]] = m["url"]
		}
	}

	if *jf != "" {
		o, err := loadJSON(*jf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load JSON: %v\n", err)
			os.Exit(1)
		}

		for _, m := range o {
			mappings[m["path"]] = m["url"]
		}
	}


	registerMappings(mappings, mux)
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func loadYAML(filename string) ([]map[string]string, error) {

	var o []map[string]string

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	err = yaml.Unmarshal(f, &o)
	if err != nil {
		return nil, err
	}
	return o, nil

}


func loadJSON(filename string) ([]map[string]string, error) {

	var o []map[string]string

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(f, &o)
	if err != nil {
		return nil, err
	}
	return o, nil

}

func registerMappings(mappings map[string]string, mux *http.ServeMux) {
	for k, v := range mappings {
		url := v
		mux.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, url, 308)
		})
	}
}
