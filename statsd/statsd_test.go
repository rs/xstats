package statsd

import (
	"bytes"
	"errors"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var tickC = make(chan time.Time)
var fakeTick = func(time.Duration) <-chan time.Time { return tickC }

func TestCounter(t *testing.T) {
	tick = fakeTick
	defer func() { tick = time.Tick }()

	buf := &bytes.Buffer{}
	c := New(buf, time.Second)

	c.Count("metric1", 1, "tag1")
	c.Count("metric2", 2, "tag1", "tag2")
	tickC <- time.Now()
	runtime.Gosched()

	assert.Equal(t, "metric1:1.000000|c\nmetric2:2.000000|c\n", buf.String())
}

func TestGauge(t *testing.T) {
	tick = fakeTick
	defer func() { tick = time.Tick }()

	buf := &bytes.Buffer{}
	c := New(buf, time.Second)

	c.Gauge("metric1", 1, "tag1")
	c.Gauge("metric2", -2.0, "tag1", "tag2")
	tickC <- time.Now()
	runtime.Gosched()

	assert.Equal(t, "metric1:1.000000|g\nmetric2:-2.000000|g\n", buf.String())
}

func TestHistogram(t *testing.T) {
	tick = fakeTick
	defer func() { tick = time.Tick }()

	buf := &bytes.Buffer{}
	c := New(buf, time.Second)

	c.Histogram("metric1", 1, "tag1")
	c.Histogram("metric2", 2, "tag1", "tag2")
	tickC <- time.Now()
	runtime.Gosched()

	assert.Equal(t, "metric1:1.000000|h\nmetric2:2.000000|h\n", buf.String())
}

func TestTiming(t *testing.T) {
	tick = fakeTick
	defer func() { tick = time.Tick }()

	buf := &bytes.Buffer{}
	c := New(buf, time.Second)

	c.Timing("metric1", time.Second, "tag1")
	c.Timing("metric2", 2*time.Second, "tag1", "tag2")
	tickC <- time.Now()
	runtime.Gosched()

	assert.Equal(t, "metric1:1.000000|ms\nmetric2:2.000000|ms\n", buf.String())
}

type errWriter struct{}

func (w errWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("i/o error")
}

func TestInvalidBuffer(t *testing.T) {
	tick = fakeTick
	defer func() { tick = time.Tick }()

	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	defer func() { log.SetOutput(os.Stderr) }()

	c := New(&errWriter{}, time.Second)

	c.Count("metric", 1)
	tickC <- time.Now()
	runtime.Gosched()

	assert.True(t, strings.HasSuffix(buf.String(), "error: could not write to statsd: i/o error\n"))
}
