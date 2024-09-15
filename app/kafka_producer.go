package app

import (
	"encoding/json"
	"kong/config"
	"kong/file"
	"kong/kafka"
	"kong/processor"

	"go.uber.org/zap"
)

type KafkaProducerApp struct {
	AppMeta      *config.AppMeta
	AppConfig    *config.KafkaProducerAppConfig
	Logger       *zap.Logger
	LogProcessor *processor.LogProcessor
	KafkaClient  *kafka.KafkaProducerClient
	FileReader   *file.FileReader
}

func NewKafkaProducerApp(
	appMeta *config.AppMeta,
	appConfig *config.KafkaProducerAppConfig,
	logger *zap.Logger,
	fileReader *file.FileReader,
	kafkaClient *kafka.KafkaProducerClient,
	processor *processor.LogProcessor,
) KafkaProducerApp {
	return KafkaProducerApp{
		AppMeta:      appMeta,
		AppConfig:    appConfig,
		Logger:       logger,
		FileReader:   fileReader,
		KafkaClient:  kafkaClient,
		LogProcessor: processor,
	}
}

func (kpa *KafkaProducerApp) Stop() {
	kpa.Logger.Info("Stopping Kafka Producer Application")
	kpa.KafkaClient.Producer.Close()
}

func (kpa *KafkaProducerApp) Run() {

	kpa.Logger.Info("Running kafka producer the application")
	lines, err := kpa.FileReader.ReadLines()
	if err != nil {
		kpa.Logger.Error("Failed to read file", zap.Any("error", err))
		return
	}

	for _, line := range lines {
		record, err := kpa.LogProcessor.ProcessMessage(line)
		if err != nil {
			kpa.Logger.Warn("Skipping invalid message", zap.String("line", line))
			continue
		}

		jsonBytes, err := json.Marshal(record)
		if err != nil {
			kpa.Logger.Warn("Failed to marshal record", zap.Any("record", record))
			continue
		}

		err = kpa.KafkaClient.SendMessage(jsonBytes)
		if err != nil {
			kpa.Logger.Error("Failed to send message to Kafka", zap.Any("error", err))
		}
		kpa.Logger.Debug("Sent message to Kafka", zap.Any("record", record))
	}
}
