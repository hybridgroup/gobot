package i2c

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*TH02Driver)(nil)

func initTestTH02DriverWithStubbedAdaptor() (*TH02Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	driver := NewTH02Driver(adaptor)
	if err := driver.Start(); err != nil {
		panic(err)
	}
	return driver, adaptor
}

func TestNewTH02Driver(t *testing.T) {
	var di interface{} = NewTH02Driver(newI2cTestAdaptor())
	d, ok := di.(*TH02Driver)
	if !ok {
		require.Fail(t, "NewTH02Driver() should have returned a *NewTH02Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "TH02"))
	assert.Equal(t, 0x40, d.defaultAddress)
}

func TestTH02Options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common options.
	// Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewTH02Driver(newI2cTestAdaptor(), WithBus(2), WithAddress(0x42))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
	assert.Equal(t, 0x42, d.GetAddressOrDefault(0x33))
}

func TestTH02SetAccuracy(t *testing.T) {
	b := NewTH02Driver(newI2cTestAdaptor())

	if b.SetAccuracy(0x42); b.Accuracy() != TH02HighAccuracy {
		t.Error("Setting an invalid accuracy should resolve to TH02HighAccuracy")
	}

	if b.SetAccuracy(TH02LowAccuracy); b.Accuracy() != TH02LowAccuracy {
		t.Error("Expected setting low accuracy to actually set to low accuracy")
	}

	if acc := b.Accuracy(); acc != TH02LowAccuracy {
		require.Fail(t, "Accuracy() didn't return what was expected")
	}
}

func TestTH02WithFastMode(t *testing.T) {
	tests := map[string]struct {
		value int
		want  bool
	}{
		"fast_on_for >0":  {value: 1, want: true},
		"fast_off_for =0": {value: 0, want: false},
		"fast_off_for <0": {value: -1, want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d := NewTH02Driver(newI2cTestAdaptor())
			// act
			WithTH02FastMode(tc.value)(d)
			// assert
			assert.Equal(t, tc.want, d.fastMode)
		})
	}
}

func TestTH02FastMode(t *testing.T) {
	// sequence to read the fast mode status
	// * write config register address (0x03)
	// * read register content
	// * if sixth bit (D5) is set, the fast mode is configured on, otherwise off
	tests := map[string]struct {
		read uint8
		want bool
	}{
		"fast on":  {read: 0x20, want: true},
		"fast off": {read: ^uint8(0x20), want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestTH02DriverWithStubbedAdaptor()
			a.i2cReadImpl = func(b []byte) (int, error) {
				b[0] = tc.read
				return len(b), nil
			}
			// act
			got, err := d.FastMode()
			// assert
			require.NoError(t, err)
			assert.Len(t, a.written, 1)
			assert.Equal(t, uint8(0x03), a.written[0])
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTH02SetHeater(t *testing.T) {
	// sequence to set the heater status
	// * set the local heater state
	// * write config register address (0x03)
	// * prepare config value by set/reset the heater bit (0x02, D1)
	// * write the config value
	tests := map[string]struct {
		heater bool
		want   uint8
	}{
		"heater on":  {heater: true, want: 0x02},
		"heater off": {heater: false, want: 0x00},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestTH02DriverWithStubbedAdaptor()
			// act
			err := d.SetHeater(tc.heater)
			// assert
			require.NoError(t, err)
			assert.Equal(t, tc.heater, d.heating)
			assert.Len(t, a.written, 2)
			assert.Equal(t, uint8(0x03), a.written[0])
			assert.Equal(t, tc.want, a.written[1])
		})
	}
}

func TestTH02Heater(t *testing.T) {
	// sequence to read the heater status
	// * write config register address (0x03)
	// * read register content
	// * if second bit (D1) is set, the heater is configured on, otherwise off
	tests := map[string]struct {
		read uint8
		want bool
	}{
		"heater on":  {read: 0x02, want: true},
		"heater off": {read: ^uint8(0x02), want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestTH02DriverWithStubbedAdaptor()
			a.i2cReadImpl = func(b []byte) (int, error) {
				b[0] = tc.read
				return len(b), nil
			}
			// act
			got, err := d.Heater()
			// assert
			require.NoError(t, err)
			assert.Len(t, a.written, 1)
			assert.Equal(t, uint8(0x03), a.written[0])
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTH02SerialNumber(t *testing.T) {
	// sequence to read SN
	// * write identification register address (0x11)
	// * read register content
	// * use the higher nibble of byte

	// arrange
	d, a := initTestTH02DriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		b[0] = 0x4F
		return len(b), nil
	}
	want := uint8(0x04)
	// act
	sn, err := d.SerialNumber()
	// assert
	require.NoError(t, err)
	assert.Len(t, a.written, 1)
	assert.Equal(t, uint8(0x11), a.written[0])
	assert.Equal(t, want, sn)
}

func TestTH02Sample(t *testing.T) {
	// sequence to read values
	// * write config register address (0x03)
	// * prepare config bits (START, HEAT, TEMP, FAST)
	// * write config register with config
	// * write status register address (0x00)
	// * read until value is "0" (means ready)
	// * write data register MSB address (0x01)
	// * read 2 bytes little-endian (MSB, LSB)
	// * shift and scale
	//    RH: 4 bits shift right, RH[%]=RH/16-24
	//    T:  2 bits shift right, T[°C]=T/32-50

	// test table according to data sheet page 15, 17
	// operating range of the temperature sensor is -40..85 °C (F-grade 0..70 °C)
	tests := map[string]struct {
		hData  uint16
		tData  uint16
		wantRH float32
		wantT  float32
	}{
		"RH 0, T -40": {
			hData: 0x0180, wantRH: 0.0,
			tData: 0x0140, wantT: -40.0,
		},
		"RH 10, T -20": {
			hData: 0x0220, wantRH: 10.0,
			tData: 0x03C0, wantT: -20.0,
		},
		"RH 20, T -10": {
			hData: 0x02C0, wantRH: 20.0,
			tData: 0x0500, wantT: -10.0,
		},
		"RH 30, T 0": {
			hData: 0x0360, wantRH: 30.0,
			tData: 0x0640, wantT: 0.0,
		},
		"RH 40, T 10": {
			hData: 0x0400, wantRH: 40.0,
			tData: 0x0780, wantT: 10.0,
		},
		"RH 50, T 20": {
			hData: 0x04A0, wantRH: 50.0,
			tData: 0x08C0, wantT: 20.0,
		},
		"RH 60, T 30": {
			hData: 0x0540, wantRH: 60.0,
			tData: 0x0A00, wantT: 30.0,
		},
		"RH 70, T 40": {
			hData: 0x05E0, wantRH: 70.0,
			tData: 0x0B40, wantT: 40.0,
		},
		"RH 80, T 50": {
			hData: 0x0680, wantRH: 80.0,
			tData: 0x0C80, wantT: 50.0,
		},
		"RH 90, T 60": {
			hData: 0x0720, wantRH: 90.0,
			tData: 0x0DC0, wantT: 60.0,
		},
		"RH 100, T 70": {
			hData: 0x07C0, wantRH: 100.0,
			tData: 0x0F00, wantT: 70.0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestTH02DriverWithStubbedAdaptor()
			var reg uint8
			var regVal uint8
			a.i2cWriteImpl = func(b []byte) (int, error) {
				reg = b[0]
				if len(b) == 2 {
					regVal = b[1]
				}
				return len(b), nil
			}
			a.i2cReadImpl = func(b []byte) (int, error) {
				switch reg {
				case 0x00:
					// status
					b[0] = 0
				case 0x01:
					// data register MSB
					var data uint16
					if (regVal & 0x10) == 0x10 {
						// temperature
						data = tc.tData << 2 // data sheet values are after shift 2 bits
					} else {
						// humidity
						data = tc.hData << 4 // data sheet values are after shift 4 bits
					}
					b[0] = byte(data >> 8)   // first read MSB from register 0x01
					b[1] = byte(data & 0xFF) // second read LSB from register 0x02
				default:
					assert.Equal(t, "only register 0 and 1 expected", fmt.Sprintf("unexpected register %d", reg))
					return 0, nil
				}
				return len(b), nil
			}
			// act
			temp, rh, err := d.Sample()
			// assert
			require.NoError(t, err)
			assert.InDelta(t, tc.wantRH, rh, 0.0)
			assert.InDelta(t, tc.wantT, temp, 0.0)
		})
	}
}

func TestTH02_readData(t *testing.T) {
	d, a := initTestTH02DriverWithStubbedAdaptor()

	var callCounter int

	tests := map[string]struct {
		rd      func([]byte) (int, error)
		wr      func([]byte) (int, error)
		rtn     uint16
		wantErr error
	}{
		"example RH": {
			rd: func(b []byte) (int, error) {
				callCounter++
				if callCounter == 1 {
					// read for ready
					b[0] = 0x00
				} else {
					copy(b, []byte{0x07, 0xC0})
				}
				return len(b), nil
			},
			rtn: 1984,
		},
		"example T": {
			rd: func(b []byte) (int, error) {
				callCounter++
				if callCounter == 1 {
					// read for ready
					b[0] = 0x00
				} else {
					copy(b, []byte{0x12, 0xC0})
				}
				return len(b), nil
			},
			rtn: 4800,
		},
		"timeout - no wait for ready": {
			rd: func(b []byte) (int, error) {
				time.Sleep(200 * time.Millisecond)
				// simulate not ready
				b[0] = 0x01
				return len(b), nil
			},
			wantErr: fmt.Errorf("timeout on \\RDY"),
			rtn:     0,
		},
		"unable to write status register": {
			rd: func(b []byte) (int, error) {
				callCounter++
				if callCounter == 1 {
					// read for ready
					b[0] = 0x00
				}
				return len(b), nil
			},
			wr: func(b []byte) (int, error) {
				return len(b), fmt.Errorf("an write error")
			},
			wantErr: fmt.Errorf("timeout on \\RDY"),
			rtn:     0,
		},
		"unable to write data register": {
			rd: func(b []byte) (int, error) {
				callCounter++
				if callCounter == 1 {
					// read for ready
					b[0] = 0x00
				}
				return len(b), nil
			},
			wr: func(b []byte) (int, error) {
				if len(b) == 1 && b[0] == 0x00 {
					// register of ready check
					return len(b), nil
				}
				// data register
				return len(b), fmt.Errorf("Nope")
			},
			wantErr: fmt.Errorf("Nope"),
			rtn:     0,
		},
		"unable to read doesn't provide enough data": {
			rd: func(b []byte) (int, error) {
				callCounter++
				if callCounter == 1 {
					// read for ready
					b[0] = 0x00
				} else {
					b = []byte{0x01}
				}
				return len(b), nil
			},
			wantErr: fmt.Errorf("Read 1 bytes from device by i2c helpers, expected 2"),
			rtn:     0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a.i2cReadImpl = tc.rd
			if tc.wr != nil {
				oldwr := a.i2cWriteImpl
				a.i2cWriteImpl = tc.wr
				defer func() { a.i2cWriteImpl = oldwr }()
			}
			callCounter = 0
			// act
			got, err := d.waitAndReadData()
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.rtn, got)
		})
	}
}

func TestTH02_waitForReadyFailOnTimeout(t *testing.T) {
	d, a := initTestTH02DriverWithStubbedAdaptor()

	a.i2cReadImpl = func(b []byte) (int, error) {
		time.Sleep(50 * time.Millisecond)
		b[0] = 0x01
		return len(b), nil
	}

	timeout := 10 * time.Microsecond
	if err := d.waitForReady(&timeout); err == nil {
		t.Error("Expected a timeout error")
	}
}

func TestTH02_waitForReadyFailOnReadError(t *testing.T) {
	d, a := initTestTH02DriverWithStubbedAdaptor()

	a.i2cReadImpl = func(b []byte) (int, error) {
		time.Sleep(50 * time.Millisecond)
		b[0] = 0x00
		wrongLength := 2
		return wrongLength, nil
	}

	timeout := 10 * time.Microsecond
	if err := d.waitForReady(&timeout); err == nil {
		t.Error("Expected a timeout error")
	}
}

func TestTH02_createConfig(t *testing.T) {
	d := &TH02Driver{}

	tests := map[string]struct {
		meas     bool
		fast     bool
		readTemp bool
		heating  bool
		want     byte
	}{
		"meas, no fast, RH, no heating":    {meas: true, fast: false, readTemp: false, heating: false, want: 0x01},
		"meas, no fast, RH, heating":       {meas: true, fast: false, readTemp: false, heating: true, want: 0x03},
		"meas, no fast, TE, no heating":    {meas: true, fast: false, readTemp: true, heating: false, want: 0x11},
		"meas, no fast, TE, heating":       {meas: true, fast: false, readTemp: true, heating: true, want: 0x13},
		"meas, fast, RH, no heating":       {meas: true, fast: true, readTemp: false, heating: false, want: 0x21},
		"meas, fast, RH, heating":          {meas: true, fast: true, readTemp: false, heating: true, want: 0x23},
		"meas, fast, TE, no heating":       {meas: true, fast: true, readTemp: true, heating: false, want: 0x31},
		"meas, fast, TE, heating":          {meas: true, fast: true, readTemp: true, heating: true, want: 0x33},
		"no meas, no fast, RH, no heating": {meas: false, fast: false, readTemp: false, heating: false, want: 0x00},
		"no meas, no fast, RH, heating":    {meas: false, fast: false, readTemp: false, heating: true, want: 0x02},
		"no meas, no fast, TE, no heating": {meas: false, fast: false, readTemp: true, heating: false, want: 0x00},
		"no meas, no fast, TE, heating":    {meas: false, fast: false, readTemp: true, heating: true, want: 0x02},
		"no meas, fast, RH, no heating":    {meas: false, fast: true, readTemp: false, heating: false, want: 0x00},
		"no meas, fast, RH, heating":       {meas: false, fast: true, readTemp: false, heating: true, want: 0x02},
		"no meas, fast, TE, no heating":    {meas: false, fast: true, readTemp: true, heating: false, want: 0x00},
		"no meas, fast, TE, heating":       {meas: false, fast: true, readTemp: true, heating: true, want: 0x02},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d.fastMode = tc.fast
			d.heating = tc.heating
			got := d.createConfig(tc.meas, tc.readTemp)
			assert.Equal(t, tc.want, got)
		})
	}
}
