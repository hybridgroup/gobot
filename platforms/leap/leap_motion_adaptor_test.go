package leap

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestLeapMotionAdaptor() *Adaptor {
	a := NewAdaptor("")
	a.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil } //nolint:nilnil // ok for test
	return a
}

func TestLeapMotionAdaptor(t *testing.T) {
	a := NewAdaptor("127.0.0.1")
	assert.Equal(t, "127.0.0.1", a.Port())
}

func TestLeapMotionAdaptorName(t *testing.T) {
	a := NewAdaptor("127.0.0.1")
	assert.True(t, strings.HasPrefix(a.Name(), "Leap"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestLeapMotionAdaptorConnect(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	require.NoError(t, a.Connect())

	a.connect = func(port string) (io.ReadWriteCloser, error) {
		return nil, errors.New("connection error")
	}
	require.ErrorContains(t, a.Connect(), "connection error")
}

func TestLeapMotionAdaptorFinalize(t *testing.T) {
	a := initTestLeapMotionAdaptor()
	require.NoError(t, a.Finalize())
}
