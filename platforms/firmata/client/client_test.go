package client

import (
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

type readWriteCloser struct{}

func (readWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

var testReadData = []byte{}

func (readWriteCloser) Read(b []byte) (int, error) {
	size := len(b)
	if len(testReadData) < size {
		size = len(testReadData)
	}
	copy(b, []byte(testReadData)[:size])
	testReadData = testReadData[size:]

	return size, nil
}

func (readWriteCloser) Close() error {
	return nil
}

func testProtocolResponse() []byte {
	// arduino uno r3 protocol response "2.3"
	return []byte{249, 2, 3}
}

func testFirmwareResponse() []byte {
	// arduino uno r3 firmware response "StandardFirmata.ino"
	return []byte{240, 121, 2, 3, 83, 0, 116, 0, 97, 0, 110, 0, 100, 0, 97,
		0, 114, 0, 100, 0, 70, 0, 105, 0, 114, 0, 109, 0, 97, 0, 116, 0, 97, 0, 46,
		0, 105, 0, 110, 0, 111, 0, 247}
}

func testCapabilitiesResponse() []byte {
	// arduino uno r3 capabilities response
	return []byte{240, 108, 127, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 3,
		8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1,
		1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0,
		1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 3, 8,
		4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 2,
		10, 127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10,
		127, 0, 1, 1, 1, 2, 10, 6, 1, 127, 0, 1, 1, 1, 2, 10, 6, 1, 127, 247}
}

func testAnalogMappingResponse() []byte {
	// arduino uno r3 analog mapping response
	return []byte{240, 106, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127,
		127, 127, 127, 127, 0, 1, 2, 3, 4, 5, 247}
}

func initTestFirmata() *Client {
	b := New()
	b.connection = readWriteCloser{}

	for _, f := range []func() []byte{
		testProtocolResponse,
		testFirmwareResponse,
		testCapabilitiesResponse,
		testAnalogMappingResponse,
	} {
		testReadData = f()
		b.process()
	}

	b.connected = true
	b.Connect(readWriteCloser{})

	return b
}

func TestReportVersion(t *testing.T) {
	b := initTestFirmata()
	//test if functions executes
	b.ProtocolVersionQuery()
}

func TestQueryFirmware(t *testing.T) {
	b := initTestFirmata()
	//test if functions executes
	b.FirmwareQuery()
}

func TestQueryPinState(t *testing.T) {
	b := initTestFirmata()
	//test if functions executes
	b.PinStateQuery(1)
}

func TestProcess(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()

	tests := []struct {
		event    string
		data     []byte
		expected interface{}
		init     func()
	}{
		{
			event:    "ProtocolVersion",
			data:     []byte{249, 2, 3},
			expected: "2.3",
			init:     func() {},
		},
		{
			event:    "AnalogRead0",
			data:     []byte{0xE0, 0x23, 0x05},
			expected: 675,
			init:     func() {},
		},
		{
			event:    "AnalogRead1",
			data:     []byte{0xE1, 0x23, 0x06},
			expected: 803,
			init:     func() {},
		},
		{
			event:    "DigitalRead2",
			data:     []byte{0x90, 0x04, 0x00},
			expected: 1,
			init:     func() { b.pins[2].Mode = Input },
		},
		{
			event:    "DigitalRead4",
			data:     []byte{0x90, 0x16, 0x00},
			expected: 1,
			init:     func() { b.pins[4].Mode = Input },
		},
		{
			event:    "PinState13",
			data:     []byte{240, 110, 13, 1, 1, 247},
			expected: Pin{[]int{0, 1, 4}, 1, 0, 1, 127},
			init:     func() {},
		},
		{
			event: "I2cReply",
			data:  []byte{240, 119, 9, 0, 0, 0, 24, 1, 1, 0, 26, 1, 247},
			expected: I2cReply{
				Address:  9,
				Register: 0,
				Data:     []byte{152, 1, 154},
			},
			init: func() {},
		},
		{
			event: "FirmwareQuery",
			data: []byte{240, 121, 2, 3, 83, 0, 116, 0, 97, 0, 110, 0, 100, 0, 97,
				0, 114, 0, 100, 0, 70, 0, 105, 0, 114, 0, 109, 0, 97, 0, 116, 0, 97, 0, 46,
				0, 105, 0, 110, 0, 111, 0, 247},
			expected: "StandardFirmata.ino",
			init:     func() {},
		},
		{
			event:    "StringData",
			data:     append([]byte{240, 0x71}, append([]byte("Hello Firmata!"), 247)...),
			expected: "Hello Firmata!",
			init:     func() {},
		},
	}

	for _, test := range tests {
		test.init()
		gobot.Once(b.Event(test.event), func(data interface{}) {
			gobot.Assert(t, data, test.expected)
			sem <- true
		})

		testReadData = test.data
		go b.process()

		select {
		case <-sem:
		case <-time.After(10 * time.Millisecond):
			t.Errorf("%v was not published", test.event)
		}
	}
}

func TestConnect(t *testing.T) {
	b := New()

	response := testProtocolResponse()

	go func() {
		for {
			testReadData = append(testReadData, response...)
			<-time.After(100 * time.Millisecond)
		}
	}()

	gobot.Once(b.Event("ProtocolVersion"), func(data interface{}) {
		response = testFirmwareResponse()
	})

	gobot.Once(b.Event("FirmwareQuery"), func(data interface{}) {
		response = testCapabilitiesResponse()
	})

	gobot.Once(b.Event("CapabilityQuery"), func(data interface{}) {
		response = testAnalogMappingResponse()
	})

	gobot.Once(b.Event("AnalogMappingQuery"), func(data interface{}) {
		response = testProtocolResponse()
	})

	gobot.Assert(t, b.Connect(readWriteCloser{}), nil)
}
