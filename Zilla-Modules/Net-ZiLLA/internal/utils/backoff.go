package utils

import (
	"math"
	"time"
)

// Backoff implements exponential backoff
type Backoff struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Factor       float64
	Attempts     int
}

func NewBackoff(initial time.Duration, max time.Duration) *Backoff {
	return &Backoff{
		InitialDelay: initial,
		MaxDelay:     max,
		Factor:       2.0,
	}
}

func (b *Backoff) Next() time.Duration {
	delay := float64(b.InitialDelay) * math.Pow(b.Factor, float64(b.Attempts))
	b.Attempts++
	
	actualDelay := time.Duration(delay)
	if actualDelay > b.MaxDelay {
		return b.MaxDelay
	}
	return actualDelay
}

func (b *Backoff) Reset() {
	b.Attempts = 0
}
