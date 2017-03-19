package xstats

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMultiSender(t *testing.T) {
	fs1 := &fakeSender{}
	fs2 := &fakeSendCloser{err: errors.New("foo")}
	fs3 := &fakeSendCloser{err: errors.New("bar")}
	m := MultiSender{fs1, fs2, fs3}

	m.Count("foo", 1, "bar", "baz")
	countCmd := cmd{"Count", "foo", 1, []string{"bar", "baz"}}
	assert.Equal(t, countCmd, fs1.last)
	assert.Equal(t, countCmd, fs2.last)
	assert.Equal(t, countCmd, fs3.last)

	m.Gauge("foo", 1, "bar", "baz")
	gaugeCmd := cmd{"Gauge", "foo", 1, []string{"bar", "baz"}}
	assert.Equal(t, gaugeCmd, fs1.last)
	assert.Equal(t, gaugeCmd, fs2.last)
	assert.Equal(t, gaugeCmd, fs3.last)

	m.Histogram("foo", 1, "bar", "baz")
	histoCmd := cmd{"Histogram", "foo", 1, []string{"bar", "baz"}}
	assert.Equal(t, histoCmd, fs1.last)
	assert.Equal(t, histoCmd, fs2.last)
	assert.Equal(t, histoCmd, fs3.last)

	m.Timing("foo", 1*time.Second, "bar", "baz")
	timingCmd := cmd{"Timing", "foo", 1, []string{"bar", "baz"}}
	assert.Equal(t, timingCmd, fs1.last)
	assert.Equal(t, timingCmd, fs2.last)
	assert.Equal(t, timingCmd, fs3.last)

	assert.Equal(t, fs2.err, CloseSender(m))
	assert.Equal(t, timingCmd, fs1.last)
	assert.Equal(t, cmd{name: "Close"}, fs2.last)
	assert.Equal(t, cmd{name: "Close"}, fs3.last)
}
