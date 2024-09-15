package kafka

import (
	"kong/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaConsumerClient struct {
	Logger        *zap.Logger
	Consumer      *kafka.Consumer
	Topic         string
	FlushInterval int
}

func NewKafkaConsumerClient(config config.KafkaProducerConfig, logger *zap.Logger) (*KafkaConsumerClient, error) {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "consumer-group3",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Error("Failed to create Kafka consumer: %s", zap.Any("error", err))
	}

	return &KafkaConsumerClient{
		Logger:        logger,
		Consumer:      kafkaConsumer,
		Topic:         config.Topic,
		FlushInterval: config.ProducerFlushInterval,
	}, nil
}

func (kc *KafkaConsumerClient) Close() {
	kc.Consumer.Close()
}
