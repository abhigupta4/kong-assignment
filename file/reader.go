package file

import (
	"bufio"
	"log"
	"os"
)

type FileReader struct {
	FilePath string
}

func NewFileReader(filePath string) *FileReader {
	return &FileReader{FilePath: filePath}
}

func (fr *FileReader) ReadLines() ([]string, error) {
	file, err := os.Open(fr.FilePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
		return nil, err
	}

	return lines, nil
}
