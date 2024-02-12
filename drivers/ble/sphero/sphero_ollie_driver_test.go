//nolint:forcetypeassert // ok here
package sphero

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
	"gobot.io/x/gobot/v2/drivers/common/spherocommon"
)

var _ gobot.Driver = (*OllieDriver)(nil)

func initTestOllieDriver() *OllieDriver {
	d := NewOllieDriver(testutil.NewBleTestAdaptor())
	return d
}

func TestNewOllieDriver(t *testing.T) {
	d := NewOllieDriver(testutil.NewBleTestAdaptor())
	assert.NotNil(t, d.Driver)
	assert.NotNil(t, d.Eventer)
	assert.Equal(t, d.defaultCollisionConfig, ollieDefaultCollisionConfig())
	assert.NotNil(t, d.packetChannel)
}

func TestNewOllieDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewOllieDriver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestOllieStartAndHalt(t *testing.T) {
	d := initTestOllieDriver()
	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestLocatorData(t *testing.T) {
	d := initTestOllieDriver()

	tables := []struct {
		x1 byte
		x2 byte
		y1 byte
		y2 byte
		x  int16
		y  int16
	}{
		{0x00, 0x05, 0x00, 0x05, 5, 5},
		{0x00, 0x00, 0x00, 0x00, 0, 0},
		{0x00, 0x0A, 0x00, 0xF0, 10, 240},
		{0x01, 0x00, 0x01, 0x00, 256, 256},
		{0xFF, 0xFE, 0xFF, 0xFE, -1, -1},
	}

	for _, point := range tables {
		// 0x0B is the locator ID
		packet := []byte{
			0xFF, 0xFF, 0x00, 0x00, 0x0B, point.x1, point.x2, point.y1, point.y2, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}

		d.GetLocatorData(func(p Point2D) {
			assert.Equal(t, point.y, p.Y)
		})
		d.handleResponses(packet)
	}
}

func TestDataStreaming(t *testing.T) {
	d := initTestOllieDriver()

	err := d.SetDataStreamingConfig(spherocommon.DefaultDataStreamingConfig())
	require.NoError(t, err)

	responseChan := make(chan bool)
	err = d.On("sensordata", func(data interface{}) {
		cont := data.(spherocommon.DataStreamingPacket)
		// fmt.Printf("got streaming packet: %+v \n", cont)
		assert.Equal(t, int16(10), cont.RawAccX)
		responseChan <- true
	})
	require.NoError(t, err)

	// example data packet
	p1 := []string{
		"FFFE030053000A003900FAFFFE0007FFFF000000",
		"000000000000000000FFECFFFB00010000004B01",
		"BD1034FFFF000300000000000000000000000000",
		"0000002701FDE500560000000000000065000000",
		"0000000000000071",
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

	// send empty packet to indicate start of next message
	d.handleResponses([]byte{0xFF})
	select {
	case <-responseChan:
	case <-time.After(10 * time.Millisecond):
		t.Error("no response received")
	}
}
