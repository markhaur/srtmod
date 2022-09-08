package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configuration struct {
	Files []FileConfiguration
}

func (c *Configuration) GetFiles() []FileConfiguration {
	return c.Files
}

type FileConfiguration struct {
	FilePath    string
	ModifyValue int
}

func GetConfiguration(filename string, filepath string) (*Configuration, error) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(filepath)
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
