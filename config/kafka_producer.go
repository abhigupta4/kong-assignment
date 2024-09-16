package config

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
			Topic:                 TOPIC,
			ProducerFlushInterval: PRODUCER_FLUSH_INTERVAL,
		},
	}
}
