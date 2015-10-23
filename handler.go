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
	next xhandler.Handler
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

// FromContext retreives the request client from a given context if any.
func FromContext(ctx context.Context) XStater {
	rc, ok := ctx.Value(xstatsKey).(XStater)
	if ok {
		return rc
	}
	return nop
}

// NewHandler creates a new handler with the provided metric client.
// If some tags are provided, the will be added to all logged metrics.
func NewHandler(s Sender, tags []string, next xhandler.Handler) *Handler {
	return &Handler{
		s:    s,
		tags: tags,
		next: next,
	}
}

// Implements xhandler.Handler interface
func (h *Handler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	xs, _ := xstatsPool.Get().(*xstats)
	xs.s = h.s
	xs.tags = append([]string{}, h.tags...)
	ctx = newContext(ctx, xs)
	h.next.ServeHTTP(ctx, w, r)
	xs.s = nil
	xs.tags = nil
	xstatsPool.Put(xs)
}
