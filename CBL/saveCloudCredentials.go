package CBL

import (
	"errors"
	"github.com/devlife20/monitoring-tool/types"
	"github.com/spf13/viper"
	"log"
)

type CloudLogConfigurations struct {
	CloudLogCredentials types.Credentials
	Schedule            string
}

func SaveCloudConfiguration(cloudConfig CloudLogConfigurations, cloudServiceProvider string) {
	// Read existing config first
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Fatalf("failed to read config: %v", err)
		}
	}

	switch cloudServiceProvider {
	case "AWS":
		viper.Set("aws", struct {
			AccessKey    string `mapstructure:"access_key"`
			AccessSecret string `mapstructure:"access_secret"`
			Region       string `mapstructure:"region"`
			Schedule     string `mapstructure:"schedule"`
		}{
			AccessKey:    cloudConfig.CloudLogCredentials.AccessKeyID,
			AccessSecret: cloudConfig.CloudLogCredentials.SecretAccessKey,
			Region:       cloudConfig.CloudLogCredentials.Region,
			Schedule:     cloudConfig.Schedule,
		})
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("failed to save credentials: %v", err)
		}
	}
}
