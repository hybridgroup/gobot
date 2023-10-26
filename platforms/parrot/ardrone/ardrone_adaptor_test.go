package ardrone

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestArdroneAdaptor() *Adaptor {
	a := NewAdaptor()
	a.connect = func(a *Adaptor) (drone, error) {
		return &testDrone{}, nil
	}
	return a
}

func TestArdroneAdaptor(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, "192.168.1.1", a.config.Ip)

	a = NewAdaptor("192.168.100.100")
	assert.Equal(t, "192.168.100.100", a.config.Ip)
}

func TestArdroneAdaptorConnect(t *testing.T) {
	a := initTestArdroneAdaptor()
	assert.NoError(t, a.Connect())

	a.connect = func(a *Adaptor) (drone, error) {
		return nil, errors.New("connection error")
	}
	assert.ErrorContains(t, a.Connect(), "connection error")
}

func TestArdroneAdaptorName(t *testing.T) {
	a := initTestArdroneAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "ARDrone"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestArdroneAdaptorFinalize(t *testing.T) {
	a := initTestArdroneAdaptor()
	assert.NoError(t, a.Finalize())
}
