package processor

import (
	"encoding/json"
	"kong/models"

	"go.uber.org/zap"
)

type LogProcessor struct {
	Logger *zap.Logger
}

func NewLogProcessor(logger *zap.Logger) *LogProcessor {
	return &LogProcessor{
		Logger: logger,
	}
}

func (lp *LogProcessor) ProcessMessage(line string) (*models.Record, error) {
	var record models.Record
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		lp.Logger.Error("Failed to unmarshal JSON: %s", zap.Any("error", err))
		return nil, err
	}
	return &record, nil
}
