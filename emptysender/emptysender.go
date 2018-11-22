// Package emptysender implements an empty xstats Sender interface
// Useful for testing locally or perf benchmarks with and without instrumentation
package emptysender

import (
	"time"
)

type sender struct {
}

// New creates an instance of type sender
func New() *sender {
	return new(sender)
}

// Gauge implements xstats.Sender interface
func (s *sender) Gauge(stat string, value float64, tags ...string) {
}

// Count implements xstats.Sender interface
func (s *sender) Count(stat string, count float64, tags ...string) {
}

// Histogram implements xstats.Sender interface
func (s *sender) Histogram(stat string, value float64, tags ...string) {
}

// Timing implements xstats.Sender interface
func (s *sender) Timing(stat string, value time.Duration, tags ...string) {
}
