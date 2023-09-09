package joystick

import (
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestAdaptor() *Adaptor {
	a := NewAdaptor(6)
	a.connect = func(j *Adaptor) (err error) {
		j.joystick = &testJoystick{}
		return nil
	}
	return a
}

func TestJoystickAdaptorName(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Joystick"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestAdaptorConnect(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.Connect(), nil)

	a = NewAdaptor(6)
	err := a.Connect()
	gobottest.Assert(t, strings.HasPrefix(err.Error(), "No joystick available"), true)
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	_ = a.Connect()
	gobottest.Assert(t, a.Finalize(), nil)
}
