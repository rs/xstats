package xstats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type fakeSender struct {
	last cmd
}

type fakeSendCloser struct {
	fakeSender
	err error
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

func (s *fakeSendCloser) Close() error {
	s.fakeSender.last = cmd{name: "Close"}
	return s.err
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

func TestNewScoping(t *testing.T) {
	xs := NewScoping(&fakeSender{}, "/")
	x, ok := xs.(*xstats)
	assert.True(t, ok)
	assert.Equal(t, "", x.prefix)

	xs = NewScoping(&fakeSender{}, "/", "prefix", "infix", "suffix")
	x, ok = xs.(*xstats)
	assert.True(t, ok)
	assert.Equal(t, "prefix/infix/suffix/", x.prefix)
}

func TestCopy(t *testing.T) {
	xs := NewPrefix(&fakeSender{}, "prefix.").(*xstats)
	xs.AddTags("foo")
	xs2 := Copy(xs).(*xstats)
	assert.Equal(t, xs.s, xs2.s)
	assert.Equal(t, xs.tags, xs2.tags)
	assert.Equal(t, xs.prefix, xs2.prefix)
	xs2.AddTags("bar", "baz")
	assert.Equal(t, []string{"foo"}, xs.tags)
	assert.Equal(t, []string{"foo", "bar", "baz"}, xs2.tags)

	assert.Equal(t, nop, Copy(nop))
	assert.Equal(t, nop, Copy(nil))
}

func TestScope(t *testing.T) {
	xs := NewScoping(&fakeSender{}, ".").(*xstats)
	xs.AddTags("foo")

	xs2 := Scope(xs, "prefix").(*xstats)
	assert.Equal(t, xs.s, xs2.s)
	assert.Equal(t, xs.tags, xs2.tags)
	assert.Equal(t, "prefix.", xs2.prefix)

	xs3 := Scope(xs2, "infix", "suffix").(*xstats)
	assert.Equal(t, xs2.s, xs3.s)
	assert.Equal(t, xs2.tags, xs3.tags)
	assert.Equal(t, "prefix.infix.suffix.", xs3.prefix)

	xs2.AddTags("bar", "baz")
	xs3.AddTags("blegga")
	assert.Equal(t, []string{"foo"}, xs.tags)
	assert.Equal(t, []string{"foo", "bar", "baz"}, xs2.tags)
	assert.Equal(t, []string{"foo", "blegga"}, xs3.tags)

	assert.Equal(t, nop, Scope(nop, "prefix"))
	assert.Equal(t, nop, Scope(nil, "prefix"))
}

func TestAddTag(t *testing.T) {
	xs := &xstats{s: &fakeSender{}}
	xs.AddTags("foo")
	assert.Equal(t, []string{"foo"}, xs.tags)
}

func TestGauge(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s, prefix: "p."}
	xs.AddTags("foo")
	xs.Gauge("bar", 1, "baz")
	assert.Equal(t, cmd{"Gauge", "p.bar", 1, []string{"baz", "foo"}}, s.last)
}

func TestCount(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s, prefix: "p."}
	xs.AddTags("foo")
	xs.Count("bar", 1, "baz")
	assert.Equal(t, cmd{"Count", "p.bar", 1, []string{"baz", "foo"}}, s.last)
}

func TestHistogram(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s, prefix: "p."}
	xs.AddTags("foo")
	xs.Histogram("bar", 1, "baz")
	assert.Equal(t, cmd{"Histogram", "p.bar", 1, []string{"baz", "foo"}}, s.last)
}

func TestTiming(t *testing.T) {
	s := &fakeSender{}
	xs := &xstats{s: s, prefix: "p."}
	xs.AddTags("foo")
	xs.Timing("bar", 1, "baz")
	assert.Equal(t, cmd{"Timing", "p.bar", 1 / float64(time.Second), []string{"baz", "foo"}}, s.last)
}

func TestNilSender(t *testing.T) {
	xs := &xstats{}
	xs.Gauge("foo", 1)
	xs.Count("foo", 1)
	xs.Histogram("foo", 1)
	xs.Timing("foo", 1)
}
