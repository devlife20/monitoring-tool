package linux

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LogFile represents a log file with its name and path.
type LogFile struct {
	Name string
	Path string
}

// FetchLogFiles retrieves a list of log files from the specified directory.
func FetchLogFiles(dir string) ([]LogFile, error) {
	var logFiles []LogFile
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Include only regular files with .log extension or any specific criteria
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".log")) {
			logFiles = append(logFiles, LogFile{Name: info.Name(), Path: path})
		}
		return nil
	})
	return logFiles, err
}

// ReadLogFile reads the contents of a log file and prints it line by line.

func ReadFileContent(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("Failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	var builder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		builder.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		builder.WriteString(fmt.Sprintf("\nError reading file: %v", err))
	}

	return builder.String()
}

//func ScanAndReadLogFile(logDir string) {
//	fmt.Printf("Scanning directory: %s\n", logDir)
//
//	// Fetch log files
//	logFiles, err := FetchLogFiles(logDir)
//	if err != nil {
//		fmt.Printf("Error fetching log files: %v\n", err)
//		return
//	}
//
//	if len(logFiles) == 0 {
//		fmt.Println("No log files found.")
//		return
//	}
//
//	// Display available log files
//	fmt.Println("Available log files:")
//	for i, logFile := range logFiles {
//		fmt.Printf("[%d] %s\n", i+1, logFile.Name)
//	}
//
//	// Prompt user to select a log file to read
//	var choice int
//	fmt.Printf("\nEnter the number of the log file to read (1-%d): ", len(logFiles))
//	_, err = fmt.Scan(&choice)
//	if err != nil || choice < 1 || choice > len(logFiles) {
//		fmt.Println("Invalid choice. Exiting.")
//		return
//	}
//
//	selectedFile := logFiles[choice-1]
//	fmt.Printf("You selected: %s\n", selectedFile.Name)
//
//	// Read the selected log file
//	//err = ReadFileContent(selectedFile.Path)
//	if err != nil {
//		fmt.Printf("Error reading log file: %v\n", err)
//	}
//}
