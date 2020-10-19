package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
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

	yf := flag.String("yaml", "", "Path to a yaml file to load from ")
	jf := flag.String("json", "", "Path to a json file to load urls from")
	df := flag.String("bolt", "", "Path to a bolt database to load from")
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

	if *df != "" {
		
		err := populateBolt(*df)
		if err != nil { 
			fmt.Fprintf(os.Stderr, "Failed to populate DB: %v\n", err)
			os.Exit(1)
		}
		
		o, err := loadBolt(*df)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load from DB: %v\n", err)
			os.Exit(1)
		}
		for _, m := range o {
			mappings[m["path"]] = m["url"]
		}
	}


	registerMappings(mappings, mux)
	
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func loadBolt(filename string) ([]map[string]string, error) {

	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}
	

	m := make([]map[string]string, 0)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("mappings"))
		c := b.Cursor()
		
		for k, v := c.First(); k != nil; k, v = c.Next() {
			m = append(m, map[string]string{
				"path": string(k),
				"url": string(v),
			})
			fmt.Println(k, v)
		}
		return nil
	})
	defer db.Close()

	return m, err
}

func populateBolt(filename string) error {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("mappings"))
		if err == nil {
			err := b.Put([]byte("/test-link-1"), []byte("https://google.com"))
			if err != nil {
				return err
			}

			err = b.Put([]byte("/test-link-2"), []byte("https://gog.com"))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failed to initialize database: %s", err)
	}

	return nil
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
