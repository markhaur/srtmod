package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Configuration struct {
	Files []FileConfiguration
}

func (c *Configuration) GetFiles() []FileConfiguration {
	return c.Files
}

type FileConfiguration struct {
	InputFile   string
	OutputFile  string
	ModifyValue time.Duration
}

func FromPath(path string) (*Configuration, error) {
	dir, file := filepath.Split(path)
	viper.SetConfigName(file)
	viper.AddConfigPath(dir)
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var configuration Configuration

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		return nil, err
	}

	err := viper.Unmarshal(&configuration)

	if err != nil {
		fmt.Printf("Unable to decode configurations into struct, %v", err)
		return nil, err
	}

	return &configuration, nil
}
