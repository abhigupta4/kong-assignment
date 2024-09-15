package log

import (
	"kong/config"

	"go.uber.org/zap"
)

func NewLogger(appMeta *config.AppMeta) (*zap.Logger, error) {

	var logger *zap.Logger
	var err error
	if appMeta.Env == config.PROD {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}
	return logger, nil
}
