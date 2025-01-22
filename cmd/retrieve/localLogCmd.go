package retrieve

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var (
	logFilePath string
	logDirPath  string
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Retrieve logs from local files or directories",
	Long: `Fetch logs from local files or directories with various filtering options.

Examples:
 monit local --file /var/log/myapp.log --pattern "ERROR"
 monit local --dir /var/log/ --pattern "WARN" --limit 50
 monit local --file /var/log/myapp.log --start 2024-12-10T15:00:00Z --end 2024-12-10T16:00:00Z`,

	Run: func(cmd *cobra.Command, args []string) {
		// Validate input
		if logFilePath == "" && logDirPath == "" {
			log.Fatal("Either --file or --dir must be specified")
		}
		// Fetch logs
		if logFilePath != "" {
			fetchLogsFromFile(logFilePath, filterPattern)
		} else if logDirPath != "" {
			fetchLogsFromDir(logDirPath, filterPattern)
		}
	},
}

func init() {
	FetchLogCmd.AddCommand(localCmd)

	// Flags for the localCmd
	localCmd.Flags().StringVarP(&logFilePath, "file", "f", "", "Path to the log file (optional)")
	localCmd.Flags().StringVarP(&logDirPath, "dir", "d", "", "Path to the log directory (optional)")
	localCmd.Flags().StringVarP(&filterPattern, "pattern", "p", "", "Pattern to filter log lines (optional)")

	// Mark either --file or --dir as required
	localCmd.MarkFlagsMutuallyExclusive("file", "dir")
	//localCmd.MarkFlagRequired("file")
	//localCmd.MarkFlagRequired("dir")
}

// fetchLogsFromFile reads logs from a file and applies filters.
func fetchLogsFromFile(filePath, pattern string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	lines := strings.Split(string(file), "\n")
	filteredLines := filterLogs(lines, pattern)

	for _, line := range filteredLines {
		fmt.Println(line)
	}
}

// fetchLogsFromDir reads logs from all files in a directory and applies filters.
func fetchLogsFromDir(dirPath, pattern string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	var allLines []string
	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := fmt.Sprintf("%s/%s", dirPath, entry.Name())
			file, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", filePath, err)
				continue
			}
			lines := strings.Split(string(file), "\n")
			allLines = append(allLines, lines...)
		}
	}

	filteredLines := filterLogs(allLines, pattern)
	for _, line := range filteredLines {
		fmt.Println(line)
	}
}

// filterLogs applies filtering based on pattern, time range, and limit.
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

	}

	return filteredLines
}
