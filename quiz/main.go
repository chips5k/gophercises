package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	duration := flag.Int("duration", 30, "time limit in seconds to complete the quiz")
	filename := flag.String("filename", "problems.csv", "the filename/path of a csv in the format \"question,answer\"")
	shuffle := flag.Bool("shuffle", false, "Randomizes the order of questions")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s: %v", *filename, err)
		os.Exit(1)
	}

	r := csv.NewReader(f)

	questions, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse %s: %v", *filename, err)
		os.Exit(1)
	}
	
	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	fmt.Print("Press enter to begin")
	fmt.Scan()
	result := make(chan bool, len(questions))
	done := make(chan bool)
	timer := time.NewTimer(time.Duration(*duration) * time.Second)
	go quiz(questions, result, done)
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

	fmt.Printf("\nResults: %d/%d\n", count, len(questions))

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
