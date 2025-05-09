package client

import (
	"bytes"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const semPublishWait = 30 * time.Millisecond

type readWriteCloser struct {
	id string
}

var (
	testWriteData  = bytes.Buffer{}
	writeDataMutex sync.Mutex
)

// do not set this data directly, use always addTestReadData()
var (
	testReadDataMap = make(map[string][]byte)
	rwDataMapMutex  sync.Mutex
)

// arduino uno r3 protocol response "2.3"
var testDataProtocolResponse = []byte{249, 2, 3}

// arduino uno r3 firmware response "StandardFirmata.ino"
var testDataFirmwareResponse = []byte{
	240, 121, 2, 3, 83, 0, 116, 0, 97, 0, 110, 0, 100, 0, 97,
	0, 114, 0, 100, 0, 70, 0, 105, 0, 114, 0, 109, 0, 97, 0, 116, 0, 97, 0, 46,
	0, 105, 0, 110, 0, 111, 0, 247,
}

// arduino uno r3 capabilities response
var testDataCapabilitiesResponse = []byte{
	240, 108, 127, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 3,
	8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1,
	1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0,
	1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 3, 8, 4, 14, 127, 0, 1, 1, 1, 3, 8,
	4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 4, 14, 127, 0, 1, 1, 1, 2,
	10, 127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10, 127, 0, 1, 1, 1, 2, 10,
	127, 0, 1, 1, 1, 2, 10, 6, 1, 127, 0, 1, 1, 1, 2, 10, 6, 1, 127, 247,
}

// arduino uno r3 analog mapping response
var testDataAnalogMappingResponse = []byte{
	240, 106, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127,
	127, 127, 127, 127, 0, 1, 2, 3, 4, 5, 247,
}

func (readWriteCloser) Write(p []byte) (int, error) {
	writeDataMutex.Lock()
	defer writeDataMutex.Unlock()
	return testWriteData.Write(p)
}

func (rwc readWriteCloser) addTestReadData(d []byte) {
	// concurrent read/change of map is not allowed
	rwDataMapMutex.Lock()
	defer rwDataMapMutex.Unlock()

	data, ok := testReadDataMap[rwc.id]
	if !ok {
		data = []byte{}
	}

	data = append(data, d...)
	testReadDataMap[rwc.id] = data
}

func (rwc readWriteCloser) Read(b []byte) (int, error) {
	// concurrent change of map is not allowed
	rwDataMapMutex.Lock()
	defer rwDataMapMutex.Unlock()

	data, ok := testReadDataMap[rwc.id]
	if !ok {
		// there was no content stored before to read
		log.Printf("no content stored in %s", rwc.id)
		return 0, nil
	}
	size := len(b)
	if len(data) < size {
		size = len(data)
	}
	copy(b, data[:size])
	testReadDataMap[rwc.id] = data[size:]
	return size, nil
}

func (readWriteCloser) Close() error {
	return nil
}

func initTestFirmataWithReadWriteCloser(name string, data ...[]byte) (*Client, readWriteCloser) {
	b := New()
	rwc := readWriteCloser{id: name}
	b.connection = rwc

	for _, d := range data {
		rwc.addTestReadData(d)
		_ = b.process()
	}

	b.setConnected(true)
	return b, rwc
}

func TestPins(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse, testDataAnalogMappingResponse)
	assert.Len(t, b.Pins(), 20)
	assert.Len(t, b.analogPins, 6)
}

func TestProtocolVersionQuery(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name())
	require.NoError(t, b.ProtocolVersionQuery())
}

func TestFirmwareQuery(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name())
	require.NoError(t, b.FirmwareQuery())
}

func TestPinStateQuery(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name())
	require.NoError(t, b.PinStateQuery(1))
}

func TestProcessProtocolVersion(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name())
	rwc.addTestReadData(testDataProtocolResponse)

	_ = b.Once(b.Event("ProtocolVersion"), func(data interface{}) {
		assert.Equal(t, "2.3", data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "ProtocolVersion was not published")
	}
}

func TestProcessAnalogRead0(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse, testDataAnalogMappingResponse)
	rwc.addTestReadData([]byte{0xE0, 0x23, 0x05})

	_ = b.Once(b.Event("AnalogRead0"), func(data interface{}) {
		assert.Equal(t, 675, data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "AnalogRead0 was not published")
	}
}

func TestProcessAnalogRead1(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse, testDataAnalogMappingResponse)
	rwc.addTestReadData([]byte{0xE1, 0x23, 0x06})

	_ = b.Once(b.Event("AnalogRead1"), func(data interface{}) {
		assert.Equal(t, 803, data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "AnalogRead1 was not published")
	}
}

func TestProcessDigitalRead2(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse)
	b.pins[2].Mode = Input
	rwc.addTestReadData([]byte{0x90, 0x04, 0x00})

	_ = b.Once(b.Event("DigitalRead2"), func(data interface{}) {
		assert.Equal(t, 1, data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "DigitalRead2 was not published")
	}
}

func TestProcessDigitalRead4(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse)
	b.pins[4].Mode = Input
	rwc.addTestReadData([]byte{0x90, 0x16, 0x00})

	_ = b.Once(b.Event("DigitalRead4"), func(data interface{}) {
		assert.Equal(t, 1, data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "DigitalRead4 was not published")
	}
}

func TestDigitalWrite(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse)
	require.NoError(t, b.DigitalWrite(13, 0))
}

func TestSetPinMode(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse)
	require.NoError(t, b.SetPinMode(13, Output))
}

func TestAnalogWrite(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse)
	require.NoError(t, b.AnalogWrite(0, 128))
}

func TestReportAnalog(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name())
	require.NoError(t, b.ReportAnalog(0, 1))
	require.NoError(t, b.ReportAnalog(0, 0))
}

func TestProcessPinState13(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name(), testDataCapabilitiesResponse, testDataAnalogMappingResponse)
	rwc.addTestReadData([]byte{240, 110, 13, 1, 1, 247})

	_ = b.Once(b.Event("PinState13"), func(data interface{}) {
		assert.Equal(t, Pin{[]int{0, 1, 4}, 1, 0, 1, 127}, data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "PinState13 was not published")
	}
}

func TestI2cConfig(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name())
	require.NoError(t, b.I2cConfig(100))
}

func TestI2cWrite(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name())
	require.NoError(t, b.I2cWrite(0x00, []byte{0x01, 0x02}))
}

func TestI2cRead(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name())
	require.NoError(t, b.I2cRead(0x00, 10))
}

func TestWriteSysex(t *testing.T) {
	b, _ := initTestFirmataWithReadWriteCloser(t.Name())
	require.NoError(t, b.WriteSysex([]byte{0x01, 0x02}))
}

func TestProcessI2cReply(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name())
	rwc.addTestReadData([]byte{240, 119, 9, 0, 0, 0, 24, 1, 1, 0, 26, 1, 247})

	_ = b.Once(b.Event("I2cReply"), func(data interface{}) {
		assert.Equal(t, I2cReply{
			Address:  9,
			Register: 0,
			Data:     []byte{152, 1, 154},
		}, data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "I2cReply was not published")
	}
}

func TestProcessFirmwareQuery(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name())
	rwc.addTestReadData(testDataFirmwareResponse)

	_ = b.Once(b.Event("FirmwareQuery"), func(data interface{}) {
		assert.Equal(t, "StandardFirmata.ino", data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "FirmwareQuery was not published")
	}
}

func TestProcessStringData(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name())
	rwc.addTestReadData(append([]byte{240, 0x71}, append([]byte("Hello Firmata!"), 247)...))

	_ = b.Once(b.Event("StringData"), func(data interface{}) {
		assert.Equal(t, "Hello Firmata!", data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "StringData was not published")
	}
}

func TestConnect(t *testing.T) {
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name())
	b.setConnected(false)

	rwc.addTestReadData(testDataProtocolResponse)

	_ = b.Once(b.Event("ProtocolVersion"), func(data interface{}) {
		rwc.addTestReadData(testDataFirmwareResponse)
	})

	_ = b.Once(b.Event("FirmwareQuery"), func(data interface{}) {
		rwc.addTestReadData(testDataCapabilitiesResponse)
	})

	_ = b.Once(b.Event("CapabilityQuery"), func(data interface{}) {
		rwc.addTestReadData(testDataAnalogMappingResponse)
	})

	_ = b.Once(b.Event("AnalogMappingQuery"), func(data interface{}) {
		rwc.addTestReadData(testDataProtocolResponse)
	})

	require.NoError(t, b.Connect(rwc))
	time.Sleep(150 * time.Millisecond)
	require.NoError(t, b.Disconnect())
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
		assert.Equal(t, test.expected, testWriteData.Bytes())
		assert.Equal(t, test.result, err)
		writeDataMutex.Unlock()
	}
}

func TestProcessSysexData(t *testing.T) {
	sem := make(chan bool)
	b, rwc := initTestFirmataWithReadWriteCloser(t.Name())
	rwc.addTestReadData([]byte{240, 17, 1, 2, 3, 247})

	_ = b.Once("SysexResponse", func(data interface{}) {
		assert.Equal(t, []byte{240, 17, 1, 2, 3, 247}, data)
		sem <- true
	})

	_ = b.process()

	select {
	case <-sem:
	case <-time.After(semPublishWait):
		require.Fail(t, "SysexResponse was not published")
	}
}
