package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	path := flag.String("path", "", "Path to html file")
	flag.Parse()

	if *path == "" {
		flag.Usage()
	}

	bb, err := ioutil.ReadFile(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v", *path, err)
		os.Exit(1)
	}

	fmt.Println(string(bb))
}
