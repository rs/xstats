// Package mock implements mock object for xstats Sender interface based on
// github.com/stretchr/testify/mock package mock implementation
package mock

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type sender struct {
	mock.Mock
}

// New creates an instance of type sender
func New() *sender {
	return new(sender)
}

// Gauge implements xstats.Sender interface
func (s *sender) Gauge(stat string, value float64, tags ...string) {
	s.Called(stat, value, tags)
}

// Count implements xstats.Sender interface
func (s *sender) Count(stat string, count float64, tags ...string) {
	s.Called(stat, count, tags)
}

// Histogram implements xstats.Sender interface
func (s *sender) Histogram(stat string, value float64, tags ...string) {
	s.Called(stat, value, tags)
}

// Timing implements xstats.Sender interface
func (s *sender) Timing(stat string, value time.Duration, tags ...string) {
	s.Called(stat, value, tags)
}
