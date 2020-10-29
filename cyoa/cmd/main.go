package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chips5k/gophercises/cyoa"
)

func main() {

	port := flag.String("port", "3000", "port of the webserver")
	path := flag.String("json", "gopher.json", "path of a json story file")
	start := flag.String("start", "intro", "key of the first story arc")
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s, error: %v", *path, err)
		os.Exit(1)
	}

	d := json.NewDecoder(f)

	s := make(cyoa.Story)
	err = d.Decode(&s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal json: %v", err)
		os.Exit(1)
	}

	h := cyoa.NewStoryHandler(s, *start)
	log.Fatal(http.ListenAndServe(":"+*port, h))
}
