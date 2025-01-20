/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/devlife20/monitoring-tool/cmd"
	utilities "github.com/devlife20/monitoring-tool/utilies"
	"log"
)

func main() {
	err := utilities.CreateConfigPath()
	if err != nil {
		log.Fatalf("Error creating config path: %v", err)
	}
	cmd.Execute()
}
