package xstats

import "time"

type nopS struct {
}

var nop = &nopS{}

// AddTags implements XStats interface
func (rc *nopS) AddTags(tags ...string) {
}

// GetTags implements XStats interface
func (rc *nopS) GetTags() []string {
	return nil
}

// Gauge implements XStats interface
func (rc *nopS) Gauge(stat string, value float64, tags ...string) {
}

// Count implements XStats interface
func (rc *nopS) Count(stat string, count float64, tags ...string) {
}

// Histogram implements XStats interface
func (rc *nopS) Histogram(stat string, value float64, tags ...string) {
}

// Timing implements XStats interface
func (rc *nopS) Timing(stat string, duration time.Duration, tags ...string) {
}
