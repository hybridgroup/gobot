package raspi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

var _ gobot.PWMPinner = (*PWMPin)(nil)

func TestPwmPin(t *testing.T) {
	const path = "/dev/pi-blaster"
	a := system.NewAccesser()
	a.UseMockFilesystem([]string{path})

	pin := NewPWMPin(a, path, "1")

	require.NoError(t, pin.Export())
	require.NoError(t, pin.SetEnabled(true))

	val, _ := pin.Polarity()
	assert.True(t, val)

	require.NoError(t, pin.SetPolarity(false))

	val, _ = pin.Polarity()
	assert.True(t, val)

	_, err := pin.Period()
	require.ErrorContains(t, err, "Raspi PWM pin period not set")
	require.ErrorContains(t, pin.SetDutyCycle(10000), "Raspi PWM pin period not set")

	require.NoError(t, pin.SetPeriod(20000000))
	period, _ := pin.Period()
	assert.Equal(t, uint32(20000000), period)
	require.ErrorContains(t, pin.SetPeriod(10000000), "Cannot set the period of individual PWM pins on Raspi")

	dc, _ := pin.DutyCycle()
	assert.Equal(t, uint32(0), dc)

	require.NoError(t, pin.SetDutyCycle(10000))

	dc, _ = pin.DutyCycle()
	assert.Equal(t, uint32(10000), dc)

	require.ErrorContains(t, pin.SetDutyCycle(999999999), "Duty cycle exceeds period")
	dc, _ = pin.DutyCycle()
	assert.Equal(t, uint32(10000), dc)

	require.NoError(t, pin.Unexport())
}
