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

// XStater is a wrapper around a Sender to inject env tags within all observations.
type XStater interface {
	Sender

	// AddTag adds a tag to the request client, this tag will be sent with all subsequent
	// stats queries.
	AddTags(tags ...string)
}

// Copier is an interface to an xstater that support coping
type Copier interface {
	Copy() XStater
}

// New returns a new xstats client with the provided backend sender.
func New(s Sender) XStater {
	return &xstats{s: s}
}

// NewPrefix returns a new xstats client with the provided backend sender.
// The prefix is prepended to all metric names.
func NewPrefix(s Sender, prefix string) XStater {
	return &xstats{
		s:      s,
		prefix: prefix,
	}
}

// Copy makes a copy of the given xstater if it implements the Copier
// interface. Otherwise it returns a nop stats.
func Copy(xs XStater) XStater {
	if c, ok := xs.(Copier); ok {
		return c.Copy()
	}
	return nop
}

type xstats struct {
	s Sender
	// tags are appended to the tags provided to commands
	tags []string
	// prefix is prepended to all metric
	prefix string
}

// Copy makes a copy of the xstats client
func (xs *xstats) Copy() XStater {
	return &xstats{
		s:      xs.s,
		tags:   xs.tags,
		prefix: xs.prefix,
	}
}

// AddTag implements XStats interface
func (xs *xstats) AddTags(tags ...string) {
	xs.tags = append(xs.tags, tags...)
}

// Gauge implements XStats interface
func (xs *xstats) Gauge(stat string, value float64, tags ...string) {
	if xs.s == nil {
		return
	}
	tags = append(tags, xs.tags...)
	xs.s.Gauge(xs.prefix+stat, value, tags...)
}

// Count implements XStats interface
func (xs *xstats) Count(stat string, count float64, tags ...string) {
	if xs.s == nil {
		return
	}
	tags = append(tags, xs.tags...)
	xs.s.Count(xs.prefix+stat, count, tags...)
}

// Histogram implements XStats interface
func (xs *xstats) Histogram(stat string, value float64, tags ...string) {
	if xs.s == nil {
		return
	}
	tags = append(tags, xs.tags...)
	xs.s.Histogram(xs.prefix+stat, value, tags...)
}

// Timing implements XStats interface
func (xs *xstats) Timing(stat string, duration time.Duration, tags ...string) {
	if xs.s == nil {
		return
	}
	tags = append(tags, xs.tags...)
	xs.s.Timing(xs.prefix+stat, duration, tags...)
}
