package kafka

import (
	"kong/config"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	Logger   *zap.Logger
	Consumer *kafka.Consumer
	Topic    string
	Timeout  time.Duration
}

func NewKafkaConsumer(config config.KafkaConsumerConfig, logger *zap.Logger, timeout time.Duration) (*KafkaConsumer, error) {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          config.ConsumerGroup,
		"auto.offset.reset": config.AutoOffsetReset,
	})
	if err != nil {
		logger.Error("Failed to create Kafka consumer: %s", zap.Any("error", err))
	}

	err = kafkaConsumer.Subscribe(config.Topic, nil)
	if err != nil {
		logger.Error("Failed to subscribe to topic", zap.Any("error", err))
		return nil, err
	}

	return &KafkaConsumer{
		Logger:   logger,
		Consumer: kafkaConsumer,
		Topic:    config.Topic,
		Timeout:  timeout,
	}, nil
}

func (kc *KafkaConsumer) Close() {
	kc.Consumer.Close()
}

func (kc *KafkaConsumer) ReadMessage() (*kafka.Message, error) {
	return kc.Consumer.ReadMessage(kc.Timeout)
}
