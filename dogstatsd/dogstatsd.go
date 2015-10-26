// Package dogstatsd implement Datadog extended StatsD protocol for github.com/rs/xstats
package dogstatsd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/rs/xstats"
)

// Inspired by https://github.com/streadway/handy statsd package

type sender chan string

// MaxPacketLen is the number of bytes filled before a packet is flushed before
// the reporting interval.
const maxPacketLen = 2 ^ 15

var tick = time.Tick

// New creates a datadog statsd sender that emit observations in the statsd
// protocol to the passed writer. Observations are buffered for the report
// interval or until the buffer exceeds a max packet size, whichever comes
// first.
func New(w io.Writer, reportInterval time.Duration) xstats.Sender {
	s := make(chan string)
	go fwd(w, reportInterval, s)
	return sender(s)
}

// Gauge implements xstats.Sender interface
func (s sender) Gauge(stat string, value float64, tags ...string) {
	s <- fmt.Sprintf("%s:%f|g%s\n", stat, value, t(tags))
}

// Count implements xstats.Sender interface
func (s sender) Count(stat string, count float64, tags ...string) {
	s <- fmt.Sprintf("%s:%f|c%s\n", stat, count, t(tags))
}

// Histogram implements xstats.Sender interface
func (s sender) Histogram(stat string, value float64, tags ...string) {
	s <- fmt.Sprintf("%s:%f|h%s\n", stat, value, t(tags))
}

// Timing implements xstats.Sender interface
func (s sender) Timing(stat string, duration time.Duration, tags ...string) {
	s <- fmt.Sprintf("%s:%f|ms%s\n", stat, duration.Seconds(), t(tags))
}

// Generate a DogStatsD tag suffix
func t(tags []string) string {
	t := ""
	if len(tags) > 0 {
		t = "|#" + strings.Join(tags, ",")
	}
	return t
}

func fwd(w io.Writer, reportInterval time.Duration, c <-chan string) {
	buf := &bytes.Buffer{}
	tick := tick(reportInterval)
	for {
		select {
		case s := <-c:
			buf.Write([]byte(s))
			if buf.Len() > maxPacketLen {
				flush(w, buf)
			}

		case <-tick:
			flush(w, buf)
		}
	}
}

func flush(w io.Writer, buf *bytes.Buffer) {
	if buf.Len() <= 0 {
		return
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("error: could not write to statsd: %v", err)
	}
	buf.Reset()
}
