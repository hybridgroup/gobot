package ollie

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/sphero"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestOllieDriver() *Driver {
	d := NewDriver(NewBleTestAdaptor())
	return d
}

func TestOllieDriver(t *testing.T) {
	d := initTestOllieDriver()
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestOllieDriverStartAndHalt(t *testing.T) {
	d := initTestOllieDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
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
		packet := []byte{0xFF, 0xFF, 0x00, 0x00, 0x0B, point.x1, point.x2, point.y1, point.y2, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

		d.GetLocatorData(func(p Point2D) {
			assert.Equal(t, point.y, p.Y)
		})
		d.HandleResponses(packet, nil)
	}
}

func TestDataStreaming(t *testing.T) {
	d := initTestOllieDriver()

	_ = d.SetDataStreamingConfig(sphero.DefaultDataStreamingConfig())

	response := false
	_ = d.On("sensordata", func(data interface{}) {
		cont := data.(DataStreamingPacket)
		fmt.Printf("got streaming packet: %+v \n", cont)
		assert.Equal(t, int16(10), cont.RawAccX)
		response = true
	})

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
			b, _ := strconv.ParseUint(string(a), 16, 16)
			c := uint16(b)
			bytes = append(bytes, byte(c))
		}
		d.HandleResponses(bytes, nil)

	}

	// send empty packet to indicate start of next message
	d.HandleResponses([]byte{0xFF}, nil)
	time.Sleep(10 * time.Millisecond)
	if response == false {
		t.Error("no response received")
	}
}
