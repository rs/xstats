package xstats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type fakeSender struct {
	last cmd
	err  error
}

type cmd struct {
	name  string
	stat  string
	value float64
	tags  []string
}

func (s *fakeSender) Gauge(stat string, value float64, tags ...string) {
	s.last = cmd{"Gauge", stat, value, tags}
}

func (s *fakeSender) Count(stat string, count float64, tags ...string) {
	s.last = cmd{"Count", stat, count, tags}
}

func (s *fakeSender) Histogram(stat string, value float64, tags ...string) {
	s.last = cmd{"Histogram", stat, value, tags}
}

func (s *fakeSender) Timing(stat string, duration time.Duration, tags ...string) {
	s.last = cmd{"Timing", stat, duration.Seconds(), tags}
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	s := FromContext(ctx)
	assert.Equal(t, nop, s)

	ctx = context.Background()
	xs := &xstats{}
	ctx = newContext(ctx, xs)
	ctxs := FromContext(ctx)
	assert.Equal(t, xs, ctxs)
}

func TestNew(t *testing.T) {
	xs := New(&fakeSender{})
	_, ok := xs.(*xstats)
	assert.True(t, ok)
}

func TestAddTag(t *testing.T) {
	xs := &xstats{s: &fakeSender{}}
	xs.AddTags("foo")
	assert.Equal(t, []string{"foo"}, xs.tags)
}

func TestGauge(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s}
	xs.AddTags("foo")
	xs.Gauge("bar", 1, "baz")
	assert.Equal(t, cmd{"Gauge", "bar", 1, []string{"baz", "foo"}}, s.last)
}

func TestCount(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s}
	xs.AddTags("foo")
	xs.Count("bar", 1, "baz")
	assert.Equal(t, cmd{"Count", "bar", 1, []string{"baz", "foo"}}, s.last)
}

func TestHistogram(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s}
	xs.AddTags("foo")
	xs.Histogram("bar", 1, "baz")
	assert.Equal(t, cmd{"Histogram", "bar", 1, []string{"baz", "foo"}}, s.last)
}

func TestTiming(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s}
	xs.AddTags("foo")
	xs.Timing("bar", 1, "baz")
	assert.Equal(t, cmd{"Timing", "bar", 1 / float64(time.Second), []string{"baz", "foo"}}, s.last)
}
