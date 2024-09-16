package config

import "time"

type RetryConfig struct {
	Count   int           `json:"count"`
	Backoff time.Duration `json:"backoff"`
}
