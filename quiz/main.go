package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	duration := flag.Int("duration", 30, "time limit in seconds to complete the quiz")
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
	
	
	fmt.Print("Press enter to begin")
	fmt.Scan()
	result := make(chan bool, len(d))
	done := make(chan bool)
	timer := time.NewTimer(time.Duration(*duration) * time.Second)
	go quiz(d, result, done)
	count := 0
	running := true
	for running {
		select {
			case r := <-result:
				if r { 
					count++
				}
			case <-timer.C: 
				fmt.Print("\nTimes up!\n")
				running = false
			case <-done:
				running = false
		}
	}

	fmt.Printf("\nResults: %d/%d\n", count, len(d))

	os.Exit(0)
}



func quiz(questions [][]string, result chan bool, done chan bool) {
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	for i, q := range questions {
		fmt.Printf("\nQuestion %d: %s? ", i+1, strings.Replace(q[0], "?", "", 1))
		s.Scan()
		result <- (s.Text() == strings.TrimSpace(q[1]))
	}
	done <- true 
}
