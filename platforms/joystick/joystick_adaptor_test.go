package joystick

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestAdaptor() *Adaptor {
	a := NewAdaptor()
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

	a = NewAdaptor()
	gobottest.Assert(t, a.Connect(), errors.New("No joystick available"))
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	a.Connect()
	gobottest.Assert(t, a.Finalize(), nil)
}
