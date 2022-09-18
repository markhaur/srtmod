package config

import (
	"os"
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

func FromPath(path string) (*Configuration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}
	defer f.Close()

	var config Configuration
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, errors.Wrap(err, "could not decode file contents")
	}

	return &config, nil
}
