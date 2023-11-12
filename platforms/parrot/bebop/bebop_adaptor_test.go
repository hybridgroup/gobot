package bebop

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestBebopAdaptor() *Adaptor {
	a := NewAdaptor()
	a.connect = func(b *Adaptor) error {
		b.drone = &testDrone{}
		return nil
	}
	return a
}

func TestBebopAdaptorName(t *testing.T) {
	a := NewAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Bebop"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestBebopAdaptorConnect(t *testing.T) {
	a := initTestBebopAdaptor()
	require.NoError(t, a.Connect())

	a.connect = func(a *Adaptor) error {
		return errors.New("connection error")
	}
	require.ErrorContains(t, a.Connect(), "connection error")
}

func TestBebopAdaptorFinalize(t *testing.T) {
	a := initTestBebopAdaptor()
	_ = a.Connect()
	require.NoError(t, a.Finalize())
}
