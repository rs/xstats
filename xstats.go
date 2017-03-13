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

import (
	"strings"
	"sync"
	"time"
)

// XStater is a wrapper around a Sender to inject env tags within all observations.
type XStater interface {
	Sender

	// AddTags adds the tags to the request client, this tag will be sent with all
	// subsequent stats queries.
	AddTags(tags ...string)

	// AddTag adds a tag to the request client, this tag will be sent with all
	// subsequent stats queries.
	AddTag(key, value string)

	// GetTags return the tags associated with the xstater, all the tags that
	// will be sent along with all the stats queries.
	GetTags() map[string]string
}

// Copier is an interface to an xstater that support coping
type Copier interface {
	Copy() XStater
}

type xstats struct {
	s Sender
	// tags are appended to the tags provided to commands
	mu   sync.RWMutex
	tags map[string]string
	// prefix is prepended to all metric
	prefix string
}

var xstatsPool = sync.Pool{
	New: func() interface{} {
		return &xstats{
			tags: make(map[string]string),
		}
	},
}

// New returns a new xstats client with the provided backend sender.
func New(s Sender) XStater {
	return NewPrefix(s, "")
}

// NewPrefix returns a new xstats client with the provided backend sender.
// The prefix is prepended to all metric names.
func NewPrefix(s Sender, prefix string) XStater {
	xs := xstatsPool.Get().(*xstats)
	xs.s = s
	xs.prefix = prefix
	return xs
}

// Copy makes a copy of the given xstater if it implements the Copier
// interface. Otherwise it returns a nop stats.
func Copy(xs XStater) XStater {
	if c, ok := xs.(Copier); ok {
		return c.Copy()
	}
	return nop
}

// Copy makes a copy of the xstats client
func (xs *xstats) Copy() XStater {
	xs2 := NewPrefix(xs.s, xs.prefix).(*xstats)
	xs.mu.RLock()
	for k, v := range xs.tags {
		xs2.tags[k] = v
	}
	xs.mu.RUnlock()
	return xs2
}

// Close returns the xstats to the sync.Pool.
func (xs *xstats) Close() error {
	xs.mu.Lock()
	defer xs.mu.Unlock()

	xs.s = nil
	xs.tags = make(map[string]string)
	xs.prefix = ""
	xstatsPool.Put(xs)
	return nil
}

// AddTag implements XStats interface
func (xs *xstats) AddTags(tags ...string) {
	xs.mu.Lock()
	defer xs.mu.Unlock()

	for _, tag := range tags {
		tagSlice := strings.Split(tag, ":")
		xs.tags[tagSlice[0]] = tagSlice[1]
	}
}

// AddTag implements XStats interface
func (xs *xstats) AddTag(k, v string) {
	xs.mu.Lock()
	defer xs.mu.Unlock()

	xs.tags[k] = v
}

// AddTag implements XStats interface
func (xs *xstats) GetTags() map[string]string {
	xs.mu.RLock()
	defer xs.mu.RUnlock()

	// copy the tags map so it cannot be altered
	tags := make(map[string]string)
	for k, v := range xs.tags {
		tags[k] = v
	}

	return tags
}

// Gauge implements XStats interface
func (xs *xstats) Gauge(stat string, value float64, tags ...string) {
	if xs.s == nil {
		return
	}
	// copy the tags map so it cannot be altered
	ts := make([]string, 0, len(xs.tags)+len(tags))
	xs.mu.RLock()
	for k, v := range xs.tags {
		ts = append(ts, k+":"+v)
	}
	xs.mu.RUnlock()
	// copy the given tags to it
	ts = append(ts, tags...)

	xs.s.Gauge(xs.prefix+stat, value, ts...)
}

// Count implements XStats interface
func (xs *xstats) Count(stat string, count float64, tags ...string) {
	if xs.s == nil {
		return
	}
	// copy the tags map so it cannot be altered
	ts := make([]string, 0, len(xs.tags)+len(tags))
	xs.mu.RLock()
	for k, v := range xs.tags {
		ts = append(ts, k+":"+v)
	}
	xs.mu.RUnlock()
	// copy the given tags to it
	ts = append(ts, tags...)

	xs.s.Count(xs.prefix+stat, count, ts...)
}

// Histogram implements XStats interface
func (xs *xstats) Histogram(stat string, value float64, tags ...string) {
	if xs.s == nil {
		return
	}
	// copy the tags map so it cannot be altered
	ts := make([]string, 0, len(xs.tags)+len(tags))
	xs.mu.RLock()
	for k, v := range xs.tags {
		ts = append(ts, k+":"+v)
	}
	xs.mu.RUnlock()
	// copy the given tags to it
	ts = append(ts, tags...)

	xs.s.Histogram(xs.prefix+stat, value, ts...)
}

// Timing implements XStats interface
func (xs *xstats) Timing(stat string, duration time.Duration, tags ...string) {
	if xs.s == nil {
		return
	}
	// copy the tags map so it cannot be altered
	ts := make([]string, 0, len(xs.tags)+len(tags))
	xs.mu.RLock()
	for k, v := range xs.tags {
		ts = append(ts, k+":"+v)
	}
	xs.mu.RUnlock()
	// copy the given tags to it
	ts = append(ts, tags...)

	xs.s.Timing(xs.prefix+stat, duration, ts...)
}
