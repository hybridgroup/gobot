package pebble

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, a.Connect())
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	assert.NoError(t, a.Finalize())
}
