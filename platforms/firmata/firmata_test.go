package firmata

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

type NullReadWriteCloser struct{}

func (NullReadWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (NullReadWriteCloser) Read(b []byte) (int, error) {
	return len(b), nil
}

var closeErr error

func (NullReadWriteCloser) Close() error {
	return closeErr
}

func initTestFirmata() *board {
	b := newBoard(NullReadWriteCloser{})
	b.initTimeInterval = 0 * time.Second
	// arduino uno r3 firmware response "StandardFirmata.ino"
	b.process([]byte{240, 121, 2, 3, 83, 0, 116, 0, 97, 0, 110, 0, 100, 0, 97,
		0, 114, 0, 100, 0, 70, 0, 105, 0, 114, 0, 109, 0, 97, 0, 116, 0, 97, 0, 46,
		0, 105, 0, 110, 0, 111, 0, 247})
	// arduino uno r3 capabilities response
	b.process([]byte{240, 108, 127, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 3,
		8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1,
		1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0,
		1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 3, 8,
		4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 2,
		10, 127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10,
		127, 0, 1, 1, 1, 2, 10, 6, 1, 127, 0, 1, 1, 1, 2, 10, 6, 1, 127, 247})
	// arduino uno r3 analog mapping response
	b.process([]byte{240, 106, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127,
		127, 127, 127, 127, 0, 1, 2, 3, 4, 5, 247})
	return b
}

func TestReportVersion(t *testing.T) {
	b := initTestFirmata()
	//test if functions executes
	b.queryReportVersion()
}

func TestQueryFirmware(t *testing.T) {
	b := initTestFirmata()
	//test if functions executes
	b.queryFirmware()
}

func TestQueryPinState(t *testing.T) {
	b := initTestFirmata()
	//test if functions executes
	b.queryPinState(byte(1))
}

func TestProcess(t *testing.T) {
	b := initTestFirmata()
	sem := make(chan bool)
	//reportVersion
	gobot.Once(b.events["report_version"], func(data interface{}) {
		gobot.Assert(t, data.(string), "1.17")
		sem <- true
	})
	b.process([]byte{0xF9, 0x01, 0x11})
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("report_version was not published")
	}
	//analogMessageRangeStart
	gobot.Once(b.events["analog_read_0"], func(data interface{}) {
		b := data.([]byte)
		gobot.Assert(t,
			int(uint(b[0])<<24|uint(b[1])<<16|uint(b[2])<<8|uint(b[3])),
			675)
		sem <- true
	})
	b.process([]byte{0xE0, 0x23, 0x05})
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("analog_read_0 was not published")
	}
	gobot.Once(b.events["analog_read_1"], func(data interface{}) {
		b := data.([]byte)
		gobot.Assert(t,
			int(uint(b[0])<<24|uint(b[1])<<16|uint(b[2])<<8|uint(b[3])),
			803)
		sem <- true
	})
	b.process([]byte{0xE1, 0x23, 0x06})
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("analog_read_1 was not published")
	}
	//digitalMessageRangeStart
	b.pins[2].mode = input
	gobot.Once(b.events["digital_read_2"], func(data interface{}) {
		gobot.Assert(t, int(data.([]byte)[0]), 1)
		sem <- true
	})
	b.process([]byte{0x90, 0x04, 0x00})
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("digital_read_2 was not published")
	}
	gobot.Once(b.events["analog_read_1"], func(data interface{}) {
		b := data.([]byte)
		gobot.Assert(t,
			int(uint(b[0])<<24|uint(b[1])<<16|uint(b[2])<<8|uint(b[3])),
			803)
		sem <- true
	})
	b.pins[4].mode = input
	gobot.Once(b.events["digital_read_4"], func(data interface{}) {
		gobot.Assert(t, int(data.([]byte)[0]), 1)
		sem <- true
	})
	b.process([]byte{0x90, 0x16, 0x00})
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("digital_read_4 was not published")
	}
	//pinStateResponse
	gobot.Once(b.events["pin_13_state"], func(data interface{}) {
		gobot.Assert(t, data, map[string]int{
			"pin":   13,
			"mode":  1,
			"value": 1,
		})
		sem <- true
	})
	b.process([]byte{240, 110, 13, 1, 1, 247})
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("pin_13_state was not published")
	}
	//i2cReply
	gobot.Once(b.events["i2c_reply"], func(data interface{}) {
		i2cReply := map[string][]byte{
			"slave_address": []byte{9},
			"register":      []byte{0},
			"data":          []byte{152, 1, 154}}
		gobot.Assert(t, data.(map[string][]byte), i2cReply)
		sem <- true
	})
	b.process([]byte{240, 119, 9, 0, 0, 0, 24, 1, 1, 0, 26, 1, 247})
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("i2c_reply was not published")
	}
	//firmwareName
	gobot.Once(b.events["firmware_query"], func(data interface{}) {
		gobot.Assert(t, data.(string), "StandardFirmata.ino")
		sem <- true
	})
	b.process([]byte{240, 121, 2, 3, 83, 0, 116, 0, 97, 0, 110, 0, 100, 0, 97,
		0, 114, 0, 100, 0, 70, 0, 105, 0, 114, 0, 109, 0, 97, 0, 116, 0, 97, 0, 46,
		0, 105, 0, 110, 0, 111, 0, 247})
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("firmware_query was not published")
	}
	//stringData
	gobot.Once(b.events["string_data"], func(data interface{}) {
		gobot.Assert(t, data.(string), "Hello Firmata!")
		sem <- true
	})
	b.process(append([]byte{240, 0x71}, []byte("Hello Firmata!")...))
	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("string_data was not published")
	}
}
