package xstats

import "time"

// Sender define an interface to a stats system like statsd or datadog to send
// service's metrics.
type Sender interface {
	// Gauge measure the value of a particular thing at a particular time,
	// like the amount of fuel in a car’s gas tank or the number of users
	// connected to a system.
	Gauge(stat string, value float64, tags ...string)

	// Count track how many times something happened per second, like
	// the number of database requests or page views.
	Count(stat string, count float64, tags ...string)

	// Histogram track the statistical distribution of a set of values,
	// like the duration of a number of database queries or the size of
	// files uploaded by users. Each histogram will track the average,
	// the minimum, the maximum, the median, the 95th percentile and the count.
	Histogram(stat string, value float64, tags ...string)

	// Timing mesures the elapsed time
	Timing(stat string, value time.Duration, tags ...string)
}

// MultiSender lets you assign more than one sender to xstats in order to
// multicast observeration to different systems.
type MultiSender []Sender

// Gauge implements xstats.Sender interface
func (s MultiSender) Gauge(stat string, value float64, tags ...string) {
	for _, ss := range s {
		ss.Gauge(stat, value, tags...)
	}
}

// Count implements xstats.Sender interface
func (s MultiSender) Count(stat string, count float64, tags ...string) {
	for _, ss := range s {
		ss.Count(stat, count, tags...)
	}
}

// Histogram implements xstats.Sender interface
func (s MultiSender) Histogram(stat string, value float64, tags ...string) {
	for _, ss := range s {
		ss.Histogram(stat, value, tags...)
	}
}

// Timing implements xstats.Sender interface
func (s MultiSender) Timing(stat string, duration time.Duration, tags ...string) {
	for _, ss := range s {
		ss.Timing(stat, duration, tags...)
	}
}
