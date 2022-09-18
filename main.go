package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Files []struct {
		InputFile  string        `yaml:"inputFile"`
		OutputFile string        `yaml:"outputFile"`
		Offset     time.Duration `yaml:"offset"`
	} `yaml:"files"`
}

var regexPattern = regexp.MustCompile("[0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3} --> [0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3}")

const timeFormat = "15:04:05,000"

func main() {
	var configPath string
	var duration time.Duration
	var inputPath string
	var outputPath string
	flag.StringVar(&configPath, "config", "./config.yml", "Path to configuration file.")
	flag.DurationVar(&duration, "duration", 0, "Time modify duration. Must be negative for now.")
	flag.StringVar(&inputPath, "i", "", "Path to input file.")
	flag.StringVar(&outputPath, "o", "", "Path to output file.")
	flag.Parse()

	if inputPath != "" && outputPath != "" && duration < 0 {
		err := process(inputPath, outputPath, duration)
		if err != nil {
			log.Fatalf("could not process file %s: %v", inputPath, err)
		}
		os.Exit(0)
	}

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("could not open file from path %s: %v\n", configPath, err)
	}
	defer f.Close()

	var config Configuration
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		log.Fatalf("could not decode file contents from path %s: %v\n", configPath, err)
	}

	var wg sync.WaitGroup
	for _, cf := range config.Files {
		wg.Add(1)
		go func(in, out string, offset time.Duration) {
			defer wg.Done()
			err := process(in, out, offset)
			if err != nil {
				log.Printf("could not process file %s: %v\n", in, err)
			}
		}(cf.InputFile, cf.OutputFile, cf.Offset)
	}
	wg.Wait()
}

func process(inputPath string, outputPath string, offset time.Duration) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return errors.Wrap(err, "could not open input file")
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "could not create output file")
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		text := scanner.Text()

		if regexPattern.MatchString(text) {
			times := strings.Split(text, " --> ")

			startTime, err := time.Parse(timeFormat, times[0])
			if err != nil {
				return errors.Wrap(err, "could not parse time")
			}

			endTime, err := time.Parse(timeFormat, times[1])
			if err != nil {
				return errors.Wrap(err, "could not parse time")
			}

			startTime = startTime.Add(offset)
			endTime = endTime.Add(offset)

			text = fmt.Sprintf("%s --> %s", startTime.Format(timeFormat), endTime.Format(timeFormat))
		}
		outFile.WriteString(text + "\n")
	}

	log.Printf("Processed %v\n", inputPath)
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "could not scan")
	}

	return nil
}
