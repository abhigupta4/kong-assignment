package main

import (
	"kong/app"
	"kong/log"
	"kong/processor"

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
	}
}

func initializeKafkaProducer(appMeta *config.AppMeta) app.KafkaProducerApp {
	logger, err := log.NewLogger(appMeta)
	if err != nil {
		panic(err)
	}
	logger.Info("Initializing kafka producer application")

	config := config.NewKafkaProducerAppConfig(appMeta)
	fileReader := file.NewFileReader(config.FilePath, logger)
	kafkaClient, err := kafka.NewKafkaClient(config.KafkaProducerConfig, logger)
	if err != nil {
		panic(err)
	}
	processor := processor.NewLogProcessor(logger)
	app := app.NewKafkaProducerApp(appMeta, config, logger, fileReader, kafkaClient, processor)

	logger.Info("Kafka producer application initialized")
	return app
}
