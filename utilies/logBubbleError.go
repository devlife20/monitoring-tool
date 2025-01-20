package utilities

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

// LogBubbleTeaError sets up Bubble Tea logging to /var/log/monit.
func LogBubbleTeaError(message error) {
	// Define the log file path
	logFilePath := "/var/log/monit/debug.log"

	// Ensure the /var/log/monit directory exists
	err := os.MkdirAll("/var/log/monit", 0755)
	if err != nil {
		fmt.Println("fatal: unable to create /var/log/monit directory:", err)
		os.Exit(1)
	}

	// Open the log file for Bubble Tea logging
	f, err := tea.LogToFile(logFilePath, "debug")
	if err != nil {
		fmt.Println("fatal: unable to open log file:", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintln(f, message)
}
