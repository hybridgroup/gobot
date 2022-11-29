package jetson

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
)

var _ gobot.PWMPinner = (*PWMPin)(nil)

func TestPwmPin(t *testing.T) {
	a := system.NewAccesser()
	mockPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm0/enable",
		"/sys/class/pwm/pwmchip0/pwm0/period",
		"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
	}
	a.UseMockFilesystem(mockPaths)

	pin := NewPWMPin(a, "/sys/class/pwm/pwmchip0", "0")
	gobottest.Assert(t, pin.Export(), nil)
	gobottest.Assert(t, pin.Enable(true), nil)
	val, _ := pin.Polarity()
	gobottest.Assert(t, val, "normal")
	gobottest.Assert(t, pin.InvertPolarity(true), nil)
	val, _ = pin.Polarity()
	gobottest.Assert(t, val, "normal")

	period, err := pin.Period()
	gobottest.Assert(t, err, errors.New("Jetson PWM pin period not set"))
	gobottest.Assert(t, pin.SetDutyCycle(10000), errors.New("Jetson PWM pin period not set"))

	gobottest.Assert(t, pin.SetPeriod(20000000), nil)
	period, _ = pin.Period()
	gobottest.Assert(t, period, uint32(20000000))
	gobottest.Assert(t, pin.SetPeriod(10000000), errors.New("Cannot set the period of individual PWM pins on Jetson"))

	dc, _ := pin.DutyCycle()
	gobottest.Assert(t, dc, uint32(0))

	gobottest.Assert(t, pin.SetDutyCycle(10000), nil)
	dc, _ = pin.DutyCycle()
	gobottest.Assert(t, dc, uint32(10000))

	gobottest.Assert(t, pin.SetDutyCycle(999999999), errors.New("Duty cycle exceeds period"))
	dc, _ = pin.DutyCycle()
	gobottest.Assert(t, dc, uint32(10000))

	gobottest.Assert(t, pin.Unexport(), nil)
}
