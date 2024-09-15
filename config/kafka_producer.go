package config

// Ideally these config should be injected via environment variables or pulled using a config service
const PRODUCER_BOOTSTRAP_SERVERS = "localhost:9092"
const PRODUCER_TOPIC = "topic"
const PRODUCER_FLUSH_INTERVAL = 15 * 1000
const FILE_PATH = "stream2.jsonl"

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
		FilePath: FILE_PATH,
		KafkaProducerConfig: KafkaProducerConfig{
			BootstrapServers:      PRODUCER_BOOTSTRAP_SERVERS,
			Topic:                 PRODUCER_TOPIC,
			ProducerFlushInterval: PRODUCER_FLUSH_INTERVAL,
		},
	}
}
