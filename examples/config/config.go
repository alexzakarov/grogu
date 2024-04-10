package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func exportConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	configPath := fmt.Sprintf("config.yaml")

	viper.SetConfigName(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// ParseConfig Parse config file
func ParseConfig() (*Config, error) {
	if err := exportConfig(); err != nil {
		return nil, err
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Printf("Unable to decode config file into struct, %v \n", err)
		return nil, err
	}

	return &c, nil
}
