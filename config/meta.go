package config

import (
	"os"
)

type ApplicationType string

const (
	KAFKA_PRODUCER ApplicationType = "kafka-producer"
	KAFKA_CONSUMER ApplicationType = "kafka-consumer"
)

type Meta struct {
	AppType ApplicationType `json:"type"`
}

func GetAppType() ApplicationType {

	if os.Getenv("APP_TYPE") == "kafka-producer" {
		return KAFKA_PRODUCER
	} else if os.Getenv("APP_TYPE") == "kafka-consumer" {
		return KAFKA_CONSUMER
	}
	panic("APP_TYPE is not set")
}

func NewMeta() *Meta {
	return &Meta{
		AppType: GetAppType(),
	}
}
