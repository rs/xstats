package expvar

import (
	"expvar"
	"strconv"
	"time"

	"github.com/rs/xstats"
)

type sender struct {
	vars *expvar.Map
}

// A expvar.Var static float
type float float64

// String implements the expvar.Var
func (f float) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

// New creates a statsd sender that publish observations in expvar under
// the given prefix "path". Will panic if the prefix is already used.
//
// Tags are ignored. Histogram and Timing methods are not supported.
func New(prefix string) xstats.Sender {
	return &sender{expvar.NewMap(prefix)}
}

// Gauge implements xstats.Sender interface
func (s sender) Gauge(stat string, value float64, tags ...string) {
	s.vars.Set(stat, float(value))
}

// Count implements xstats.Sender interface
func (s sender) Count(stat string, count float64, tags ...string) {
	s.vars.AddFloat(stat, count)
}

// Histogram implements xstats.Sender interface
func (s sender) Histogram(stat string, value float64, tags ...string) {
	// Not supported, just ignored
}

// Timing implements xstats.Sender interface
func (s sender) Timing(stat string, duration time.Duration, tags ...string) {
	// Not supported, just ignored
}
