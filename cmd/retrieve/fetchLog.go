/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package retrieve

import (
	"fmt"
	"github.com/spf13/cobra"
)

// FetchLogCmd represents the fetchLog command
var FetchLogCmd = &cobra.Command{
	Use:   "fetch",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fetchLog called")
	},
}
