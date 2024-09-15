package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kong/config"
	"kong/models"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
)

func main() {

	meta := config.NewMeta()

	switch meta.AppType {
	case config.KAFKA_PRODUCER:
		fileStreamToKafkaApplication()
	case config.KAFKA_CONSUMER:
		kafkaToElasticSearchApplication()
	}
}

func kafkaToElasticSearchApplication() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Starting the application")

	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "consumer-group3",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Error("Failed to create Kafka consumer: %s", zap.Any("error", err))
	}
	defer kafkaConsumer.Close()

	topic := "topic"
	err = kafkaConsumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logger.Error("Failed to subscribe to topics: %s", zap.Any("error", err))
	}

	osClient, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{
			"http://localhost:9200", // OpenSearch address
		},
	})
	if err != nil {
		log.Fatalf("Error creating the OpenSearch client: %s", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Starting Kafka consumer...")
	run := true
	for run {
		select {
		case sig := <-sigChan:
			logger.Info("Caught signal %v: terminating\n", zap.Any("signal", sig))
			run = false
		default:
			msg, err := kafkaConsumer.ReadMessage(100 * time.Millisecond)
			if err == nil {
				logger.Info("Received message from Kafka", zap.String("message", string(msg.Value)))

				var record models.LogRecord
				err := json.Unmarshal(msg.Value, &record)
				if err != nil {
					logger.Error("Failed to unmarshal message", zap.Any("error", err))
					continue
				}

				err = writeToElasticsearch(osClient, record, logger)
				if err != nil {
					logger.Error("Failed to write to Elasticsearch", zap.Any("error", err))
				}

			} else if err.(kafka.Error).Code() != kafka.ErrTimedOut {
				logger.Error("Consumer error", zap.Any("error", err))
			}
		}
	}
}

func fileStreamToKafkaApplication() {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

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

func writeToElasticsearch(osClient *opensearch.Client, record models.LogRecord, logger *zap.Logger) error {
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal change to JSON: %w", err)
	}

	docID := url.PathEscape(record.After.Key)

	req := opensearchapi.IndexRequest{
		Index:      "kafka-index",
		DocumentID: docID,
		Body:       bytes.NewReader(jsonBytes),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), osClient)
	if err != nil {
		return fmt.Errorf("failed to execute OpenSearch request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s: %s", docID, res.String())
	}

	fmt.Printf("Successfully wrote document ID=%s to OpenSearch\n", docID)
	return nil
}
