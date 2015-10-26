package expvar

import (
	"expvar"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Publishes prefix in expvar, panics the second time
	assert.Nil(t, expvar.Get("name"))
	New("name")
	assert.NotNil(t, expvar.Get("name"))
	assert.IsType(t, expvar.Get("name"), &expvar.Map{})
	assert.Panics(t, func() {
		New("name")
	})
}

func TestGauge(t *testing.T) {
	s := New("gauge")
	v := expvar.Get("gauge").(*expvar.Map)
	s.Gauge("test", 1)
	assert.Equal(t, "1", v.Get("test").String())
	s.Gauge("test", -1)
	assert.Equal(t, "-1", v.Get("test").String())
}

func TestCount(t *testing.T) {
	s := New("count")
	v := expvar.Get("count").(*expvar.Map)
	s.Count("test", 1)
	assert.Equal(t, "1", v.Get("test").String())
	s.Count("test", -1)
	assert.Equal(t, "0", v.Get("test").String())
}

func TestHistogram(t *testing.T) {
	s := New("histogram")
	s.Histogram("test", 1)
}

func TestTiming(t *testing.T) {
	s := New("timing")
	s.Timing("test", 1)
}
