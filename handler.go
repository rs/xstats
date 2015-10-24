package xstats

import (
	"net/http"
	"sync"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// Handler injects a per request metrics client in the net/context which can be
// retrived using xstats.FromContext(ctx)
type Handler struct {
	s    Sender
	tags []string
	next xhandler.HandlerC
}

type key int

const xstatsKey key = 0

var xstatsPool = sync.Pool{
	New: func() interface{} {
		return &xstats{}
	},
}

// newContext returns a context with the given stats request client stored as value.
func newContext(ctx context.Context, xs *xstats) context.Context {
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

// NewHandler creates a new handler with the provided metric client.
// If some tags are provided, the will be added to all logged metrics.
func NewHandler(s Sender, tags []string, next xhandler.HandlerC) *Handler {
	return &Handler{
		s:    s,
		tags: tags,
		next: next,
	}
}

// ServeHTTPC implements xhandler.HandlerC interface
func (h *Handler) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	xs, _ := xstatsPool.Get().(*xstats)
	xs.s = h.s
	xs.tags = append([]string{}, h.tags...)
	ctx = newContext(ctx, xs)
	h.next.ServeHTTPC(ctx, w, r)
	xs.s = nil
	xs.tags = nil
	xstatsPool.Put(xs)
}
