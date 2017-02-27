package i2c

import (
	"bytes"
	"encoding/binary"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*DRV2605LDriver)(nil)

// --------- HELPERS

func initTestDriverAndAdaptor() (*DRV2605LDriver, *i2cTestAdaptor) {
	adaptor := newI2cTestAdaptor()
	// Prime adapter reader to make "Start()" call happy
	adaptor.i2cReadImpl = func(b []byte) (int, error) {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, uint8(42))
		copy(b, buf.Bytes())
		return buf.Len(), nil
	}
	return NewDRV2605LDriver(adaptor), adaptor
}

// --------- TESTS

func TestDRV2605LDriver(t *testing.T) {
	d, _ := initTestDriverAndAdaptor()
	gobottest.Refute(t, d.Connection(), nil)
}

func TestDRV2605LDriverStart(t *testing.T) {
	d, _ := initTestDriverAndAdaptor()
	gobottest.Assert(t, d.Start(), nil)
}

func TestDRV2605LDriverHalt(t *testing.T) {
	d, adaptor := initTestDriverAndAdaptor()
	gobottest.Assert(t, d.Start(), nil)
	adaptor.written = []byte{}
	gobottest.Assert(t, d.Halt(), nil)
	gobottest.Assert(t, adaptor.written, []byte{drv2605RegGo, 0, drv2605RegMode, 42 | drv2605Standby})
}

func TestDRV2605LDriverGetPause(t *testing.T) {
	d, _ := initTestDriverAndAdaptor()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.GetPauseWaveform(0), uint8(0x80))
	gobottest.Assert(t, d.GetPauseWaveform(1), uint8(0x81))
	gobottest.Assert(t, d.GetPauseWaveform(128), d.GetPauseWaveform(127))
}

func TestDRV2605LDriverSequenceTermination(t *testing.T) {
	d, adaptor := initTestDriverAndAdaptor()
	gobottest.Assert(t, d.Start(), nil)
	adaptor.written = []byte{}
	gobottest.Assert(t, d.SetSequence([]byte{1, 2}), nil)
	gobottest.Assert(t, adaptor.written, []byte{
		drv2605RegWaveSeq1, 1,
		drv2605RegWaveSeq2, 2,
		drv2605RegWaveSeq3, 0,
	})
}

func TestDRV2605LDriverSequenceTruncation(t *testing.T) {
	d, adaptor := initTestDriverAndAdaptor()
	gobottest.Assert(t, d.Start(), nil)
	adaptor.written = []byte{}
	gobottest.Assert(t, d.SetSequence([]byte{1, 2, 3, 4, 5, 6, 7, 8, 99, 100}), nil)
	gobottest.Assert(t, adaptor.written, []byte{
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

func TestDRV2605LDriverSetName(t *testing.T) {
	d, _ := initTestDriverAndAdaptor()
	d.SetName("TESTME")
	gobottest.Assert(t, d.Name(), "TESTME")
}

func TestDRV2605DriverOptions(t *testing.T) {
	d := NewDRV2605LDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}
