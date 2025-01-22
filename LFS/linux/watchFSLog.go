package linux

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func TailLogsFromFile(filePath, pattern string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				file.Seek(0, io.SeekCurrent)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return err
		}

		if pattern != "" && !strings.Contains(line, pattern) {
			continue
		}

		fmt.Print(line)
	}
}
