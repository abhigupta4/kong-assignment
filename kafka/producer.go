package kafka

import (
	"kong/config"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaProducerClient struct {
	Logger        *zap.Logger
	Producer      *kafka.Producer
	Topic         string
	FlushInterval int
}

func NewKafkaProducerClient(config config.KafkaProducerConfig, logger *zap.Logger) (*KafkaProducerClient, error) {
	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": config.BootstrapServers,
		},
	)
	if err != nil {
		logger.Error("Failed to create Kafka producer: %v", zap.Any("error", err))
		return nil, err
	}

	return &KafkaProducerClient{
		Logger:        logger,
		Producer:      producer,
		Topic:         config.Topic,
		FlushInterval: config.ProducerFlushInterval,
	}, nil
}

func (kc *KafkaProducerClient) SendMessage(value []byte, partition int32) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kc.Topic, Partition: partition},
		Value:          value,
		Timestamp:      time.Now(),
	}

	err := kc.Producer.Produce(msg, nil)
	if err != nil {
		kc.Logger.Error("Failed to produce message: %s", zap.Any("error", err))
		return err
	}

	go func() {
		e := <-kc.Producer.Events()
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				kc.Logger.Error("Delivery failed: %v", zap.Any("error", ev.TopicPartition.Error))
			} else {
				kc.Logger.Debug("Delivered message to %v", zap.Any("topicPartition", ev.TopicPartition))
			}
		}
	}()

	return nil
}

func (kc *KafkaProducerClient) Flush() {
	kc.Logger.Info("Flushing remaining messages")
	kc.Producer.Flush(kc.FlushInterval)
}

func (kc *KafkaProducerClient) Shutdown() {
	kc.Producer.Flush(kc.FlushInterval)
	kc.Producer.Close()
}
