package resource

import "kong/config"

type Kafka struct {
	KafkaConfig config.KafkaProducerConfig
}
