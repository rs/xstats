// +build !go1.7

package xstats

import (
	"net/http"
	"testing"

	"github.com/rs/xhandler"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestHandler(t *testing.T) {
	s := &fakeSender{}
	n := xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		xs, ok := FromContext(ctx).(*xstats)
		assert.True(t, ok)
		assert.Equal(t, s, xs.s)
		assert.Equal(t, []string{"envtag"}, xs.tags)
	})
	h := NewHandler(s, []string{"envtag"})(n)
	h.ServeHTTPC(context.Background(), nil, nil)
}

func TestHandlerPrefix(t *testing.T) {
	s := &fakeSender{}
	n := xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		xs, ok := FromContext(ctx).(*xstats)
		assert.True(t, ok)
		assert.Equal(t, s, xs.s)
		assert.Equal(t, []string{"envtag"}, xs.tags)
		assert.Equal(t, "prefix.", xs.prefix)
	})
	h := NewHandlerPrefix(s, []string{"envtag"}, "prefix.")(n)
	h.ServeHTTPC(context.Background(), nil, nil)
}
