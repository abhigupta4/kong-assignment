package config

// Ideally these config should be injected via environment variables or pulled using a config service
var bootstrapServers = "localhost:9092"
var topic = "topic"
var producerFlushInterval = 15 * 1000
var filePath = "stream2.jsonl"

type KafkaProducerAppConfig struct {
	FilePath            string              `json:"filePath"`
	KafkaProducerConfig KafkaProducerConfig `json:"kafkaIngestorConfig"`
}

type KafkaProducerConfig struct {
	BootstrapServers      string `json:"bootstrap_servers"`
	Topic                 string `json:"topic"`
	ProducerFlushInterval int    `json:"producer_flush_interval"`
}

// App meta can be used to fetch the relevant config for the application
func NewKafkaProducerAppConfig(appMeta *AppMeta) *KafkaProducerAppConfig {
	return &KafkaProducerAppConfig{
		FilePath: filePath,
		KafkaProducerConfig: KafkaProducerConfig{
			BootstrapServers:      bootstrapServers,
			Topic:                 topic,
			ProducerFlushInterval: producerFlushInterval,
		},
	}
}
