package prometheus

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func get(w io.Writer, h http.Handler, mark byte) {
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, &http.Request{Method: "GET", URL: &url.URL{Path: "/metrics"}})
	if rr.Code > 299 {
		fmt.Fprintf(w, "%d\n%v\n\n", rr.Code, rr.HeaderMap)
		io.Copy(w, rr.Body)
		return
	}
	scanner := bufio.NewScanner(rr.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, []byte("metric")) &&
			line[bytes.IndexByte(line, '_')+1] == mark {
			w.Write(line)
			w.Write([]byte{'\n'})
		}
	}
}

func TestCounter(t *testing.T) {
	c := NewHandler()
	c.Count("metric1_c", 1, "tag:1")
	c.Count("metric2_c", 2, "tag:1", "gat:2")
	buf := &bytes.Buffer{}
	get(buf, c, 'c')

	assert.Equal(t, "metric1_c{tag=\"1\"} 1\nmetric2_c{gat=\"2\",tag=\"1\"} 2\n", buf.String())
}

func TestGauge(t *testing.T) {
	c := NewHandler()
	c.Gauge("metric1_g", 1, "tag:1")
	c.Gauge("metric2_g", -2.0, "tag:1", "gat:2")
	buf := &bytes.Buffer{}
	get(buf, c, 'g')

	assert.Equal(t, "metric1_g{tag=\"1\"} 1\nmetric2_g{gat=\"2\",tag=\"1\"} -2\n", buf.String())
}

func TestHistogram(t *testing.T) {
	c := NewHandler()
	c.Histogram("metric1_h", 1, "tag:1")
	c.Histogram("metric2_h", 2, "tag:1", "gat:2")
	buf := &bytes.Buffer{}
	get(buf, c, 'h')

	assert.Equal(t, "metric1_h_bucket{tag=\"1\",le=\"0.005\"} 0\nmetric1_h_bucket{tag=\"1\",le=\"0.01\"} 0\nmetric1_h_bucket{tag=\"1\",le=\"0.025\"} 0\nmetric1_h_bucket{tag=\"1\",le=\"0.05\"} 0\nmetric1_h_bucket{tag=\"1\",le=\"0.1\"} 0\nmetric1_h_bucket{tag=\"1\",le=\"0.25\"} 0\nmetric1_h_bucket{tag=\"1\",le=\"0.5\"} 0\nmetric1_h_bucket{tag=\"1\",le=\"1\"} 1\nmetric1_h_bucket{tag=\"1\",le=\"2.5\"} 1\nmetric1_h_bucket{tag=\"1\",le=\"5\"} 1\nmetric1_h_bucket{tag=\"1\",le=\"10\"} 1\nmetric1_h_bucket{tag=\"1\",le=\"+Inf\"} 1\nmetric1_h_sum{tag=\"1\"} 1\nmetric1_h_count{tag=\"1\"} 1\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"0.005\"} 0\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"0.01\"} 0\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"0.025\"} 0\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"0.05\"} 0\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"0.1\"} 0\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"0.25\"} 0\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"0.5\"} 0\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"1\"} 0\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"2.5\"} 1\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"5\"} 1\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"10\"} 1\nmetric2_h_bucket{gat=\"2\",tag=\"1\",le=\"+Inf\"} 1\nmetric2_h_sum{gat=\"2\",tag=\"1\"} 2\nmetric2_h_count{gat=\"2\",tag=\"1\"} 1\n", buf.String())
}

func TestTiming(t *testing.T) {
	c := NewHandler()
	c.Timing("metric1_t", time.Second, "tag:1")
	c.Timing("metric2_t", 2*time.Second, "tag:1", "gat:2")
	buf := &bytes.Buffer{}
	get(buf, c, 't')

	assert.Equal(t, "metric1_t{tag=\"1\"} 1000\nmetric2_t{gat=\"2\",tag=\"1\"} 2000\n", buf.String())
}
