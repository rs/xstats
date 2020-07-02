// Package telegrafstatsd implement telegraf extended StatsD protocol for github.com/deciphernow/xstats
package telegraf

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/deciphernow/xstats"
)

// Inspired by https://github.com/streadway/handy statsd package

type sender struct {
	c    chan string
	quit chan struct{}
	done chan struct{}
}

// defaultMaxPacketLen is the default number of bytes filled before a packet is
// flushed before the reporting interval.
const defaultMaxPacketLen = 1 << 15

var tick = time.Tick

// New creates a telegraf statsd sender that emits observations in the statsd
// protocol to the passed writer. Observations are buffered for the report
// interval or until the buffer exceeds a max packet size, whichever comes
// first.
func New(w io.Writer, reportInterval time.Duration) xstats.Sender {
	return NewMaxPacket(w, reportInterval, defaultMaxPacketLen)
}

// NewMaxPacket creates a telegraf statsd sender that emits observations in the
// statsd protocol to the passed writer. Observations are buffered for the
// report interval or until the buffer exceeds the max packet size, whichever
// comes first.
func NewMaxPacket(w io.Writer, reportInterval time.Duration, maxPacketLen int) xstats.Sender {
	s := &sender{
		c:    make(chan string),
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
	go s.fwd(w, reportInterval, maxPacketLen)
	return s
}

// Gauge implements xstats.Sender interface
func (s *sender) Gauge(stat string, value float64, tags ...string) {
	s.c <- fmt.Sprintf("%s,%s:%f|g\n", stat, t(tags), value)
}

// Count implements xstats.Sender interface
func (s *sender) Count(stat string, count float64, tags ...string) {
	s.c <- fmt.Sprintf("%s,%s:%f|c\n", stat, t(tags), count)
}

// Histogram implements xstats.Sender interface
func (s *sender) Histogram(stat string, value float64, tags ...string) {
	s.c <- fmt.Sprintf("%s,%s:%f|h\n", stat, t(tags), value)
}

// Timing implements xstats.Sender interface
func (s *sender) Timing(stat string, duration time.Duration, tags ...string) {
	s.c <- fmt.Sprintf("%s,%s:%f|ms\n", stat, t(tags), duration.Seconds())
}

// Close implements xstats.Sender interface
func (s *sender) Close() error {
	close(s.quit)
	<-s.done
	close(s.c)

	return nil
}

// Generate a telegraf tag suffix
func t(tags []string) string {
	for i, v := range tags {
		tags[i] = strings.Replace(v, ":", "=", 1)
	}
	t := ""
	if len(tags) > 0 {
		t = "" + strings.Join(tags, ",")
	}
	return t
}

func (s *sender) fwd(w io.Writer, reportInterval time.Duration, maxPacketLen int) {
	defer close(s.done)

	buf := &bytes.Buffer{}
	tick := tick(reportInterval)
	for {
		select {
		case m := <-s.c:
			newLen := buf.Len() + len(m)
			if newLen > maxPacketLen {
				flush(w, buf)
			}

			buf.Write([]byte(m))

			if newLen == maxPacketLen {
				flush(w, buf)
			}

		case <-tick:
			flush(w, buf)
		case <-s.quit:
			flush(w, buf)
			return
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
