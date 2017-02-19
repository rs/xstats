// +build go1.7

package xstats

import (
	"net/http"

	"context"
)

// Handler injects a per request metrics client in the net/context which can be
// retrived using xstats.FromContext(ctx)
type Handler struct {
	s      Sender
	tags   []string
	prefix string
}

type key int

const xstatsKey key = 0

// NewContext returns a copy of the parent context and associates it with passed stats.
func NewContext(ctx context.Context, xs XStater) context.Context {
	return context.WithValue(ctx, xstatsKey, xs)
}

// FromContext retreives the request's xstats client from a given context if any.
// If no xstats is embeded in the context, a nop instance is returned so you can
// use it safely without having to test for it's presence.
func FromContext(ctx context.Context) XStater {
	rc, ok := ctx.Value(xstatsKey).(XStater)
	if ok {
		return rc
	}
	return nop
}

// FromRequest gets the xstats client in the request's context.
// This is a shortcut for xstats.FromContext(r.Context())
func FromRequest(r *http.Request) XStater {
	if r == nil {
		return nop
	}
	return FromContext(r.Context())
}

// NewHandler creates a new handler with the provided metric client.
// If some tags are provided, the will be added to all logged metrics.
func NewHandler(s Sender, tags []string) func(http.Handler) http.Handler {
	return NewHandlerPrefix(s, tags, "")
}

// NewHandlerPrefix creates a new handler with the provided metric client.
// If some tags are provided, the will be added to all logged metrics.
// If the prefix argument is provided, all produced metrics will have this
// prefix prepended.
func NewHandlerPrefix(s Sender, tags []string, prefix string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			xs := NewPrefix(s, prefix).(*xstats)
			xs.AddTags(tags...)
			ctx := NewContext(r.Context(), xs)
			next.ServeHTTP(w, r.WithContext(ctx))
			xs.Close()
		})
	}
}
