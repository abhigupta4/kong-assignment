package factory

import (
	"kong/app"
	"kong/config"
	"kong/file"
	"kong/kafka"
	"kong/limiter"
	"kong/processor"

	"go.uber.org/zap"
)

func InitializeKafkaProducer(appMeta *config.AppMeta, logger *zap.Logger) (*app.KafkaProducerApp, error) {
	logger.Info("Initializing kafka producer application")

	config := config.NewKafkaProducerAppConfig(appMeta)
	fileReader := file.NewFileReader(config.FilePath, logger)
	kafkaClient, err := kafka.NewKafkaProducerClient(config.KafkaProducerConfig, logger)
	if err != nil {
		return nil, err
	}
	processor := processor.NewLogProcessor(logger)
	rateLimiter := limiter.NewRateLimiter(config.RateLimiterConfig.Rate)
	app := app.NewKafkaProducerApp(appMeta, config, logger, fileReader, kafkaClient, processor, rateLimiter)

	logger.Info("Kafka producer application initialized")
	return app, nil
}
