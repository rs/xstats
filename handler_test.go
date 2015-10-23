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
	n := xhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		xs, ok := FromContext(ctx).(*xstats)
		assert.True(t, ok)
		assert.Equal(t, s, xs.s)
		assert.Equal(t, []string{"envtag"}, xs.tags)
	})
	h := NewHandler(s, []string{"envtag"}, n)
	h.ServeHTTP(context.Background(), nil, nil)
}
