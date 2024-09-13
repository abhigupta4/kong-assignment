package main

import (
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

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
)

type Record struct {
	Before *ChangeData `json:"before"`
	After  ChangeData  `json:"after"`
	Op     string      `json:"op"`
	TsMs   int64       `json:"ts_ms"`
}

// ChangeData represents the structure for the "before" and "after" parts of the change event
type ChangeData struct {
	Key   string     `json:"key"`
	Value ObjectData `json:"value"`
}

// ObjectData represents the data stored in the "value" field
type ObjectData struct {
	Type   int        `json:"type"`
	Object ObjectInfo `json:"object"`
}

// ObjectInfo represents the details of the object being changed
type ObjectInfo struct {
	ID               string                 `json:"id"`
	Host             string                 `json:"host,omitempty"`
	Name             string                 `json:"name"`
	Path             string                 `json:"path,omitempty"`
	Port             int                    `json:"port,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	Enabled          bool                   `json:"enabled,omitempty"`
	Retries          int                    `json:"retries,omitempty"`
	Protocol         string                 `json:"protocol,omitempty"`
	CreatedAt        int64                  `json:"created_at"`
	UpdatedAt        int64                  `json:"updated_at"`
	ReadTimeout      int                    `json:"read_timeout,omitempty"`
	WriteTimeout     int                    `json:"write_timeout,omitempty"`
	ConnectTimeout   int                    `json:"connect_timeout,omitempty"`
	LastPing         int64                  `json:"last_ping,omitempty"`
	Version          string                 `json:"version,omitempty"`
	Hostname         string                 `json:"hostname,omitempty"`
	ConfigHash       string                 `json:"config_hash,omitempty"`
	ProcessConf      map[string]interface{} `json:"process_conf,omitempty"`
	ConnectionState  *ConnectionState       `json:"connection_state,omitempty"`
	DataPlaneCertID  string                 `json:"data_plane_cert_id,omitempty"`
	Labels           map[string]string      `json:"labels,omitempty"`
	ResourceType     string                 `json:"resource_type,omitempty"`
	Value            string                 `json:"value,omitempty"`
	HealthChecks     *HealthChecks          `json:"healthchecks,omitempty"`
	HashOn           string                 `json:"hash_on,omitempty"`
	HashFallback     string                 `json:"hash_fallback,omitempty"`
	UseSrvName       bool                   `json:"use_srv_name,omitempty"`
	Slots            int                    `json:"slots,omitempty"`
	Algorithm        string                 `json:"algorithm,omitempty"`
	HashOnCookiePath string                 `json:"hash_on_cookie_path,omitempty"`
}

// ConnectionState represents the state of the connection
type ConnectionState struct {
	IsConnected bool `json:"is_connected"`
}

// HealthChecks represents the structure for upstream health checks
type HealthChecks struct {
	Active    ActiveHealthCheck  `json:"active"`
	Passive   PassiveHealthCheck `json:"passive"`
	Threshold int                `json:"threshold"`
}

// ActiveHealthCheck represents active health check configuration
type ActiveHealthCheck struct {
	Type            string         `json:"type"`
	Healthy         HealthyCheck   `json:"healthy"`
	Unhealthy       UnhealthyCheck `json:"unhealthy"`
	Timeout         int            `json:"timeout"`
	HTTPPath        string         `json:"http_path"`
	Concurrency     int            `json:"concurrency"`
	HTTPSVerifyCert bool           `json:"https_verify_certificate"`
}

// PassiveHealthCheck represents passive health check configuration
type PassiveHealthCheck struct {
	Type      string         `json:"type"`
	Healthy   HealthyCheck   `json:"healthy"`
	Unhealthy UnhealthyCheck `json:"unhealthy"`
}

// HealthyCheck represents the healthy status configuration
type HealthyCheck struct {
	Interval     int   `json:"interval"`
	Successes    int   `json:"successes"`
	HTTPStatuses []int `json:"http_statuses"`
}

// UnhealthyCheck represents the unhealthy status configuration
type UnhealthyCheck struct {
	Interval     int   `json:"interval"`
	Timeouts     int   `json:"timeouts"`
	TCPFailures  int   `json:"tcp_failures"`
	HTTPFailures int   `json:"http_failures"`
	HTTPStatuses []int `json:"http_statuses"`
}

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Starting the application")

	// file, err := os.Open("stream2.jsonl")
	// if err != nil {
	// 	log.Fatalf("Failed to open JSONL file: %s", err)
	// }
	// defer file.Close()

	// producer, err := kafka.NewProducer(&kafka.ConfigMap{
	// 	"bootstrap.servers": "localhost:9092",
	// })
	// if err != nil {
	// 	logger.Error("Failed to create producer: %s", zap.Any("error", err))
	// }
	// defer producer.Close()

	// topic := "topic"

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	line := scanner.Text()

	// 	var record Record
	// 	if err := json.Unmarshal([]byte(line), &record); err != nil {
	// 		log.Printf("Failed to unmarshal JSON: %s", err)
	// 		continue
	// 	}

	// 	logger.Info("Record", zap.Any("record", record))

	// 	jsonBytes, err := json.Marshal(record)
	// 	if err != nil {
	// 		log.Printf("Failed to marshal Change struct: %s", err)
	// 		continue
	// 	}

	// 	msg := &kafka.Message{
	// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 		Value:          jsonBytes,
	// 		Timestamp:      time.Now(),
	// 	}

	// 	err = producer.Produce(msg, nil)
	// 	if err != nil {
	// 		log.Printf("Failed to produce message: %s", err)
	// 	}

	// 	go func() {
	// 		e := <-producer.Events()
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			if ev.TopicPartition.Error != nil {
	// 				log.Printf("Delivery failed: %v", ev.TopicPartition.Error)
	// 			} else {
	// 				log.Printf("Delivered message to %v", ev.TopicPartition)
	// 			}
	// 		}
	// 	}()
	// }

	// if err := scanner.Err(); err != nil {
	// 	log.Fatalf("Error reading file: %s", err)
	// }

	// producer.Flush(15 * 1000)

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

				var record Record
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

func writeToElasticsearch(osClient *opensearch.Client, record Record, logger *zap.Logger) error {
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
