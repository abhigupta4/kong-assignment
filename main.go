package main

import (
	"kong/app"
	"kong/factory"
	"kong/log"
	"os"
	"os/signal"
	"syscall"

	"kong/config"

	"go.uber.org/zap"
)

func main() {

	meta := config.NewAppMeta()
	logger, err := log.NewLogger(meta)
	if err != nil {
		panic(err)
	}

	var app app.App

	app, err = initializeApplication(meta, logger)
	if err != nil {
		panic(err)
	}

	if app == nil {
		panic("Application not initialized")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger.Info("Running application...")
	go func() {
		if err := app.Run(); err != nil {
			logger.Error("Application failed", zap.Any("err", err))
			logger.Info("Shutting down application...")
			app.Shutdown()
			return
		}
	}()

	receivedSignal := <-signalChan
	logger.Info("Received shutdown signal, shutting down gracefully...", zap.Any("signal", receivedSignal))

	if err := app.Shutdown(); err != nil {
		logger.Error("Failed to shutdown application gracefully", zap.Any("err", err))
		return
	}

	logger.Info("Application shutdown successfully")
}

func initializeApplication(appMeta *config.AppMeta, logger *zap.Logger) (app.App, error) {
	switch appMeta.Type {
	case config.KAFKA_PRODUCER:
		return factory.InitializeKafkaProducer(appMeta, logger)
	case config.KAFKA_CONSUMER:
		return factory.InitializeKafkaConsumer(appMeta, logger)
	}
	return nil, nil
}
