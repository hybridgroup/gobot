package jetson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

var _ gobot.PWMPinner = (*PWMPin)(nil)

func TestPwmPin(t *testing.T) {
	a := system.NewAccesser()
	const (
		exportPath    = "/sys/class/pwm/pwmchip0/export"
		unexportPath  = "/sys/class/pwm/pwmchip0/unexport"
		enablePath    = "/sys/class/pwm/pwmchip0/pwm3/enable"
		periodPath    = "/sys/class/pwm/pwmchip0/pwm3/period"
		dutyCyclePath = "/sys/class/pwm/pwmchip0/pwm3/duty_cycle"
	)
	mockPaths := []string{
		exportPath,
		unexportPath,
		enablePath,
		periodPath,
		dutyCyclePath,
	}
	fs := a.UseMockFilesystem(mockPaths)

	pin := NewPWMPin(a, "/sys/class/pwm/pwmchip0", "3")
	require.Equal(t, "", fs.Files[exportPath].Contents)
	require.Equal(t, "", fs.Files[unexportPath].Contents)
	require.Equal(t, "", fs.Files[enablePath].Contents)
	require.Equal(t, "", fs.Files[periodPath].Contents)
	require.Equal(t, "", fs.Files[dutyCyclePath].Contents)

	assert.Nil(t, pin.Export())
	assert.Equal(t, "3", fs.Files[exportPath].Contents)

	assert.Nil(t, pin.SetEnabled(true))
	assert.Equal(t, "1", fs.Files[enablePath].Contents)

	val, _ := pin.Polarity()
	assert.True(t, val)
	assert.Nil(t, pin.SetPolarity(false))
	val, _ = pin.Polarity()
	assert.True(t, val)

	_, err := pin.Period()
	assert.Error(t, err, "Jetson PWM pin period not set")
	assert.Error(t, pin.SetDutyCycle(10000), "Jetson PWM pin period not set")
	assert.Equal(t, "", fs.Files[dutyCyclePath].Contents)

	assert.Nil(t, pin.SetPeriod(20000000))
	// TODO: see PR #990 assert.Equal(t, "20000000", fs.Files[periodPath].Contents)
	period, _ := pin.Period()
	assert.Equal(t, uint32(20000000), period)
	assert.Error(t, pin.SetPeriod(10000000), "Cannot set the period of individual PWM pins on Jetson")
	// TODO: see PR #990 assert.Equal(t, "20000000", fs.Files[periodPath].Contents)

	dc, _ := pin.DutyCycle()
	assert.Equal(t, uint32(0), dc)

	assert.Nil(t, pin.SetDutyCycle(10000))
	assert.Equal(t, "10000", fs.Files[dutyCyclePath].Contents)
	dc, _ = pin.DutyCycle()
	assert.Equal(t, uint32(10000), dc)

	assert.Error(t, pin.SetDutyCycle(999999999), "Duty cycle exceeds period")
	dc, _ = pin.DutyCycle()
	assert.Equal(t, "10000", fs.Files[dutyCyclePath].Contents)
	assert.Equal(t, uint32(10000), dc)

	assert.Nil(t, pin.Unexport())
	assert.Equal(t, "3", fs.Files[unexportPath].Contents)
}
