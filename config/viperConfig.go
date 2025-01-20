package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ElasticAPIKey string `mapstructure:"elastic_api_key"`
}

func SetConfig() {

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
	viper.Set("elastic_api_key", apiKey)
	viper.Set("elastic_url", elasticURL)

	// Save config
	err := viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}
