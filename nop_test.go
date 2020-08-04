package xstats

import (
	"testing"
	"time"
)

func TestNop(t *testing.T) {
	Nop.AddTags("tag")
	Nop.Gauge("metric", 1)
	Nop.Count("metric", 1)
	Nop.Histogram("metric", 1)
	Nop.Timing("metric", 1*time.Second)
}
