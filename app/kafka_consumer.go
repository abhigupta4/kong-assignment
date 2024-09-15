package app

import (
	"encoding/json"
	"kong/config"
	"kong/elasticsearch"
	"kong/kafka"
	"kong/models"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaConsumerApp struct {
	AppMeta             *config.AppMeta
	AppConfig           *config.KafkaConsumerAppConfig
	Logger              *zap.Logger
	KafkaClient         *kafka.KafkaConsumer
	ElasticSearchClient *elasticsearch.ElasticSearchClient
}

func NewKafkaConsumerApp(
	appMeta *config.AppMeta,
	appConfig *config.KafkaConsumerAppConfig,
	logger *zap.Logger,
	kafkaClient *kafka.KafkaConsumer,
	elasticSearchClient *elasticsearch.ElasticSearchClient,
) *KafkaConsumerApp {
	return &KafkaConsumerApp{
		AppMeta:             appMeta,
		AppConfig:           appConfig,
		Logger:              logger,
		KafkaClient:         kafkaClient,
		ElasticSearchClient: elasticSearchClient,
	}
}

func (kca *KafkaConsumerApp) Run() {
	kca.Logger.Info("Running kafka consumer application")
	for {
		msg, err := kca.KafkaClient.ReadMessage()
		if err == nil {
			kca.Logger.Info("Received message from Kafka", zap.String("message", string(msg.Value)))

			var record models.LogRecord
			err := json.Unmarshal(msg.Value, &record)
			if err != nil {
				kca.Logger.Error("Failed to unmarshal message", zap.Error(err))
			}

			err = kca.ElasticSearchClient.Write(record, "kafka-index-1", record.After.Key)
			if err != nil {
				kca.Logger.Error("Failed to write to Elasticsearch", zap.Error(err))
			}
			kca.Logger.Debug("Successfully written to ElasticSearch")
		} else if err.(confluentKafka.Error).Code() != confluentKafka.ErrTimedOut {
			kca.Logger.Error("Consumer error", zap.Error(err))
		} else {
			kca.Logger.Debug("Timed out waiting for message")
		}
	}
}
