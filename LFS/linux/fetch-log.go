package linux

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func FetchLogsFromFile(filePath, pattern string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	lines := strings.Split(string(file), "\n")
	filteredLines := filterLogs(lines, pattern)

	// If no logs match, print a message
	if len(filteredLines) == 0 {
		fmt.Println("No logs found matching the pattern.")
		return
	}

	// Print filtered log lines
	for _, line := range filteredLines {
		fmt.Println(line)
	}
}

// filterLogs applies filtering based on pattern.
func filterLogs(lines []string, pattern string) []string {
	var filteredLines []string
	for _, line := range lines {
		// Skip empty lines
		if line == "" {
			continue
		}

		// Apply pattern filter
		if pattern != "" && !strings.Contains(line, pattern) {
			continue
		}

		// Add matching line to filteredLines
		filteredLines = append(filteredLines, line)
	}

	return filteredLines
}
