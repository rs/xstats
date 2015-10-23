package xstats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

type fakeClient struct {
	last cmd
	err  error
}

type cmd struct {
	name  string
	stat  string
	value float64
	tags  []string
}

func (c *fakeClient) Gauge(stat string, value float64, tags ...string) {
	c.last = cmd{"Gauge", stat, value, tags}
}

func (c *fakeClient) Count(stat string, count float64, tags ...string) {
	c.last = cmd{"Count", stat, count, tags}
}

func (c *fakeClient) Histogram(stat string, value float64, tags ...string) {
	c.last = cmd{"Histogram", stat, value, tags}
}

func (c *fakeClient) Timing(stat string, duration time.Duration, tags ...string) {
	c.last = cmd{"Timing", stat, duration.Seconds(), tags}
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	s := FromContext(ctx)
	assert.Equal(t, nopClient, s)

	ctx = context.Background()
	rc := &requestClient{}
	ctx = newContext(ctx, rc)
	ctxrc := FromContext(ctx)
	assert.Equal(t, rc, ctxrc)
}

func TestNew(t *testing.T) {
	rc := New(&fakeClient{})
	_, ok := rc.(*requestClient)
	assert.True(t, ok)
}

func TestAddTag(t *testing.T) {
	rc := &requestClient{c: &fakeClient{}}
	rc.AddTags("foo")
	assert.Equal(t, []string{"foo"}, rc.tags)
}

func TestGauge(t *testing.T) {
	c := &fakeClient{}
	rc := &requestClient{c: c}
	rc.AddTags("foo")
	rc.Gauge("bar", 1, "baz")
	assert.Equal(t, cmd{"Gauge", "bar", 1, []string{"baz", "foo"}}, c.last)
}

func TestCount(t *testing.T) {
	c := &fakeClient{}
	rc := &requestClient{c: c}
	rc.AddTags("foo")
	rc.Count("bar", 1, "baz")
	assert.Equal(t, cmd{"Count", "bar", 1, []string{"baz", "foo"}}, c.last)
}

func TestHistogram(t *testing.T) {
	c := &fakeClient{}
	rc := &requestClient{c: c}
	rc.AddTags("foo")
	rc.Histogram("bar", 1, "baz")
	assert.Equal(t, cmd{"Histogram", "bar", 1, []string{"baz", "foo"}}, c.last)
}

func TestTiming(t *testing.T) {
	c := &fakeClient{}
	rc := &requestClient{c: c}
	rc.AddTags("foo")
	rc.Timing("bar", 1, "baz")
	assert.Equal(t, cmd{"Timing", "bar", 1 / float64(time.Second), []string{"baz", "foo"}}, c.last)
}
