// Package statsd implement the StatsD protocol for github.com/deciphernow/xstats
package statsd

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

// New creates a statsd sender that emits observations in the statsd
// protocol to the passed writer. Observations are buffered for the report
// interval or until the buffer exceeds a max packet size, whichever comes
// first.
func New(w io.Writer, reportInterval time.Duration) xstats.Sender {
	return NewMaxPacket(w, reportInterval, defaultMaxPacketLen)
}

// NewMaxPacket creates a statsd sender that emits observations in the statsd
// protocol to the passed writer. Observations are buffered for the report
// interval or until the buffer exceeds the max packet size, whichever comes
// first. Tags are ignored.
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
	s.c <- fmt.Sprintf("%s:%f|g\n", stat, value)
}

// Count implements xstats.Sender interface
func (s *sender) Count(stat string, count float64, tags ...string) {
	s.c <- fmt.Sprintf("%s:%f|c\n", stat, count)
}

// Histogram implements xstats.Sender interface
func (s *sender) Histogram(stat string, value float64, tags ...string) {
	s.c <- fmt.Sprintf("%s:%f|h\n", stat, value)
}

// Timing implements xstats.Sender interface
func (s *sender) Timing(stat string, duration time.Duration, tags ...string) {
	s.c <- fmt.Sprintf("%s:%f|ms\n", stat, duration.Seconds()*1000)
}

// Close implements xstats.Sender interface
func (s *sender) Close() error {
	close(s.quit)
	<-s.done
	close(s.c)

	return nil
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
