/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	Long:  `Display the current configuration settings (excluding sensitive values).`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				fmt.Println("No config file found. Run 'monit config init' to create one.")
				return
			}
			fmt.Printf("Error reading config: %v\n", err)
			return
		}

		fmt.Println("Current Configuration:")
		fmt.Printf("Elasticsearch URL: %s\n", viper.GetString("elastic_url"))
		fmt.Println("API Key: [HIDDEN]")

	},
}

func init() {
	MonitConfig.AddCommand(viewCmd)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(home, ".config", "monit"))
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
