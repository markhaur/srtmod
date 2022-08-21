package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {

	// if len(os.Args) < 3 {
	// 	fmt.Println("Usage: go run modifier.go <inputfile> <outputfile>")
	// }

	// inputFile := os.Args[1]
	// outputFile := os.Args[2]

	// fmt.Printf("input file: %v\n", inputFile)
	// fmt.Printf("output file: %v\n", outputFile)

	f, err := os.Open("./input/Hz.Omer.S01E01.srt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	subtitleCounter := 0

	for scanner.Scan() {
		// fmt.Println(scanner.Text())

		if subtitleCounter%4 == 1 {
			fmt.Printf("time -> %v\n", scanner.Text())
			// TODO: split the time into start time and end time

			// TODO: subtract time from start time

			// TODO: subtract time from end time

			// TODO: combines the time
		}
		// TODO: write processed information into file

		subtitleCounter += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
