package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ElasticAPIKey string `mapstructure:"elastic_api_key"`
}

func SetConfig() {
	// Read existing config first
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Fatalf("failed to read config: %v", err)
		}
	}
	// Interactive prompt for configuration
	fmt.Print("Enter Elasticsearch API Key: ")
	var apiKey string
	fmt.Scanln(&apiKey)

	fmt.Print("Enter Elasticsearch URL (default: http://localhost:9200): ")
	var elasticURL string
	fmt.Scanln(&elasticURL)
	if elasticURL == "" {
		elasticURL = "http://localhost:9200"
	}

	// Set config values
	viper.Set("elk.api_token", apiKey)
	viper.Set("elk.elastic_url", elasticURL)

	// Save config
	err := viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}
