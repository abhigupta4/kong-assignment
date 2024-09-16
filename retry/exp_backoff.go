package retry

import "time"

func RetryWithExponentialBackoff(retries int, backoff time.Duration, fn func() error) error {
	var err error
	for i := 0; i < retries; i++ {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return err
}
