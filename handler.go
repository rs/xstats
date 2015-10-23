package xmetrics

import (
	"net/http"
	"sync"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// Handler injects a per request metrics client in the net/context which can be
// retrived using xmetrics.FromContext(ctx)
type Handler struct {
	c    Client
	tags []string
	next xhandler.Handler
}

type key int

const requestClientKey key = 0

var requestClientPool = sync.Pool{
	New: func() interface{} {
		return &requestClient{}
	},
}

// newContext returns a context with the given stats request client stored as value.
func newContext(ctx context.Context, rc RequestClient) context.Context {
	return context.WithValue(ctx, requestClientKey, rc)
}

// FromContext retreives the request client from a given context if any.
func FromContext(ctx context.Context) RequestClient {
	rc, ok := ctx.Value(requestClientKey).(RequestClient)
	if ok {
		return rc
	}
	return nopClient
}

// NewHandler creates a new handler with the provided metric client.
// If some tags are provided, the will be added to all logged metrics.
func NewHandler(c Client, tags []string, next xhandler.Handler) *Handler {
	return &Handler{
		c:    c,
		tags: tags,
		next: next,
	}
}

// Implements xhandler.Handler interface
func (h *Handler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	rc, _ := requestClientPool.Get().(*requestClient)
	rc.c = h.c
	rc.tags = append([]string{}, h.tags...)
	ctx = newContext(ctx, rc)
	h.next.ServeHTTP(ctx, w, r)
	rc.tags = nil
	rc.c = nil
	requestClientPool.Put(rc)
}
