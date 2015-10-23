package xstats

import "time"

type nopC struct {
}

var nopClient = &nopC{}

// AddTag implements RequestClient interface
func (rc *nopC) AddTags(tags ...string) {
}

// Gauge implements RequestClient interface
func (rc *nopC) Gauge(stat string, value float64, tags ...string) {
}

// Count implements RequestClient interface
func (rc *nopC) Count(stat string, count float64, tags ...string) {
}

// Histogram implements RequestClient interface
func (rc *nopC) Histogram(stat string, value float64, tags ...string) {
}

// Timing implements RequestClient interface
func (rc *nopC) Timing(stat string, duration time.Duration, tags ...string) {
}
