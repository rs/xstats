// Package xstats is a generic client for service instrumentation.
//
// xstats is inspired from Go-kit's metrics (https://github.com/go-kit/kit/tree/master/metrics)
// package but it takes a slightly different path. Instead of having to create
// an instance for each metric, xstats use a single instance to log every metrics
// you want. This reduces the boiler plate when you have a lot a metrics in your app.
// It's also easier in term of dependency injection.
//
// Talking about dependency injection, xstats comes with a xhandler.Handler
// integration so it can automatically inject the xstats client within the net/context
// of each request. Each request's xstats instance have its own tags storage ;
// This let you inject some per request contextual tags to be included with all
// observations sent within the lifespan of the request.
//
// xstats is pluggable and comes with integration for StatsD and DogStatsD,
// the Datadog (http://datadoghq.com) augmented version of StatsD with support for tags.
// More integration may come later (PR welcome).
package xstats // import "github.com/rs/xstats"

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

// XStater is a wrapper around a Sender to inject env tags within all observations.
type XStater interface {
	Sender

	// AddTag adds a tag to the request client, this tag will be sent with all subsequent
	// stats queries.
	AddTags(tags ...string)
}

// New returns a new client with the provided backend client implementation.
func New(s Sender) XStater {
	return &xstats{s: s}
}

type xstats struct {
	s Sender
	// Tags are appended to the tags provided to commands
	tags []string
}

// AddTag implements XStats interface
func (xs *xstats) AddTags(tags ...string) {
	xs.tags = append(xs.tags, tags...)
}

// Gauge implements XStats interface
func (xs *xstats) Gauge(stat string, value float64, tags ...string) {
	tags = append(tags, xs.tags...)
	xs.s.Gauge(stat, value, tags...)
}

// Count implements XStats interface
func (xs *xstats) Count(stat string, count float64, tags ...string) {
	tags = append(tags, xs.tags...)
	xs.s.Count(stat, count, tags...)
}

// Histogram implements XStats interface
func (xs *xstats) Histogram(stat string, value float64, tags ...string) {
	tags = append(tags, xs.tags...)
	xs.s.Histogram(stat, value, tags...)
}

// Timing implements XStats interface
func (xs *xstats) Timing(stat string, duration time.Duration, tags ...string) {
	tags = append(tags, xs.tags...)
	xs.s.Timing(stat, duration, tags...)
}
