package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/markhaur/srt-time-modifier/config"
	"github.com/pkg/errors"
)

var regexPattern = regexp.MustCompile("[0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3} --> [0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3}")

func main() {
	var configPath string
	var duration time.Duration
	var inputFile string
	var outputFile string
	flag.StringVar(&configPath, "config", "./config.yml", "Path to configuration file.")
	flag.DurationVar(&duration, "duration", 0, "Time modify duration. Must be negative for now.")
	flag.StringVar(&inputFile, "i", "", "Path to input file.")
	flag.StringVar(&outputFile, "o", "", "Path to output file.")
	flag.Parse()

	if inputFile != "" && outputFile != "" && duration < 0 {
		err := process(inputFile, outputFile, duration)
		if err != nil {
			log.Fatalf("could not process file %s: %v", inputFile, err)
		}
		os.Exit(0)
	}

	configuration, err := config.GetConfiguration(configPath)
	if err != nil {
		log.Fatalf("could not load config from path %s: %v\n", configPath, err)
	}

	var wg sync.WaitGroup
	for _, cf := range configuration.Files {
		wg.Add(1)
		go func(in, out string, offset time.Duration) {
			defer wg.Done()
			err := process(iFile, oFile, mVal)
			if err != nil {
				log.Printf("could not process file %s: %v\n", iFile, err)
			}
		}(cf.InputFile, cf.OutputFile, cf.ModifyValue)
	}
	wg.Wait()
}

func process(inputPath string, outputPath string, offset time.Duration) error {
	inpFile, err := os.Open(inputFile)
	if err != nil {
		return errors.Wrap(err, "could not open input file")
	}
	defer inpFile.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		return errors.Wrap(err, "could not create output file")
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inpFile)
	for scanner.Scan() {
		text := scanner.Text()

		if regexPattern.MatchString(text) {
			times := strings.Split(text, " ")
			startTime := times[0]
			endTime := times[2]

			modifiedStartTime := addValueToTime(startTime, modifyValue)
			modifiedEndTime := addValueToTime(endTime, modifyValue)
			text = fmt.Sprintf("%s --> %s", modifiedStartTime, modifiedEndTime)
		}
		outFile.WriteString(text + "\n")
	}

	fmt.Printf("Processed %v\n", inputFile)
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "could not scan")
	}

	return nil

}

// currently supports subtraction of time
func addValueToTime(time string, value time.Duration) string {
	parts := strings.Split(time, ",")

	timeParts := strings.Split(parts[0], ":")

	sec, _ := strconv.Atoi(timeParts[2])
	min, _ := strconv.Atoi(timeParts[1])
	hrs, _ := strconv.Atoi(timeParts[0])

	sec += int(value.Seconds())

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

	return fmt.Sprintf("%s:%s:%s,%s", timeParts[0], timeParts[1], timeParts[2], parts[1])
}
