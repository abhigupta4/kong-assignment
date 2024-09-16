package app

import (
	"encoding/json"
	"kong/config"
	"kong/file"
	"kong/kafka"
	"kong/processor"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaProducerApp struct {
	AppMeta      *config.AppMeta
	AppConfig    *config.KafkaProducerAppConfig
	Logger       *zap.Logger
	LogProcessor *processor.LogProcessor
	KafkaClient  *kafka.KafkaProducerClient
	FileReader   *file.FileReader
	StopChan     chan struct{}
}

func NewKafkaProducerApp(
	appMeta *config.AppMeta,
	appConfig *config.KafkaProducerAppConfig,
	logger *zap.Logger,
	fileReader *file.FileReader,
	kafkaClient *kafka.KafkaProducerClient,
	processor *processor.LogProcessor,
) *KafkaProducerApp {
	return &KafkaProducerApp{
		AppMeta:      appMeta,
		AppConfig:    appConfig,
		Logger:       logger,
		FileReader:   fileReader,
		KafkaClient:  kafkaClient,
		LogProcessor: processor,
		StopChan:     make(chan struct{}),
	}
}

// Ordering matters hence not sending each message as a go routine
func (kpa *KafkaProducerApp) Run() error {

	kpa.Logger.Info("Running kafka producer the application")
	lines, err := kpa.FileReader.ReadLines()
	if err != nil {
		kpa.Logger.Error("Failed to read file", zap.Any("error", err))
		return err
	}

	for _, line := range lines {
		select {
		case <-kpa.StopChan:
			kpa.Logger.Info("Received stop signal, exiting message processing")
			return nil
		default:
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

			err = kpa.KafkaClient.SendMessage(jsonBytes, kpa.getKafkaPartition(record.After.Key))
			if err != nil {
				kpa.Logger.Error("Failed to send message to Kafka", zap.Any("error", err))
			}
			kpa.Logger.Debug("Sent message to Kafka", zap.Any("record", record))
		}
	}
	return nil
}

// Can change this to get the right partition based on the key
func (kpa *KafkaProducerApp) getKafkaPartition(key string) int32 {
	return confluentKafka.PartitionAny
}

func (kpa *KafkaProducerApp) Shutdown() error {
	close(kpa.StopChan)
	kpa.KafkaClient.Shutdown()

	return nil
}
