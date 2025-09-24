package service

import (
	"math"
	"math/rand"
	"time"
)

type BackoffWithJitter struct {
	jitter     float64 // 0.0 - 1.0
	maxRetries int
	attempt    int
}

func NewBackoffWithJitter(maxRetries int) *BackoffWithJitter {
	return &BackoffWithJitter{
		jitter:     0.2,
		maxRetries: maxRetries,
		attempt:    0,
	}
}

func (b *BackoffWithJitter) Next() (time.Duration, bool) {
	if b.attempt >= b.maxRetries {
		return 0, false
	}

	delay := float64(100*time.Millisecond) * math.Pow(2, float64(b.attempt))

	jitter := (rand.Float64()*2 - 1) * b.jitter * delay
	delayWithJitter := delay + jitter

	if delayWithJitter > float64(5*time.Second) {
		delayWithJitter = float64(5 * time.Second)
	}

	b.attempt++
	return time.Duration(delayWithJitter), true
}

func (b *BackoffWithJitter) Reset() {
	b.attempt = 0
}

func (b *BackoffWithJitter) GetAttempt() int {
	return b.attempt
}
