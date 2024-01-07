package pebble

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestAdaptor() *Adaptor {
	return NewAdaptor()
}

func TestAdaptor(t *testing.T) {
	a := initTestAdaptor()
	assert.Equal(t, "Pebble", a.Name())
}

func TestAdaptorConnect(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.Connect())
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.Finalize())
}
