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
	ctx = NewContext(ctx, xs)
	ctxs := FromContext(ctx)
	assert.Equal(t, xs, ctxs)
}

func TestNew(t *testing.T) {
	xs := New(&fakeSender{})
	_, ok := xs.(*xstats)
	assert.True(t, ok)
}

func TestNewPrefix(t *testing.T) {
	xs := NewPrefix(&fakeSender{}, "prefix.")
	x, ok := xs.(*xstats)
	assert.True(t, ok)
	assert.Equal(t, "prefix.", x.prefix)
}

func TestCopy(t *testing.T) {
	xs := NewPrefix(&fakeSender{}, "prefix.").(*xstats)
	xs.AddTags("k1:v1")
	xs2 := Copy(xs).(*xstats)
	assert.Equal(t, xs.s, xs2.s)
	assert.Equal(t, xs.tags, xs2.tags)
	assert.Equal(t, xs.prefix, xs2.prefix)
	xs2.AddTags("k2:v2", "k3:v3")
	assert.Equal(t, map[string]string{"k1": "v1"}, xs.tags)
	assert.Equal(t, map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"}, xs2.tags)

	assert.Equal(t, nop, Copy(nop))
	assert.Equal(t, nop, Copy(nil))
}

func TestAddTag(t *testing.T) {
	xs := &xstats{s: &fakeSender{}, tags: make(map[string]string)}
	xs.AddTag("k1", "v1")
	assert.Equal(t, map[string]string{"k1": "v1"}, xs.tags)
}

func TestAddTags(t *testing.T) {
	xs := &xstats{s: &fakeSender{}, tags: make(map[string]string)}
	xs.AddTags("k1:v1")
	assert.Equal(t, map[string]string{"k1": "v1"}, xs.tags)
}

func TestGauge(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s, prefix: "p.", tags: make(map[string]string)}
	xs.AddTags("k1:v1")
	xs.Gauge("bar", 1, "k2:v2")
	assert.Equal(t, cmd{"Gauge", "p.bar", 1, []string{"k1:v1", "k2:v2"}}, s.last)
}

func TestCount(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s, prefix: "p.", tags: make(map[string]string)}
	xs.AddTags("k1:v1")
	xs.Count("bar", 1, "k2:v2")
	assert.Equal(t, cmd{"Count", "p.bar", 1, []string{"k1:v1", "k2:v2"}}, s.last)
}

func TestHistogram(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s, prefix: "p.", tags: make(map[string]string)}
	xs.AddTags("k1:v1")
	xs.Histogram("bar", 1, "k2:v2")
	assert.Equal(t, cmd{"Histogram", "p.bar", 1, []string{"k1:v1", "k2:v2"}}, s.last)
}

func TestTiming(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s, prefix: "p.", tags: make(map[string]string)}
	xs.AddTags("k1:v1")
	xs.Timing("bar", 1, "k2:v2")
	assert.Equal(t, cmd{"Timing", "p.bar", 1 / float64(time.Second), []string{"k1:v1", "k2:v2"}}, s.last)
}

func TestNilSender(t *testing.T) {
	xs := &xstats{}
	xs.Gauge("foo", 1)
	xs.Count("foo", 1)
	xs.Histogram("foo", 1)
	xs.Timing("foo", 1)
}
