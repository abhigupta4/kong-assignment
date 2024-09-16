package factory

import (
	"kong/app"
	"kong/config"
	"kong/elasticsearch"
	"kong/kafka"
	"time"

	"go.uber.org/zap"
)

func InitializeKafkaConsumer(appMeta *config.AppMeta, logger *zap.Logger) (*app.KafkaConsumerApp, error) {
	logger.Info("Initializing kafka consumer application")

	config := config.NewKafkaConsumerAppConfig(appMeta)
	kafkaClient, err := kafka.NewKafkaConsumer(config.KafkaConsumerConfig, logger, time.Minute*10)
	if err != nil {
		return nil, err
	}
	elasticSearchClient, err := elasticsearch.NewElasticSearchClient(config.ElasticSearchConfig, logger)
	if err != nil {
		return nil, err
	}
	app := app.NewKafkaConsumerApp(appMeta, config, logger, kafkaClient, elasticSearchClient)

	logger.Info("Kafka consumer application initialized")
	return app, nil
}
