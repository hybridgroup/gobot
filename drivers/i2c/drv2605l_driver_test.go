package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*DRV2605LDriver)(nil)

func initTestDRV2605LDriverWithStubbedAdaptor() (*DRV2605LDriver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	// Prime adapter reader to make "Start()" call happy
	a.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, uint8(42))
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	d := NewDRV2605LDriver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewDRV2605LDriver(t *testing.T) {
	var di interface{} = NewDRV2605LDriver(newI2cTestAdaptor())
	d, ok := di.(*DRV2605LDriver)
	if !ok {
		t.Errorf("NewDRV2605LDriver() should have returned a *DRV2605LDriver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "DRV2605L"), true)
	gobottest.Assert(t, d.defaultAddress, 0x5a)
}

func TestDRV2605LOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewDRV2605LDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestDRV2605LStart(t *testing.T) {
	d := NewDRV2605LDriver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestDRV2605LHalt(t *testing.T) {
	writeStopPlaybackData := []byte{drv2605RegGo, 0}
	// single-byte-read starts with a write operation to set the register for reading
	// see section 8.5.3.5 of data sheet
	readCurrentStandbyModeData := byte(drv2605RegMode)
	writeNewStandbyModeData := []byte{drv2605RegMode, 42 | drv2605Standby}
	d, a := initTestDRV2605LDriverWithStubbedAdaptor()
	a.written = []byte{}
	gobottest.Assert(t, d.Halt(), nil)
	gobottest.Assert(t, a.written, append(append(writeStopPlaybackData, readCurrentStandbyModeData), writeNewStandbyModeData...))
}

func TestDRV2605LGetPause(t *testing.T) {
	d, _ := initTestDRV2605LDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.GetPauseWaveform(0), uint8(0x80))
	gobottest.Assert(t, d.GetPauseWaveform(1), uint8(0x81))
	gobottest.Assert(t, d.GetPauseWaveform(128), d.GetPauseWaveform(127))
}

func TestDRV2605LSequenceTermination(t *testing.T) {
	d, a := initTestDRV2605LDriverWithStubbedAdaptor()
	a.written = []byte{}
	gobottest.Assert(t, d.SetSequence([]byte{1, 2}), nil)
	gobottest.Assert(t, a.written, []byte{
		drv2605RegWaveSeq1, 1,
		drv2605RegWaveSeq2, 2,
		drv2605RegWaveSeq3, 0,
	})
}

func TestDRV2605LSequenceTruncation(t *testing.T) {
	d, a := initTestDRV2605LDriverWithStubbedAdaptor()
	a.written = []byte{}
	gobottest.Assert(t, d.SetSequence([]byte{1, 2, 3, 4, 5, 6, 7, 8, 99, 100}), nil)
	gobottest.Assert(t, a.written, []byte{
		drv2605RegWaveSeq1, 1,
		drv2605RegWaveSeq2, 2,
		drv2605RegWaveSeq3, 3,
		drv2605RegWaveSeq4, 4,
		drv2605RegWaveSeq5, 5,
		drv2605RegWaveSeq6, 6,
		drv2605RegWaveSeq7, 7,
		drv2605RegWaveSeq8, 8,
	})
}

func TestDRV2605LSetMode(t *testing.T) {
	d, _ := initTestDRV2605LDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.SetMode(DRV2605ModeIntTrig), nil)
}

func TestDRV2605LSetModeReadError(t *testing.T) {
	d, a := initTestDRV2605LDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, d.SetMode(DRV2605ModeIntTrig), errors.New("read error"))
}

func TestDRV2605LSetStandbyMode(t *testing.T) {
	d, _ := initTestDRV2605LDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.SetStandbyMode(true), nil)
}

func TestDRV2605LSetStandbyModeReadError(t *testing.T) {
	d, a := initTestDRV2605LDriverWithStubbedAdaptor()
	a.i2cReadImpl = func(b []byte) (int, error) {
		return 0, errors.New("read error")
	}
	gobottest.Assert(t, d.SetStandbyMode(true), errors.New("read error"))
}

func TestDRV2605LSelectLibrary(t *testing.T) {
	d, _ := initTestDRV2605LDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.SelectLibrary(1), nil)
}

func TestDRV2605LGo(t *testing.T) {
	d, _ := initTestDRV2605LDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Go(), nil)
}
