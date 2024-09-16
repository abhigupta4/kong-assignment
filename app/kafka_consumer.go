package app

import (
	"encoding/json"
	"kong/config"
	"kong/elasticsearch"
	"kong/kafka"
	"kong/limiter"
	"kong/models"
	"kong/retry"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaConsumerApp struct {
	AppMeta             *config.AppMeta
	AppConfig           *config.KafkaConsumerAppConfig
	Logger              *zap.Logger
	KafkaClient         *kafka.KafkaConsumer
	ElasticSearchClient *elasticsearch.ElasticSearchClient
	RateLimiter         *limiter.RateLimiter
	StopChan            chan struct{}
}

func NewKafkaConsumerApp(
	appMeta *config.AppMeta,
	appConfig *config.KafkaConsumerAppConfig,
	logger *zap.Logger,
	kafkaClient *kafka.KafkaConsumer,
	elasticSearchClient *elasticsearch.ElasticSearchClient,
	rateLimiter *limiter.RateLimiter,
) *KafkaConsumerApp {
	return &KafkaConsumerApp{
		AppMeta:             appMeta,
		AppConfig:           appConfig,
		Logger:              logger,
		KafkaClient:         kafkaClient,
		ElasticSearchClient: elasticSearchClient,
		RateLimiter:         rateLimiter,
		StopChan:            make(chan struct{}),
	}
}

func (kca *KafkaConsumerApp) Run() error {
	kca.Logger.Info("Running kafka consumer application")
	for {
		select {
		case <-kca.StopChan:
			kca.Logger.Info("Received stop signal, exiting message processing")
			return nil
		default:
			msg, err := kca.KafkaClient.ReadMessage()
			if err == nil {
				kca.Logger.Info("Received message from Kafka", zap.String("message", string(msg.Value)))

				var record models.Record
				err := json.Unmarshal(msg.Value, &record)
				if err != nil {
					kca.Logger.Error("Failed to unmarshal message", zap.Error(err))
				}

				if kca.RateLimiter.Allow() {
					err := retry.RetryWithExponentialBackoff(
						kca.AppConfig.RetryConfig.Count,
						kca.AppConfig.RetryConfig.Backoff,
						func() error {
							return kca.ElasticSearchClient.Write(record, kca.AppConfig.ElasticSearchConfig.Index, record.After.Key)
						})
					if err != nil {
						kca.Logger.Error("Failed to write to Elasticsearch", zap.Error(err))
					}
					kca.Logger.Debug("Successfully written to ElasticSearch")
				} else {
					kca.Logger.Warn("Rate limit exceeded, skipping writing to ElasticSearch")
				}
			} else if err.(confluentKafka.Error).Code() != confluentKafka.ErrTimedOut {
				kca.Logger.Error("Consumer error", zap.Error(err))
			} else {
				kca.Logger.Debug("Timed out waiting for message")
			}

		}
	}
}

func (kca *KafkaConsumerApp) Shutdown() error {
	kca.Logger.Info("Stopping Kafka Consumer Application")
	close(kca.StopChan)
	kca.KafkaClient.Shutdown()
	kca.RateLimiter.Stop()

	return nil
}
