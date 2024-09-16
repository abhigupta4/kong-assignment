package config

import "time"

// Ideally these config should be injected via environment variables or pulled using a config service
const TOPIC = "topic"

const FILE_PATH = "stream.jsonl"

const PRODUCER_BOOTSTRAP_SERVERS = "localhost:9092"
const PRODUCER_FLUSH_INTERVAL = 15 * 1000

const AUTO_OFFSET_RESET = "earliest"
const CONSUMER_MESSAGE_READ_TIMEOUT = 1000
const CONSUMER_GROUP = "consumer-group-1"
const CONSUMER_BOOTSTRAP_SERVERS = "localhost:9092"

const ELASTIC_SEARCH_ADDRESS = "http://localhost:9200"
const ES_INDEX = "cdc"

const RATE = 1000000

const RETRY_COUNT = 3
const RETRY_BACKOFF = time.Second * 5
