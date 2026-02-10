package metrics

import (
	"sync"
	"time"
)

type Tracker struct {
	counters map[string]int64
	durations map[string][]time.Duration
	mu sync.RWMutex
}

func NewTracker() *Tracker {
	return &Tracker{
		counters: make(map[string]int64),
		durations: make(map[string][]time.Duration),
	}
}

func (t *Tracker) IncrementCounter(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counters[name]++
}

func (t *Tracker) ObserveDuration(name string, duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.durations[name] = append(t.durations[name], duration)
}
