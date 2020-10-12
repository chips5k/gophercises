package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
)

func main() {

	f, err := os.Open("problems.csv")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open problems.csv: %v", err)
		os.Exit(1)
	}

	r := csv.NewReader(f)

	d, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse problems.csv: %v", err)
		os.Exit(1)
	}

	c := 0
	s := bufio.NewScanner(os.Stdin)
	for _, q := range d {
		fmt.Print(q[0] + "? ")
		s.Scan()
		if s.Text() == q[1] {
			c++
		}
	}

	fmt.Printf("%d/%d\n", c, len(d))

	os.Exit(0)
}
