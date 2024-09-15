package main

import (
	"kong/app"
	"kong/elasticsearch"
	"kong/log"
	"kong/processor"
	"time"

	"kong/config"
	"kong/file"
	kafka "kong/kafka"
)

func main() {

	meta := config.NewAppMeta()

	switch meta.Type {
	case config.KAFKA_PRODUCER:
		app := initializeKafkaProducer(meta)
		app.Run()
	case config.KAFKA_CONSUMER:
		app := initializeKafkaConsumer(meta)
		app.Run()
	}
}

func initializeKafkaProducer(appMeta *config.AppMeta) *app.KafkaProducerApp {
	logger, err := log.NewLogger(appMeta)
	if err != nil {
		panic(err)
	}
	logger.Info("Initializing kafka producer application")

	config := config.NewKafkaProducerAppConfig(appMeta)
	fileReader := file.NewFileReader(config.FilePath, logger)
	kafkaClient, err := kafka.NewKafkaProducerClient(config.KafkaProducerConfig, logger)
	if err != nil {
		panic(err)
	}
	processor := processor.NewLogProcessor(logger)
	app := app.NewKafkaProducerApp(appMeta, config, logger, fileReader, kafkaClient, processor)

	logger.Info("Kafka producer application initialized")
	return app
}

func initializeKafkaConsumer(appMeta *config.AppMeta) *app.KafkaConsumerApp {
	logger, err := log.NewLogger(appMeta)
	if err != nil {
		panic(err)
	}
	logger.Info("Initializing kafka consumer application")

	config := config.NewKafkaConsumerAppConfig(appMeta)
	kafkaClient, err := kafka.NewKafkaConsumer(config.KafkaConsumerConfig, logger, time.Minute*10)
	if err != nil {
		panic(err)
	}
	elasticSearchClient, err := elasticsearch.NewElasticSearchClient(config.ElasticSearchConfig, logger)
	if err != nil {
		panic(err)
	}
	app := app.NewKafkaConsumerApp(appMeta, config, logger, kafkaClient, elasticSearchClient)

	logger.Info("Kafka consumer application initialized")
	return app
}
