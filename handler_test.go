// +build go1.7

package xstats

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	s := &fakeSender{}
	n := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xs, ok := FromRequest(r).(*xstats)
		if assert.True(t, ok) {
			assert.Equal(t, s, xs.s)
			assert.Equal(t, map[string]string{"env": "prod"}, xs.tags)
		}
	})
	h := NewHandler(s, map[string]string{"env": "prod"})(n)
	h.ServeHTTP(nil, &http.Request{})
}

func TestHandlerPrefix(t *testing.T) {
	s := &fakeSender{}
	n := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xs, ok := FromRequest(r).(*xstats)
		if assert.True(t, ok) {
			assert.Equal(t, s, xs.s)
			assert.Equal(t, map[string]string{"env": "prod"}, xs.tags)
			assert.Equal(t, "prefix.", xs.prefix)
		}
	})
	h := NewHandlerPrefix(s, map[string]string{"env": "prod"}, "prefix.")(n)
	h.ServeHTTP(nil, &http.Request{})
}
