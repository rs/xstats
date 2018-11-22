package emptysender

import (
	"testing"

	"github.com/rs/xstats"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := New()
	assert.Implements(t, (*xstats.Sender)(nil), s)
}
