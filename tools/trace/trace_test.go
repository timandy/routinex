package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTraceId(t *testing.T) {
	const iterations = 100
	mp := map[string]bool{}
	for i := 0; i < 100; i++ {
		s := NewTraceId()
		assert.Len(t, s, 8)
		assert.False(t, mp[s])
		mp[s] = true
	}
	assert.Equal(t, iterations, len(mp))
}
