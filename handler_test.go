package xstats

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	s := &fakeSender{}
	n := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xs, ok := FromContext(r.Context()).(*xstats)
		assert.True(t, ok)
		assert.Equal(t, s, xs.s)
		assert.Equal(t, []string{"envtag"}, xs.tags)
	})
	h := NewHandler(s, []string{"envtag"})(n)
	req := httptest.NewRequest("", "/", nil)
	h.ServeHTTP(nil, req.WithContext(context.Background()))
}

func TestHandlerPrefix(t *testing.T) {
	s := &fakeSender{}
	n := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xs, ok := FromContext(r.Context()).(*xstats)
		assert.True(t, ok)
		assert.Equal(t, s, xs.s)
		assert.Equal(t, []string{"envtag"}, xs.tags)
		assert.Equal(t, "prefix.", xs.prefix)
	})
	h := NewHandlerPrefix(s, []string{"envtag"}, "prefix.")(n)
	req := httptest.NewRequest("", "/", nil)
	h.ServeHTTP(nil, req.WithContext(context.Background()))
}
