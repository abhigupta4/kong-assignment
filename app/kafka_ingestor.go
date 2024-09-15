package app

import (
	"bufio"
	"encoding/json"
	"kong/models"
	"kong/resource"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaIngestor struct {
	Kafka  resource.Kafka
	Logger zap.Logger
}

func NewKafkaIngestor(
	logger zap.Logger,
	kafka resource.Kafka,
) *KafkaIngestor {
	return &KafkaIngestor{
		Kafka: kafka,
	}
}

func (ki *KafkaIngestor) Start() {
	ki.Logger.Info("Starting the Kafka ingestor")
}

func (ki *KafkaIngestor) Stop() {
	ki.Logger.Info("Stopping the Kafka ingestor")
}

func (ki *KafkaIngestor) Run() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger.Info("Starting the application")

	file, err := os.Open("stream2.jsonl")
	if err != nil {
		log.Fatalf("Failed to open JSONL file: %s", err)
	}
	defer file.Close()

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		logger.Error("Failed to create producer: %s", zap.Any("error", err))
	}
	defer producer.Close()

	topic := "topic"

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		var record models.LogRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			log.Printf("Failed to unmarshal JSON: %s", err)
			continue
		}

		logger.Info("Record", zap.Any("record", record))

		jsonBytes, err := json.Marshal(record)
		if err != nil {
			log.Printf("Failed to marshal Change struct: %s", err)
			continue
		}

		msg := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          jsonBytes,
			Timestamp:      time.Now(),
		}

		err = producer.Produce(msg, nil)
		if err != nil {
			log.Printf("Failed to produce message: %s", err)
		}

		go func() {
			e := <-producer.Events()
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v", ev.TopicPartition.Error)
				} else {
					log.Printf("Delivered message to %v", ev.TopicPartition)
				}
			}
		}()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	producer.Flush(15 * 1000)
}
