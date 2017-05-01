package raspi

import (
	"errors"
	"testing"

	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

var _ sysfs.PWMPinner = (*PWMPin)(nil)

func TestPwmPin(t *testing.T) {
	pin := NewPWMPin("1")
	gobottest.Assert(t, pin.Export(), nil)
	gobottest.Assert(t, pin.Enable(true), nil)
	val, _ := pin.Polarity()
	gobottest.Assert(t, val, "normal")
	gobottest.Assert(t, pin.InvertPolarity(true), nil)
	val, _ = pin.Polarity()
	gobottest.Assert(t, val, "normal")

	period, _ := pin.Period()
	gobottest.Assert(t, period, uint32(10000000))
	gobottest.Assert(t, pin.SetPeriod(1000), nil)
	period, _ = pin.Period()
	gobottest.Assert(t, period, uint32(10000000))

	dc, _ := pin.DutyCycle()
	gobottest.Assert(t, dc, uint32(0))

	// call currently fails in test
	gobottest.Refute(t, pin.SetDutyCycle(10000), nil)
	dc, _ = pin.DutyCycle()
	gobottest.Assert(t, dc, uint32(10000))

	gobottest.Assert(t, pin.SetDutyCycle(999999999), errors.New("Duty cycle exceeds period."))
	dc, _ = pin.DutyCycle()
	gobottest.Assert(t, dc, uint32(10000))

	// call currently fails in test
	gobottest.Refute(t, pin.Unexport(), nil)
}
