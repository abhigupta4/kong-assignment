package service

import "kong/resource"

type Producer struct {
	Kafka resource.Kafka
}

func NewProducer(kafka resource.Kafka) *Producer {
	return &Producer{
		Kafka: kafka,
	}
}

