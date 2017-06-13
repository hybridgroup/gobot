package client

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"gobot.io/x/gobot/gobottest"
)

type readWriteCloser struct{}

func (readWriteCloser) Write(p []byte) (int, error) {
	writeDataMutex.Lock()
	defer writeDataMutex.Unlock()
	return testWriteData.Write(p)
}

var clientMutex sync.Mutex
var writeDataMutex sync.Mutex
var readDataMutex sync.Mutex
var testReadData = []byte{}
var testWriteData = bytes.Buffer{}

func SetTestReadData(d []byte) {
	readDataMutex.Lock()
	defer readDataMutex.Unlock()
	testReadData = d
	return
}

func (readWriteCloser) Read(b []byte) (int, error) {
	readDataMutex.Lock()
	defer readDataMutex.Unlock()

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
		SetTestReadData(f())
		b.process()
	}

	return b
}

func TestPins(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	//test if functions executes
	gobottest.Assert(t, len(b.Pins()), 19)
}

func TestReportVersion(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	//test if functions executes
	gobottest.Assert(t, b.ProtocolVersionQuery(), nil)
}

func TestQueryFirmware(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	//test if functions executes
	gobottest.Assert(t, b.FirmwareQuery(), nil)
}

func TestQueryPinState(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	//test if functions executes
	gobottest.Assert(t, b.PinStateQuery(1), nil)
}

func TestProcessProtocolVersion(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	SetTestReadData([]byte{249, 2, 3})

	b.Once(b.Event("ProtocolVersion"), func(data interface{}) {
		gobottest.Assert(t, data, "2.3")
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("ProtocolVersion was not published")
	}
}

func TestProcessAnalogRead0(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	SetTestReadData([]byte{0xE0, 0x23, 0x05})

	b.Once(b.Event("AnalogRead0"), func(data interface{}) {
		gobottest.Assert(t, data, 675)
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("AnalogRead0 was not published")
	}
}

func TestProcessAnalogRead1(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	SetTestReadData([]byte{0xE1, 0x23, 0x06})

	b.Once(b.Event("AnalogRead1"), func(data interface{}) {
		gobottest.Assert(t, data, 803)
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("AnalogRead1 was not published")
	}
}

func TestProcessDigitalRead2(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	b.pins[2].Mode = Input
	SetTestReadData([]byte{0x90, 0x04, 0x00})

	b.Once(b.Event("DigitalRead2"), func(data interface{}) {
		gobottest.Assert(t, data, 1)
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("DigitalRead2 was not published")
	}
}

func TestProcessDigitalRead4(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	b.pins[4].Mode = Input
	SetTestReadData([]byte{0x90, 0x16, 0x00})

	b.Once(b.Event("DigitalRead4"), func(data interface{}) {
		gobottest.Assert(t, data, 1)
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("DigitalRead4 was not published")
	}
}

func TestDigitalWrite(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	gobottest.Assert(t, b.DigitalWrite(13, 0), nil)
}

func TestSetPinMode(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	gobottest.Assert(t, b.SetPinMode(13, Output), nil)
}

func TestAnalogWrite(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	gobottest.Assert(t, b.AnalogWrite(0, 128), nil)
}

func TestReportAnalog(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	gobottest.Assert(t, b.ReportAnalog(0, 1), nil)
	gobottest.Assert(t, b.ReportAnalog(0, 0), nil)
}

func TestProcessPinState13(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	SetTestReadData([]byte{240, 110, 13, 1, 1, 247})

	b.Once(b.Event("PinState13"), func(data interface{}) {
		gobottest.Assert(t, data, Pin{[]int{0, 1, 4}, 1, 0, 1, 127})
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("PinState13 was not published")
	}
}

func TestI2cConfig(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	gobottest.Assert(t, b.I2cConfig(100), nil)
}

func TestI2cWrite(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	gobottest.Assert(t, b.I2cWrite(0x00, []byte{0x01, 0x02}), nil)
}

func TestI2cRead(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	gobottest.Assert(t, b.I2cRead(0x00, 10), nil)
}

func TestWriteSysex(t *testing.T) {
	b := initTestFirmata()
	b.setConnected(true)
	gobottest.Assert(t, b.WriteSysex([]byte{0x01, 0x02}), nil)
}

func TestProcessI2cReply(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	SetTestReadData([]byte{240, 119, 9, 0, 0, 0, 24, 1, 1, 0, 26, 1, 247})

	b.Once(b.Event("I2cReply"), func(data interface{}) {
		gobottest.Assert(t, data, I2cReply{
			Address:  9,
			Register: 0,
			Data:     []byte{152, 1, 154},
		})
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("I2cReply was not published")
	}
}

func TestProcessFirmwareQuery(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	SetTestReadData([]byte{240, 121, 2, 3, 83, 0, 116, 0, 97, 0, 110, 0, 100, 0, 97,
		0, 114, 0, 100, 0, 70, 0, 105, 0, 114, 0, 109, 0, 97, 0, 116, 0, 97, 0, 46,
		0, 105, 0, 110, 0, 111, 0, 247})

	b.Once(b.Event("FirmwareQuery"), func(data interface{}) {
		gobottest.Assert(t, data, "StandardFirmata.ino")
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("FirmwareQuery was not published")
	}
}

func TestProcessStringData(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	SetTestReadData(append([]byte{240, 0x71}, append([]byte("Hello Firmata!"), 247)...))

	b.Once(b.Event("StringData"), func(data interface{}) {
		gobottest.Assert(t, data, "Hello Firmata!")
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("StringData was not published")
	}
}

func TestConnect(t *testing.T) {
	b := New()

	var responseMutex sync.Mutex
	responseMutex.Lock()
	response := testProtocolResponse()
	responseMutex.Unlock()

	b.Once(b.Event("ProtocolVersion"), func(data interface{}) {
		responseMutex.Lock()
		response = testFirmwareResponse()
		responseMutex.Unlock()
	})

	b.Once(b.Event("FirmwareQuery"), func(data interface{}) {
		responseMutex.Lock()
		response = testCapabilitiesResponse()
		responseMutex.Unlock()
	})

	b.Once(b.Event("CapabilityQuery"), func(data interface{}) {
		responseMutex.Lock()
		response = testAnalogMappingResponse()
		responseMutex.Unlock()
	})

	b.Once(b.Event("AnalogMappingQuery"), func(data interface{}) {
		responseMutex.Lock()
		response = testProtocolResponse()
		responseMutex.Unlock()
	})

	go func() {
		for {
			responseMutex.Lock()
			readDataMutex.Lock()
			testReadData = append(testReadData, response...)
			readDataMutex.Unlock()
			responseMutex.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	gobottest.Assert(t, b.Connect(readWriteCloser{}), nil)
	time.Sleep(150 * time.Millisecond)
	gobottest.Assert(t, b.Disconnect(), nil)
}

func TestServoConfig(t *testing.T) {
	b := New()
	b.connection = readWriteCloser{}

	tests := []struct {
		description string
		arguments   [3]int
		expected    []byte
		result      error
	}{
		{
			description: "Min values for min & max",
			arguments:   [3]int{9, 0, 0},
			expected:    []byte{0xF0, 0x70, 9, 0, 0, 0, 0, 0xF7},
		},
		{
			description: "Max values for min & max",
			arguments:   [3]int{9, 0x3FFF, 0x3FFF},
			expected:    []byte{0xF0, 0x70, 9, 0x7F, 0x7F, 0x7F, 0x7F, 0xF7},
		},
		{
			description: "Clipped max values for min & max",
			arguments:   [3]int{9, 0xFFFF, 0xFFFF},
			expected:    []byte{0xF0, 0x70, 9, 0x7F, 0x7F, 0x7F, 0x7F, 0xF7},
		},
	}

	for _, test := range tests {
		writeDataMutex.Lock()
		testWriteData.Reset()
		writeDataMutex.Unlock()
		err := b.ServoConfig(test.arguments[0], test.arguments[1], test.arguments[2])
		writeDataMutex.Lock()
		gobottest.Assert(t, testWriteData.Bytes(), test.expected)
		gobottest.Assert(t, err, test.result)
		writeDataMutex.Unlock()
	}
}

func TestProcessSysexData(t *testing.T) {
	sem := make(chan bool)
	b := initTestFirmata()
	b.setConnected(true)
	SetTestReadData([]byte{240, 17, 1, 2, 3, 247})

	b.Once("SysexResponse", func(data interface{}) {
		gobottest.Assert(t, data, []byte{240, 17, 1, 2, 3, 247})
		sem <- true
	})

	b.process()

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("SysexResponse was not published")
	}
}
