package config

// Ideally these config should be injected via environment variables or pulled using a config service
var bootstrapServers = []string{"localhost:9092"}
var topic = "topic"
var producerFlushInterval = 15 * 1000
var filePath = "stream2.jsonl"

type KafkaIngestorAppConfig struct {
	AppMeta             Meta                `json:"appMeta"`
	FilePath            string              `json:"filePath"`
	KafkaProducerConfig KafkaProducerConfig `json:"kafkaIngestorConfig"`
}

type KafkaProducerConfig struct {
	BootstrapServers      []string `json:"bootstrap_servers"`
	Topic                 string   `json:"topic"`
	ProducerFlushInterval int      `json:"producer_flush_interval"`
}

func NewKafkaIngestorAppConfig() *KafkaIngestorAppConfig {
	return &KafkaIngestorAppConfig{
		AppMeta: Meta{
			AppType: "kafka_ingestor",
		},
		FilePath: filePath,
		KafkaProducerConfig: KafkaProducerConfig{
			BootstrapServers:      bootstrapServers,
			Topic:                 topic,
			ProducerFlushInterval: producerFlushInterval,
		},
	}
}
