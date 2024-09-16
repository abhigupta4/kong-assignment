package limiter

import (
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
	stopCh chan struct{}
}

func NewRateLimiter(rate int) *RateLimiter {
	limiter := &RateLimiter{
		tokens: make(chan struct{}, rate),
		stopCh: make(chan struct{}),
	}

	go limiter.refill(rate)

	return limiter
}

func (rl *RateLimiter) refill(rate int) {
	interval := time.Second / time.Duration(rate)
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		case <-rl.stopCh:
			ticker.Stop()
			return
		}
	}
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}
