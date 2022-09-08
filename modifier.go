package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/markhaur/srt-time-modifier/config"
)

func main() {
	filename := "config"
	filepath := "."

	configuration, err := config.GetConfiguration(filename, filepath)

	if err != nil {
		fmt.Printf("error reading configs: %v", err)
	}

	for _, cf := range configuration.Files {
		fmt.Printf("filename: %v with modifValue: %v\n", cf.FilePath, cf.ModifyValue)
	}
}

/*
func main() {

	if len(os.Args) < 4 {
		fmt.Println("Usage: go run modifier.go <inputfile> <outputfile> <modify value>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	modifyValue, _ := strconv.Atoi(os.Args[3])

	fmt.Printf("Processing %v\n", inputFile)

	inpFile, err := os.Open(inputFile)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Creating %v if not exists\n", outputFile)
	outFile, err := os.Create(outputFile)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer inpFile.Close()
	defer outFile.Close()

	scanner := bufio.NewScanner(inpFile)

	subtitleCounter := 0
	regexPattern, _ := regexp.Compile("[0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3} --> [0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3}")

	for scanner.Scan() {

		text := scanner.Text()

		if regexPattern.MatchString(text) {
			// split the time into start time and end time
			times := strings.Split(text, " ")
			startTime := times[0]
			endTime := times[2]

			// subtract time from start time
			modifiedStartTime := addValueToTime(startTime, modifyValue)

			// subtract time from end time
			modifiedEndTime := addValueToTime(endTime, modifyValue)

			// combines the time
			text = modifiedStartTime + " --> " + modifiedEndTime
		}
		// write processed information into file
		outFile.WriteString(text + "\n")
		subtitleCounter += 1
	}

	fmt.Printf("Processed %v\n", inputFile)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
*/
// currently supports subtraction of time
func addValueToTime(time string, value int) string {
	parts := strings.Split(time, ",")

	timeParts := strings.Split(parts[0], ":")

	sec, _ := strconv.Atoi(timeParts[2])
	min, _ := strconv.Atoi(timeParts[1])
	hrs, _ := strconv.Atoi(timeParts[0])

	sec += value

	if sec < 0 {
		sec += 60
		min -= 1
		if min < 0 {
			min += 60
			hrs -= 1
		}
	}

	if hrs < 10 {
		timeParts[0] = "0" + strconv.Itoa(hrs)
	} else {
		timeParts[0] = strconv.Itoa(hrs)
	}

	if min < 10 {
		timeParts[1] = "0" + strconv.Itoa(min)
	} else {
		timeParts[1] = strconv.Itoa(min)
	}

	if sec < 10 {
		timeParts[2] = "0" + strconv.Itoa(sec)
	} else {
		timeParts[2] = strconv.Itoa(sec)
	}

	return timeParts[0] + ":" + timeParts[1] + ":" + timeParts[2] + "," + parts[1]
}
