package i2c

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Adafruit1109Driver)(nil)

func initTestAdafruit1109WithStubbedAdaptor() (*Adafruit1109Driver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	return NewAdafruit1109Driver(adaptor), adaptor
}

func TestNewAdafruit1109Driver(t *testing.T) {
	var di interface{} = NewAdafruit1109Driver(newI2cTestAdaptor())
	d, ok := di.(*Adafruit1109Driver)
	if !ok {
		require.Fail(t, "NewAdafruit1109Driver() should have returned a *Adafruit1109Driver")
	}
	assert.NotNil(t, d.Driver)
	assert.NotNil(t, d.Connection())
	assert.True(t, strings.HasPrefix(d.Name(), "Adafruit1109"))
	assert.Contains(t, d.Name(), "MCP23017")
	assert.Contains(t, d.Name(), "HD44780")
	assert.NotNil(t, d.MCP23017Driver)
	assert.NotNil(t, d.HD44780Driver)
	assert.NotNil(t, d.redPin)
	assert.NotNil(t, d.greenPin)
	assert.NotNil(t, d.bluePin)
	assert.NotNil(t, d.selectPin)
	assert.NotNil(t, d.upPin)
	assert.NotNil(t, d.downPin)
	assert.NotNil(t, d.leftPin)
	assert.NotNil(t, d.rightPin)
	assert.NotNil(t, d.rwPin)
	assert.NotNil(t, d.rsPin)
	assert.NotNil(t, d.enPin)
	assert.NotNil(t, d.dataPinD4)
	assert.NotNil(t, d.dataPinD5)
	assert.NotNil(t, d.dataPinD6)
	assert.NotNil(t, d.dataPinD7)
}

func TestAdafruit1109Connect(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	require.NoError(t, d.Connect())
}

func TestAdafruit1109Finalize(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	require.NoError(t, d.Finalize())
}

func TestAdafruit1109SetName(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	d.SetName("foo")
	assert.Equal(t, "foo", d.name)
}

func TestAdafruit1109Start(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	require.NoError(t, d.Start())
}

func TestAdafruit1109StartWriteErr(t *testing.T) {
	d, adaptor := initTestAdafruit1109WithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	require.ErrorContains(t, d.Start(), "write error")
}

func TestAdafruit1109StartReadErr(t *testing.T) {
	d, adaptor := initTestAdafruit1109WithStubbedAdaptor()
	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	require.ErrorContains(t, d.Start(), "MCP write-read: MCP write-ReadByteData(reg=0): read error")
}

func TestAdafruit1109Halt(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Halt())
}

func TestAdafruit1109DigitalRead(t *testing.T) {
	tests := map[string]struct {
		read    uint8
		wantReg uint8
	}{
		"A_0": {read: 0x01, wantReg: 0x12},
		"A_1": {read: 0x02, wantReg: 0x12},
		"A_2": {read: 0x04, wantReg: 0x12},
		"A_3": {read: 0x08, wantReg: 0x12},
		"A_4": {read: 0x10, wantReg: 0x12},
		"A_5": {read: 0x20, wantReg: 0x12},
		"A_6": {read: 0x40, wantReg: 0x12},
		"A_7": {read: 0x80, wantReg: 0x12},
		"B_0": {read: 0x01, wantReg: 0x13},
		"B_1": {read: 0x02, wantReg: 0x13},
		"B_2": {read: 0x04, wantReg: 0x13},
		"B_3": {read: 0x08, wantReg: 0x13},
		"B_4": {read: 0x10, wantReg: 0x13},
		"B_5": {read: 0x20, wantReg: 0x13},
		"B_6": {read: 0x40, wantReg: 0x13},
		"B_7": {read: 0x80, wantReg: 0x13},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestAdafruit1109WithStubbedAdaptor()
			_ = d.Start()
			a.written = []byte{} // reset writes of Start() and former test
			// arrange reads
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				b[0] = tc.read
				return len(b), nil
			}
			// act
			got, err := d.DigitalRead(name)
			// assert
			require.NoError(t, err)
			assert.Equal(t, 1, numCallsRead)
			assert.Len(t, a.written, 1)
			assert.Equal(t, tc.wantReg, a.written[0])
			assert.Equal(t, 1, got)
		})
	}
}

func TestAdafruit1109SelectButton(t *testing.T) {
	tests := map[string]struct {
		read uint8
		want uint8
	}{
		"A0_not_pressed": {read: 0xFE, want: 0},
		"A0_pressed":     {read: 0x01, want: 1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestAdafruit1109WithStubbedAdaptor()
			_ = d.Start()
			// arrange reads
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				b[0] = tc.read
				return len(b), nil
			}
			// act
			got, err := d.SelectButton()
			// assert
			require.NoError(t, err)
			assert.Equal(t, 1, numCallsRead)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAdafruit1109UpButton(t *testing.T) {
	tests := map[string]struct {
		read uint8
		want uint8
	}{
		"A3_not_pressed": {read: 0xF7, want: 0},
		"A3_pressed":     {read: 0x08, want: 1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestAdafruit1109WithStubbedAdaptor()
			_ = d.Start()
			// arrange reads
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				b[0] = tc.read
				return len(b), nil
			}
			// act
			got, err := d.UpButton()
			// assert
			require.NoError(t, err)
			assert.Equal(t, 1, numCallsRead)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAdafruit1109DownButton(t *testing.T) {
	tests := map[string]struct {
		read uint8
		want uint8
	}{
		"A2_not_pressed": {read: 0xFB, want: 0},
		"A2_pressed":     {read: 0x04, want: 1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestAdafruit1109WithStubbedAdaptor()
			_ = d.Start()
			// arrange reads
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				b[0] = tc.read
				return len(b), nil
			}
			// act
			got, err := d.DownButton()
			// assert
			require.NoError(t, err)
			assert.Equal(t, 1, numCallsRead)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAdafruit1109LeftButton(t *testing.T) {
	tests := map[string]struct {
		read uint8
		want uint8
	}{
		"A4_not_pressed": {read: 0xEF, want: 0},
		"A4_pressed":     {read: 0x10, want: 1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestAdafruit1109WithStubbedAdaptor()
			_ = d.Start()
			// arrange reads
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				b[0] = tc.read
				return len(b), nil
			}
			// act
			got, err := d.LeftButton()
			// assert
			require.NoError(t, err)
			assert.Equal(t, 1, numCallsRead)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAdafruit1109RightButton(t *testing.T) {
	tests := map[string]struct {
		read uint8
		want uint8
	}{
		"A1_not_pressed": {read: 0xFD, want: 0},
		"A1_pressed":     {read: 0x02, want: 1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			d, a := initTestAdafruit1109WithStubbedAdaptor()
			_ = d.Start()
			// arrange reads
			numCallsRead := 0
			a.i2cReadImpl = func(b []byte) (int, error) {
				numCallsRead++
				b[0] = tc.read
				return len(b), nil
			}
			// act
			got, err := d.RightButton()
			// assert
			require.NoError(t, err)
			assert.Equal(t, 1, numCallsRead)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAdafruit1109_parseID(t *testing.T) {
	// arrange
	ports := []string{"A", "B"}
	for _, port := range ports {
		for pin := uint8(0); pin <= 7; pin++ {
			id := fmt.Sprintf("%s_%d", port, pin)
			t.Run(id, func(t *testing.T) {
				// act
				got := adafruit1109ParseID(id)
				// assert
				assert.Equal(t, adafruit1109PortPin{port, pin}, got)
			})
		}
	}
}
