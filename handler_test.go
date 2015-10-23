package xmetrics

import (
	"net/http"
	"testing"

	"github.com/rs/xhandler"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestHandler(t *testing.T) {
	c := &fakeClient{}
	n := xhandler.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		s := FromContext(ctx)
		rc, ok := s.(*requestClient)
		assert.True(t, ok)
		assert.Equal(t, c, rc.c)
		assert.Equal(t, []string{"envtag"}, rc.tags)
	})
	h := NewHandler(c, n, []string{"envtag"})
	h.ServeHTTP(context.Background(), nil, nil)
}
