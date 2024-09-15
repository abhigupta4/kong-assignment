package config

const CONSUMER_BOOTSTRAP_SERVERS = "localhost:9092"
const CONSUMER_TOPIC = "topic"
const CONSUMER_GROUP = "consumer-group-2"
const AUTO_OFFSET_RESET = "earliest"
const ELASTIC_SEARCH_ADDRESS = "http://localhost:9200"
const CONSUMER_MESSAGE_READ_TIMEOUT = 1000

type KafkaConsumerAppConfig struct {
	KafkaConsumerConfig KafkaConsumerConfig `json:"kafkaConsumerConfig"`
	ElasticSearchConfig ElasticSearchConfig `json:"elasticSearchConfig"`
}

type KafkaConsumerConfig struct {
	BootstrapServers           string `json:"bootstrap_servers"`
	Topic                      string `json:"topic"`
	ConsumerGroup              string `json:"consumer_group"`
	AutoOffsetReset            string `json:"auto_offset_reset"`
	ConsumerMessageReadTimeout int    `json:"consumer_message_read_timeout"`
}

type ElasticSearchConfig struct {
	Address string `json:"address"`
}

// App meta can be used to fetch the relevant config for the application
func NewKafkaConsumerAppConfig(appMeta *AppMeta) *KafkaConsumerAppConfig {
	return &KafkaConsumerAppConfig{
		KafkaConsumerConfig: KafkaConsumerConfig{
			BootstrapServers:           CONSUMER_BOOTSTRAP_SERVERS,
			Topic:                      CONSUMER_TOPIC,
			ConsumerGroup:              CONSUMER_GROUP,
			AutoOffsetReset:            AUTO_OFFSET_RESET,
			ConsumerMessageReadTimeout: CONSUMER_MESSAGE_READ_TIMEOUT,
		},
		ElasticSearchConfig: ElasticSearchConfig{
			Address: ELASTIC_SEARCH_ADDRESS,
		},
	}
}
