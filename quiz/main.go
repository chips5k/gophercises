package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	filename := flag.String("filename", "problems.csv", "the filename/path of a csv in the format \"question,answer\"")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s: %v", *filename, err)
		os.Exit(1)
	}

	r := csv.NewReader(f)

	d, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %s: %v", *filename, err)
		os.Exit(1)
	}
	
	c := 0
	s := bufio.NewScanner(os.Stdin)
	for _, q := range d {
		fmt.Print(strings.Replace(q[0], "?", "", 1) + "? ")
		s.Scan()
		if s.Text() == strings.TrimSpace(q[1]) {
			c++
		}
	}

	fmt.Printf("%d/%d\n", c, len(d))

	os.Exit(0)
}

//Test cases
/*
- reads in csv file defaults to problem.csv
- specify filename via flag


*/
