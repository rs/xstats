package xstats_test

import (
	"net"
	"time"

	"github.com/deciphernow/xstats"
	"github.com/deciphernow/xstats/dogstatsd"
	"github.com/deciphernow/xstats/expvar"
)

func ExampleMultiSender() {
	// Create an expvar sender
	s1 := expvar.New("stats")

	// Create the stats sender
	statsdWriter, _ := net.Dial("udp", "127.0.0.1:8126")
	s2 := dogstatsd.New(statsdWriter, 5*time.Second)

	// Create a xstats with a sender composed of the previous two.
	// You may also create a NewHandler() the same way.
	s := xstats.New(xstats.MultiSender{s1, s2})

	// Send some observations
	s.Count("requests", 1, "tag")
	s.Timing("something", 5*time.Millisecond, "tag")
}
