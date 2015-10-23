// Package xmetrics provides a generic client to log metrics from go web services.
package xmetrics // import "github.com/rs/xmetrics"

import "time"

// Client define an interface to a stats system like statsd or datadog to track
// service's metrics.
type Client interface {
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

// RequestClient is a per request metrics client with it's own tags storage
// persistent for the request lifetime.
type RequestClient interface {
	Client

	// AddTag adds a tag to the request client, this tag will be sent with all subsequent
	// stats queries.
	AddTags(tags ...string)
}

// New returns a new client with the provided backend client implementation.
func New(c Client) RequestClient {
	return &requestClient{c: c}
}

type requestClient struct {
	c Client
	// Tags are appended to the tags provided to commands
	tags []string
}

// AddTag implements RequestClient interface
func (rc *requestClient) AddTags(tags ...string) {
	rc.tags = append(rc.tags, tags...)
}

// Gauge implements RequestClient interface
func (rc *requestClient) Gauge(stat string, value float64, tags ...string) {
	tags = append(tags, rc.tags...)
	rc.c.Gauge(stat, value, tags...)
}

// Count implements RequestClient interface
func (rc *requestClient) Count(stat string, count float64, tags ...string) {
	tags = append(tags, rc.tags...)
	rc.c.Count(stat, count, tags...)
}

// Histogram implements RequestClient interface
func (rc *requestClient) Histogram(stat string, value float64, tags ...string) {
	tags = append(tags, rc.tags...)
	rc.c.Histogram(stat, value, tags...)
}

// Timing implements RequestClient interface
func (rc *requestClient) Timing(stat string, duration time.Duration, tags ...string) {
	tags = append(tags, rc.tags...)
	rc.c.Timing(stat, duration, tags...)
}
