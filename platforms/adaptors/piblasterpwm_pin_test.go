package adaptors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

var _ gobot.PWMPinner = (*piBlasterPWMPin)(nil)

func TestPiBlasterPWMPin(t *testing.T) {
	// arrange
	const path = "/dev/pi-blaster"
	a := system.NewAccesser()
	a.UseMockFilesystem([]string{path})
	pin := newPiBlasterPWMPin(a, 1)
	// act & assert: activate pin for usage
	require.NoError(t, pin.Export())
	require.NoError(t, pin.SetEnabled(true))
	// act & assert: get and set polarity
	val, err := pin.Polarity()
	require.NoError(t, err)
	assert.True(t, val)
	require.NoError(t, pin.SetPolarity(false))
	polarity, err := pin.Polarity()
	assert.True(t, polarity)
	require.NoError(t, err)
	// act & assert: get and set period
	period, err := pin.Period()
	require.NoError(t, err)
	assert.Equal(t, uint32(0), period)
	require.NoError(t, pin.SetPeriod(20000000))
	period, err = pin.Period()
	require.NoError(t, err)
	assert.Equal(t, uint32(20000000), period)
	err = pin.SetPeriod(10000000)
	require.EqualError(t, err, "the period of PWM pins needs to be set to '10000000' in pi-blaster source code")
	// act & assert: cleanup
	require.NoError(t, pin.Unexport())
}

func TestPiBlasterPWMPin_DutyCycle(t *testing.T) {
	// arrange
	const path = "/dev/pi-blaster"
	a := system.NewAccesser()
	a.UseMockFilesystem([]string{path})
	pin := newPiBlasterPWMPin(a, 1)
	// act & assert: activate pin for usage
	require.NoError(t, pin.Export())
	require.NoError(t, pin.SetEnabled(true))
	// act & assert zero
	dc, err := pin.DutyCycle()
	require.NoError(t, err)
	assert.Equal(t, uint32(0), dc)
	// act & assert error without period set, the value remains zero
	err = pin.SetDutyCycle(10000)
	require.EqualError(t, err, "pi-blaster PWM pin period not set while try to set duty cycle to '10000'")
	dc, err = pin.DutyCycle()
	require.NoError(t, err)
	assert.Equal(t, uint32(0), dc)
	// arrange, act & assert a value
	pin.period = 20000000
	require.NoError(t, pin.SetDutyCycle(10000))
	dc, err = pin.DutyCycle()
	require.NoError(t, err)
	assert.Equal(t, uint32(10000), dc)
	// act & assert error on over limit, the value remains
	err = pin.SetDutyCycle(20000001)
	require.EqualError(t, err, "the duty cycle (20000001) exceeds period (20000000) for pi-blaster")
	dc, err = pin.DutyCycle()
	require.NoError(t, err)
	assert.Equal(t, uint32(10000), dc)
	// act & assert: cleanup
	require.NoError(t, pin.Unexport())
}
