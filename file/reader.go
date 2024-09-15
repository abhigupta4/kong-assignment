package file

import (
	"bufio"
	"os"

	"go.uber.org/zap"
)

type FileReader struct {
	FilePath string
	logger   *zap.Logger
}

func NewFileReader(filePath string, logger *zap.Logger) *FileReader {
	return &FileReader{
		FilePath: filePath,
		logger:   logger,
	}
}

func (fr *FileReader) ReadLines() ([]string, error) {
	file, err := os.Open(fr.FilePath)
	if err != nil {
		fr.logger.Fatal("Failed to open file: %s", zap.Any("err", err))
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fr.logger.Fatal("Error reading file: %s", zap.Any("err", err))
		return nil, err
	}

	return lines, nil
}
