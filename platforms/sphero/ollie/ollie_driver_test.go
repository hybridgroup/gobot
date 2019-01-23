package ollie

import (
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/platforms/sphero"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestOllieDriver() *Driver {
	d := NewDriver(NewBleTestAdaptor())
	return d
}

func TestOllieDriver(t *testing.T) {
	d := initTestOllieDriver()
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestOllieDriverStartAndHalt(t *testing.T) {
	d := initTestOllieDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
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
		//0x0B is the locator ID
		packet := []byte{0xFF, 0xFF, 0x00, 0x00, 0x0B, point.x1, point.x2, point.y1, point.y2, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

		d.GetLocatorData(func(p Point2D) {
			gobottest.Assert(t, p.Y, point.y)
		})
		d.HandleResponses(packet, nil)
	}
}

func TestDataStreaming(t *testing.T) {
	d := initTestOllieDriver()

	d.SetDataStreamingConfig(sphero.DefaultDataStreamingConfig())

	response := false
	d.On("sensordata", func(data interface{}) {
		cont := data.(DataStreamingPacket)
		fmt.Printf("got streaming packet: %+v \n", cont)
		gobottest.Assert(t, cont.RawAccX, int16(10))
		response = true
	})

	//example data packet
	p1 := []string{"FFFE030053000A003900FAFFFE0007FFFF000000",
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

	//send empty packet to indicate start of next message
	d.HandleResponses([]byte{0xFF}, nil)
	time.Sleep(10 * time.Millisecond)
	if response == false {
		t.Error("no response recieved")
	}
}

func parseBytes(s string) (f byte) {
	i, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return
	}

	f = byte(math.Float32frombits(uint32(i)))

	return
}
