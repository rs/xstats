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
package xstats // import "github.com/deciphernow/xstats"

import (
	"io"
	"strings"
	"sync"
	"time"
)

const (
	defaultDelimiter = "."
)

var (
	// DisablePooling will disable the use of sync.Pool fo resource management when true.
	// This allows for XStater instances to persist beyond the scope of an HTTP request
	// handler. However, using this option puts a greater pressure on GC and changes
	// the memory usage patterns of the library. Use only if there is a requirement
	// for persistent stater references.
	DisablePooling = false
)

// XStater is a wrapper around a Sender to inject env tags within all observations.
type XStater interface {
	Sender

	// AddTag adds a tag to the request client, this tag will be sent with all
	// subsequent stats queries.
	AddTags(tags ...string)

	// GetTags returns the tags associated with the XStater, all the tags that
	// will be sent along with all the stats queries.
	GetTags() []string
}

// Copier is an interface to an XStater that supports coping
type Copier interface {
	Copy() XStater
}

// Scoper is an interface to an XStater, that supports scoping
type Scoper interface {
	Scope(scope string, scopes ...string) XStater
}

var xstatsPool = &sync.Pool{
	New: func() interface{} {
		return &xstats{}
	},
}

// New returns a new xstats client with the provided backend sender.
func New(s Sender) XStater {
	return NewPrefix(s, "")
}

// NewPrefix returns a new xstats client with the provided backend sender.
// The prefix is prepended to all metric names.
func NewPrefix(s Sender, prefix string) XStater {
	return NewScoping(s, "", prefix)
}

// NewScoping returns a new xstats client with the provided backend sender.
// The delimiter is used to delimit scopes. Initial scopes can be provided.
func NewScoping(s Sender, delimiter string, scopes ...string) XStater {
	var xs *xstats
	if DisablePooling {
		xs = &xstats{}
	} else {
		xs = xstatsPool.Get().(*xstats)
	}
	xs.s = s
	if len(scopes) > 0 {
		xs.prefix = strings.Join(scopes, delimiter) + delimiter
	} else {
		xs.prefix = ""
	}
	xs.delimiter = delimiter
	return xs
}

// Copy makes a copy of the given XStater if it implements the Copier
// interface. Otherwise it returns a nop stats.
func Copy(xs XStater) XStater {
	if c, ok := xs.(Copier); ok {
		return c.Copy()
	}
	return nop
}

// Scope makes a scoped copy of the given XStater if it implements the Scoper
// interface. Otherwise it returns a nop stats.
func Scope(xs XStater, scope string, scopes ...string) XStater {
	if c, ok := xs.(Scoper); ok {
		return c.Scope(scope, scopes...)
	}
	return nop
}

// Close will call Close() on any xstats.XStater that implements io.Closer
func Close(xs XStater) error {
	if c, ok := xs.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

type xstats struct {
	s Sender
	// tags are appended to the tags provided to commands
	tags []string
	// prefix is prepended to all metric
	prefix string
	// delimiter is used to delimit scopes
	delimiter string
}

// Copy implements the Copier interface
func (xs *xstats) Copy() XStater {
	xs2 := NewScoping(xs.s, xs.delimiter, xs.prefix).(*xstats)
	xs2.tags = xs.tags
	return xs2
}

// Scope implements Scoper interface
func (xs *xstats) Scope(scope string, scopes ...string) XStater {
	var scs []string
	if xs.prefix == "" {
		scs = make([]string, 0, 1+len(scopes))
	} else {
		scs = make([]string, 0, 2+len(scopes))
		scs = append(scs, strings.TrimRight(xs.prefix, xs.delimiter))
	}
	scs = append(scs, scope)
	scs = append(scs, scopes...)
	xs2 := NewScoping(xs.s, xs.delimiter, scs...).(*xstats)
	xs2.tags = xs.tags
	return xs2
}

// Close returns the xstats to the sync.Pool
func (xs *xstats) Close() error {
	if !DisablePooling {
		xs.s = nil
		xs.tags = nil
		xs.prefix = ""
		xs.delimiter = ""
		xstatsPool.Put(xs)
	}
	return nil
}

// AddTag implements XStater interface
func (xs *xstats) AddTags(tags ...string) {
	if xs.tags == nil {
		xs.tags = tags
	} else {
		xs.tags = append(xs.tags, tags...)
	}
}

// AddTag implements XStater interface
func (xs *xstats) GetTags() []string {
	return xs.tags
}

// Gauge implements XStater interface
func (xs *xstats) Gauge(stat string, value float64, tags ...string) {
	if xs.s == nil {
		return
	}
	tags = append(tags, xs.tags...)
	xs.s.Gauge(xs.prefix+stat, value, tags...)
}

// Count implements XStater interface
func (xs *xstats) Count(stat string, count float64, tags ...string) {
	if xs.s == nil {
		return
	}
	tags = append(tags, xs.tags...)
	xs.s.Count(xs.prefix+stat, count, tags...)
}

// Histogram implements XStater interface
func (xs *xstats) Histogram(stat string, value float64, tags ...string) {
	if xs.s == nil {
		return
	}
	tags = append(tags, xs.tags...)
	xs.s.Histogram(xs.prefix+stat, value, tags...)
}

// Timing implements XStater interface
func (xs *xstats) Timing(stat string, duration time.Duration, tags ...string) {
	if xs.s == nil {
		return
	}
	tags = append(tags, xs.tags...)
	xs.s.Timing(xs.prefix+stat, duration, tags...)
}
