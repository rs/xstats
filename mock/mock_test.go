package mock

import (
	"testing"
	"time"

	"github.com/deciphernow/xstats"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := New()
	assert.Implements(t, (*xstats.Sender)(nil), s)
}

func TestSender(t *testing.T) {
	s := New()
	d := time.Second
	tags := []string{"2", "3"}
	s.On("Timing", "1", d, tags)
	s.Timing("1", d, tags...)
	s.On("Count", "1", 2.0, tags)
	s.Count("1", 2.0, tags...)
	s.On("Gauge", "1", 2.0, tags)
	s.Gauge("1", 2.0, tags...)
	s.On("Histogram", "1", 2.0, tags)
	s.Histogram("1", 2.0, tags...)
	s.AssertExpectations(t)
}
