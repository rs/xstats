package prometheus

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/xstats"
)

type sender struct {
	http.Handler

	counters   map[string]*prometheus.CounterVec
	gauges     map[string]*prometheus.GaugeVec
	histograms map[string]*prometheus.HistogramVec
	sync.RWMutex
}

// New creates a prometheus publisher at the given HTTP address.
func New(addr string) xstats.Sender {
	s := NewHandler()
	go func() {
		http.ListenAndServe(addr, s)
	}()
	return s
}

// NewHandler creates a prometheus publisher - a http.Handler and an xstats.Sender.
func NewHandler() *sender {
	return &sender{
		Handler:    prometheus.Handler(),
		counters:   make(map[string]*prometheus.CounterVec),
		gauges:     make(map[string]*prometheus.GaugeVec),
		histograms: make(map[string]*prometheus.HistogramVec),
	}
}

// Gauge implements xstats.Sender interface
//
// Mark the tags as "key:value".
func (s *sender) Gauge(stat string, value float64, tags ...string) {
	s.RLock()
	m, ok := s.gauges[stat]
	s.RUnlock()
	keys, values := splitTags(tags)
	if !ok {
		s.Lock()
		if m, ok = s.gauges[stat]; !ok {
			m = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{Name: stat, Help: stat},
				keys)
			prometheus.MustRegister(m)
			s.gauges[stat] = m
		}
		s.Unlock()
	}
	m.WithLabelValues(values...).Set(value)
}

// Count implements xstats.Sender interface
//
// Mark the tags as "key:value".
func (s *sender) Count(stat string, count float64, tags ...string) {
	s.RLock()
	m, ok := s.counters[stat]
	s.RUnlock()
	keys, values := splitTags(tags)
	if !ok {
		s.Lock()
		if m, ok = s.counters[stat]; !ok {
			m = prometheus.NewCounterVec(
				prometheus.CounterOpts{Name: stat, Help: stat},
				keys)
			prometheus.MustRegister(m)
			s.counters[stat] = m
		}
		s.Unlock()
	}
	m.WithLabelValues(values...).Add(count)
}

// Histogram implements xstats.Sender interface
//
// Mark the tags as "key:value".
func (s *sender) Histogram(stat string, value float64, tags ...string) {
	s.RLock()
	m, ok := s.histograms[stat]
	s.RUnlock()
	keys, values := splitTags(tags)
	if !ok {
		s.Lock()
		if m, ok = s.histograms[stat]; !ok {
			m = prometheus.NewHistogramVec(
				prometheus.HistogramOpts{Name: stat, Help: stat},
				keys)
			prometheus.MustRegister(m)
			s.histograms[stat] = m
		}
		s.Unlock()
	}
	m.WithLabelValues(values...).Observe(value)
}

// Timing implements xstats.Sender interface - simulates Timing with Gauge.
//
// Mark the tags as "key:value".
func (s *sender) Timing(stat string, duration time.Duration, tags ...string) {
	s.Gauge(stat, float64(duration/time.Millisecond), tags...)
}

func splitTags(tags []string) ([]string, []string) {
	keys, values := make([]string, len(tags)), make([]string, len(tags))
	for i, t := range tags {
		if j := strings.IndexByte(t, ':'); j >= 0 {
			keys[i] = t[:j]
			values[i] = t[j+1:]
		} else {
			keys[i] = t
		}
	}
	return keys, values
}
