package xstats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultiSender(t *testing.T) {
	fs1 := &fakeSender{}
	fs2 := &fakeSender{}
	m := MultiSender{fs1, fs2}
	m.Count("foo", 1, "bar", "baz")
	assert.Equal(t, cmd{"Count", "foo", 1, []string{"bar", "baz"}}, fs1.last)
	assert.Equal(t, cmd{"Count", "foo", 1, []string{"bar", "baz"}}, fs2.last)
	m.Gauge("foo", 1, "bar", "baz")
	assert.Equal(t, cmd{"Gauge", "foo", 1, []string{"bar", "baz"}}, fs1.last)
	assert.Equal(t, cmd{"Gauge", "foo", 1, []string{"bar", "baz"}}, fs2.last)
	m.Histogram("foo", 1, "bar", "baz")
	assert.Equal(t, cmd{"Histogram", "foo", 1, []string{"bar", "baz"}}, fs1.last)
	assert.Equal(t, cmd{"Histogram", "foo", 1, []string{"bar", "baz"}}, fs2.last)
	m.Timing("foo", 1*time.Second, "bar", "baz")
	assert.Equal(t, cmd{"Timing", "foo", 1, []string{"bar", "baz"}}, fs1.last)
	assert.Equal(t, cmd{"Timing", "foo", 1, []string{"bar", "baz"}}, fs2.last)
}
