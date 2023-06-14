package fileutil

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

// ReadLines reads the lines from the file at the given path and returns them as a slice of strings.
func ReadLines(path string) ([]string, error) {
	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new scanner for the file.
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // Use a 1MB buffer

	// Read the lines from the file and store them in a bytes.Buffer. Also discard the first word, and only take the string withing the single quotes
	var buf bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		buf.WriteString(strings.Trim(fields[1], "'"))
		buf.WriteByte('\n')
	}

	// Check for any errors that occurred while scanning the file.
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Convert the bytes.Buffer to a slice of strings.
	lines := strings.Split(buf.String(), "\n")

	// Return the slice of lines.
	return lines, nil
}
