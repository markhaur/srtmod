package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type configuration struct {
	Files []file `yaml:"files"`
}

type file struct {
	InputFile  string        `yaml:"inputFile"`
	OutputFile string        `yaml:"outputFile"`
	Offset     time.Duration `yaml:"offset"`
}

var regexPattern = regexp.MustCompile("[0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3} --> [0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3}")

const timeFormat = "15:04:05,000"

func main() {
	var configPath string
	var offset time.Duration
	var inputPath string
	var outputPath string
	flag.StringVar(&configPath, "config", "config.yml", "Path to configuration file.")
	flag.DurationVar(&offset, "offset", 0, "Time modify duration like 1s, 1m, -1m")
	flag.StringVar(&inputPath, "i", "", "Path to input file.")
	flag.StringVar(&outputPath, "o", "", "Path to output file.")
	flag.Parse()

	var config configuration
	switch {
	case inputPath != "" && outputPath != "":
		config.Files = append(config.Files, file{inputPath, outputPath, offset})
	default:
		f, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("could not open config file from path %s: %v\n", configPath, err)
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(&config); err != nil {
			log.Fatalf("could not decode config file contents from path %s: %v\n", configPath, err)
		}
	}

	var wg sync.WaitGroup
	for _, cf := range config.Files {
		in, err := os.Open(cf.InputFile)
		if err != nil {
			log.Printf("could not open input file at path %s: %v", cf.InputFile, err)
		}
		defer in.Close()

		out, err := os.Create(cf.OutputFile)
		if err != nil {
			log.Printf("could not create output file at path %s: %v", cf.OutputFile, err)
		}
		defer out.Close()

		wg.Add(1)
		go func(in, out *os.File, offset time.Duration) {
			defer wg.Done()
			err := process(in, out, offset)
			if err != nil {
				log.Printf("could not process file %s: %v\n", in.Name(), err)
				return
			}
			log.Printf("Processed %s\n", in.Name())
		}(in, out, cf.Offset)
	}
	wg.Wait()
}

func process(input io.Reader, output io.Writer, offset time.Duration) error {
	scanner := bufio.NewScanner(input)
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
		fmt.Fprintf(output, "%s\n", text)
	}

	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "could not scan")
	}

	return nil
}
