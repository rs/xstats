package xstats_test

import (
	"log"
	"net"
	"time"

	"github.com/rs/xstats"
	"github.com/rs/xstats/dogstatsd"
)

func ExampleNew() {
	// Defines interval between flushes to statsd server
	flushInterval := 5 * time.Second

	// Connection to the statsd server
	statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
	if err != nil {
		log.Fatal(err)
	}

	// Create the stats client
	s := xstats.New(dogstatsd.New(statsdWriter, flushInterval))

	// Global tags sent with all metrics (only with supported clients like datadog's)
	s.AddTags("role:my-service", "dc:sv6")

	// Send some observations
	s.Count("requests", 1, "tag")
	s.Timing("something", 5*time.Millisecond, "tag")
}
