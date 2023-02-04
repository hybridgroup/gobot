package i2c

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
		t.Errorf("NewAdafruit1109Driver() should have returned a *Adafruit1109Driver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Refute(t, d.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Adafruit1109"), true)
	gobottest.Assert(t, strings.Contains(d.Name(), "MCP23017"), true)
	gobottest.Assert(t, strings.Contains(d.Name(), "HD44780"), true)
	gobottest.Refute(t, d.MCP23017Driver, nil)
	gobottest.Refute(t, d.HD44780Driver, nil)
	gobottest.Refute(t, d.redPin, nil)
	gobottest.Refute(t, d.greenPin, nil)
	gobottest.Refute(t, d.bluePin, nil)
	gobottest.Refute(t, d.selectPin, nil)
	gobottest.Refute(t, d.upPin, nil)
	gobottest.Refute(t, d.downPin, nil)
	gobottest.Refute(t, d.leftPin, nil)
	gobottest.Refute(t, d.rightPin, nil)
	gobottest.Refute(t, d.rwPin, nil)
	gobottest.Refute(t, d.rsPin, nil)
	gobottest.Refute(t, d.enPin, nil)
	gobottest.Refute(t, d.dataPinD4, nil)
	gobottest.Refute(t, d.dataPinD5, nil)
	gobottest.Refute(t, d.dataPinD6, nil)
	gobottest.Refute(t, d.dataPinD7, nil)
}

func TestAdafruit1109Connect(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	gobottest.Assert(t, d.Connect(), nil)
}

func TestAdafruit1109Finalize(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	gobottest.Assert(t, d.Finalize(), nil)
}

func TestAdafruit1109SetName(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	d.SetName("foo")
	gobottest.Assert(t, d.name, "foo")
}

func TestAdafruit1109Start(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	gobottest.Assert(t, d.Start(), nil)
}

func TestAdafruit1109StartWriteErr(t *testing.T) {
	d, adaptor := initTestAdafruit1109WithStubbedAdaptor()
	adaptor.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}
	gobottest.Assert(t, d.Start(), errors.New("write error"))
}

func TestAdafruit1109StartReadErr(t *testing.T) {
	d, adaptor := initTestAdafruit1109WithStubbedAdaptor()
	adaptor.i2cReadImpl = func([]byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, d.Start(), errors.New("MCP write-read: MCP write-ReadByteData(reg=0): read error"))
}

func TestAdafruit1109Halt(t *testing.T) {
	d, _ := initTestAdafruit1109WithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestAdafruit1109DigitalRead(t *testing.T) {
	var tests = map[string]struct {
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
			d.Start()
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, numCallsRead, 1)
			gobottest.Assert(t, len(a.written), 1)
			gobottest.Assert(t, a.written[0], tc.wantReg)
			gobottest.Assert(t, got, 1)
		})
	}
}

func TestAdafruit1109SelectButton(t *testing.T) {
	var tests = map[string]struct {
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
			d.Start()
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, numCallsRead, 1)
			gobottest.Assert(t, got, tc.want)
		})
	}
}

func TestAdafruit1109UpButton(t *testing.T) {
	var tests = map[string]struct {
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
			d.Start()
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, numCallsRead, 1)
			gobottest.Assert(t, got, tc.want)
		})
	}
}

func TestAdafruit1109DownButton(t *testing.T) {
	var tests = map[string]struct {
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
			d.Start()
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, numCallsRead, 1)
			gobottest.Assert(t, got, tc.want)
		})
	}
}

func TestAdafruit1109LeftButton(t *testing.T) {
	var tests = map[string]struct {
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
			d.Start()
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, numCallsRead, 1)
			gobottest.Assert(t, got, tc.want)
		})
	}
}

func TestAdafruit1109RightButton(t *testing.T) {
	var tests = map[string]struct {
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
			d.Start()
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
			gobottest.Assert(t, err, nil)
			gobottest.Assert(t, numCallsRead, 1)
			gobottest.Assert(t, got, tc.want)
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
				gobottest.Assert(t, got, adafruit1109PortPin{port, pin})
			})
		}
	}
}
