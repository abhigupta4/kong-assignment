package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kong/config"
	"kong/models"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
)

type KafkaConsumerApp struct {
	AppMeta   *config.AppMeta
	AppConfig *config.KafkaProducerAppConfig
	Logger    *zap.Logger
}

func NewKafkaConsumerApp() *KafkaConsumerApp {
	return &KafkaConsumerApp{}
}

func (ei *KafkaConsumerApp) Ingest() {
	// Ingest data into ElasticSearch
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
	// if err != nil {
	// 	log.Fatalf("Error creating the OpenSearch client: %s", err)
	// }

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
