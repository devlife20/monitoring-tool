package utilities

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// CreateConfigPath creates the config directory if it doesn't exist and adds it to Viper's config path
func CreateConfigPath() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".config", "monit")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := SetDefaultConfig(); err != nil {
		return err
	}

	// viper.SetDefault("key", "value")

	// Create config file if it doesn't exist
	if err := viper.SafeWriteConfig(); err != nil {
		// If the error is that the config file already exists, that's okay
		var configFileAlreadyExistsError viper.ConfigFileAlreadyExistsError
		if !errors.As(err, &configFileAlreadyExistsError) {
			return err
		}
	}

	return nil
}

func SetDefaultConfig() error {
	// Set default values for AWS
	viper.SetDefault("aws.access_key", "your-access-key")
	viper.SetDefault("aws.access_secret", "your-secret")
	viper.SetDefault("aws.region", "your-region")

	// Set default values for Azure
	viper.SetDefault("azure.tenantID", "your-tenant-id")
	viper.SetDefault("azure.clientID", "your-client-id")
	viper.SetDefault("azure.clientSecret", "your-client-secret")

	// Set default values for GCP
	viper.SetDefault("gcp.projectID", "your-project-id")
	viper.SetDefault("gcp.keyFilePath", "your-key-file-path")

	// Set default values for ELK
	viper.SetDefault("elk.api_token", "your-api-token")
	viper.SetDefault("elk.elastic_url", "your-elastic-url")

	// Create the config file with default values if it doesn't exist
	if err := viper.SafeWriteConfig(); err != nil {
		// Ignore error if config file already exists
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			return fmt.Errorf("failed to write default config: %v", err)
		}
	}

	return nil
}
