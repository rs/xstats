package xstats_test

import (
	"log"
	"net"
	"time"

	"github.com/deciphernow/xstats"
	"github.com/deciphernow/xstats/dogstatsd"
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

func ExampleNewScoping() {
	// Defines interval between flushes to statsd server
	flushInterval := 5 * time.Second

	// Connection to the statsd server
	statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
	if err != nil {
		log.Fatal(err)
	}

	// Create the stats client
	s := xstats.NewScoping(dogstatsd.New(statsdWriter, flushInterval), ".", "my-thing")

	// Global tags sent with all metrics (only with supported clients like datadog's)
	s.AddTags("role:my-service", "dc:sv6")

	// Send some observations
	s.Count("requests", 1, "tag")
	s.Timing("something", 5*time.Millisecond, "tag")

	// Scope the client
	ss := xstats.Scope(s, "my-sub-thing")
	ss.Histogram("latency", 50, "tag")
}

func ExampleNewMaxPacket() {
	// Defines interval between flushes to statsd server
	flushInterval := 5 * time.Second

	// Defines the largest packet sent to the statsd server
	maxPacketLen := 8192

	// Connection to the statsd server
	statsdWriter, err := net.Dial("udp", "127.0.0.1:8126")
	if err != nil {
		log.Fatal(err)
	}

	// Create the stats client
	s := xstats.New(dogstatsd.NewMaxPacket(statsdWriter, flushInterval, maxPacketLen))

	// Global tags sent with all metrics (only with supported clients like datadog's)
	s.AddTags("role:my-service", "dc:sv6")

	// Send some observations
	s.Count("requests", 1, "tag")
	s.Timing("something", 5*time.Millisecond, "tag")
}
