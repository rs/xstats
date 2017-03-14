package xstats

import (
	"testing"
	"time"
)

func TestNop(t *testing.T) {
	nop.AddTags("tag:value")
	nop.Gauge("metric", 1)
	nop.Count("metric", 1)
	nop.Histogram("metric", 1)
	nop.Timing("metric", 1*time.Second)
}
