package config

const CONSUMER_BOOTSTRAP_SERVERS = "localhost:9092"
const CONSUMER_TOPIC = "topic"
const CONSUMER_GROUP = "consumer-group-1"
const AUTO_OFFSET_RESET = "earliest"

type KafkaConsumerAppConfig struct {
}

type KafkaConsumerConfig struct {
	BootstrapServers string `json:"bootstrap_servers"`
	Topic            string `json:"topic"`
}
