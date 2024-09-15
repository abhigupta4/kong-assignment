package config

import (
	"os"
)

type ApplicationType string
type Env string

const (
	KAFKA_PRODUCER ApplicationType = "kafka-producer"
	KAFKA_CONSUMER ApplicationType = "kafka-consumer"
)

const (
	DEV  Env = "dev"
	PROD Env = "prod"
)

type AppMeta struct {
	Type ApplicationType `json:"type"`
	Env  Env             `json:"env"`
}

func GetAppType() ApplicationType {

	if os.Getenv("APP_TYPE") == "kafka-producer" {
		return KAFKA_PRODUCER
	} else if os.Getenv("APP_TYPE") == "kafka-consumer" {
		return KAFKA_CONSUMER
	}
	panic("APP_TYPE is not set")
}

func GetEnv() Env {
	if os.Getenv("ENV") == "prod" {
		return PROD
	}

	return DEV
}

func NewAppMeta() *AppMeta {
	return &AppMeta{
		Type: GetAppType(),
		Env:  GetEnv(),
	}
}
