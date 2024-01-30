package bit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	const (
		oldVal = 1
		bitPos = 7
		want   = 129
	)
	require.False(t, IsSet(oldVal, bitPos))

	got := Set(1, 7)
	assert.Equal(t, want, got)
	assert.True(t, IsSet(got, bitPos))
}

func TestClear(t *testing.T) {
	const (
		oldVal     = 128
		bitPos     = 7
		want   int = 0
	)
	require.True(t, IsSet(oldVal, bitPos))

	got := Clear(128, 7)
	assert.Equal(t, want, got)
	assert.False(t, IsSet(got, bitPos))
}
