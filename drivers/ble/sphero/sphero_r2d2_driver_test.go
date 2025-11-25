package sphero

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*R2D2Driver)(nil)

func initTestR2D2Driver() *R2D2Driver {
	d := NewR2D2Driver(testutil.NewBleTestAdaptor())
	return d
}

func TestNewR2D2Driver(t *testing.T) {
	d := NewR2D2Driver(testutil.NewBleTestAdaptor())
	assert.NotNil(t, d.Driver)
	assert.NotNil(t, d.Eventer)
	assert.Equal(t, d.defaultCollisionConfig, r2d2DefaultCollisionConfig())
	assert.NotNil(t, d.packetChannel)
}

func TestNewR2D2DriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewR2D2Driver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestR2D2StartAndHalt(t *testing.T) {
	d := initTestR2D2Driver()
	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestR2D2HandleResponsesOutOfOrder(t *testing.T) {
	d := initTestR2D2Driver()

	// example data packet that's out of order where end of packet (D8) is not at the end
	p1 := []string{
		"8D", "09", "13", "0D", "00", "D8", "D6", "00",
	}

	for _, elem := range p1 {
		var bytes []byte
		for i := 0; i < len([]rune(elem)); i += 2 {
			a := []rune(elem)[i : i+2]
			b, err := strconv.ParseUint(string(a), 16, 16)
			require.NoError(t, err)

			c := uint16(b)
			bytes = append(bytes, byte(c))
		}
		d.handleResponses(bytes)
	}

	// assert buffer contains packets except end of packet
	assert.Equal(t, []byte{0x8D, 0x09, 0x13, 0x0D, 0x00, 0xD6, 0x00}, d.asyncBuffer)
	// send start packet to indicate start of next message
	d.handleResponses([]byte{sop})
	assert.Equal(t, []byte{sop}, d.asyncBuffer)
	assert.Equal(t, []byte{0x8D, 0x09, 0x13, 0x0D, 0x00, 0xD6, 0x00, 0xD8}, d.asyncMessage)
}
