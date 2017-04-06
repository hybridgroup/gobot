package neurosky

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestNeuroskyDriver() *Driver {
	a := NewAdaptor("/dev/null")
	a.connect = func(n *Adaptor) (io.ReadWriteCloser, error) {
		return &NullReadWriteCloser{}, nil
	}
	a.Connect()
	return NewDriver(a)
}

func TestNeuroskyDriver(t *testing.T) {
	d := initTestNeuroskyDriver()
	gobottest.Refute(t, d.Connection(), nil)
}

func TestNeuroskyDriverName(t *testing.T) {
	d := initTestNeuroskyDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Neurosky"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestNeuroskyDriverStart(t *testing.T) {
	sem := make(chan bool, 0)

	rwc := &NullReadWriteCloser{}
	a := NewAdaptor("/dev/null")
	a.connect = func(n *Adaptor) (io.ReadWriteCloser, error) {
		return rwc, nil
	}
	a.Connect()

	d := NewDriver(a)
	e := errors.New("read error")
	d.Once(d.Event(Error), func(data interface{}) {
		gobottest.Assert(t, data.(error), e)
		sem <- true
	})

	gobottest.Assert(t, d.Start(), nil)

	time.Sleep(50 * time.Millisecond)
	rwc.ReadError(e)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		{
			t.Errorf("error was not emitted")
		}

	}
}

func TestNeuroskyDriverHalt(t *testing.T) {
	d := initTestNeuroskyDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestNeuroskyDriverParse(t *testing.T) {
	sem := make(chan bool)
	d := initTestNeuroskyDriver()

	// CodeEx
	go func() {
		time.Sleep(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 1, 0x55, 0x00}))
	}()

	d.On(d.Event(Extended), func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Event \"extended\" was not published")
	}

	// CodeSignalQuality
	go func() {
		time.Sleep(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x02, 100, 0x00}))
	}()

	d.On(d.Event(Signal), func(data interface{}) {
		gobottest.Assert(t, data.(byte), byte(100))
		sem <- true
	})

	<-sem

	// CodeAttention
	go func() {
		time.Sleep(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x04, 40, 0x00}))
	}()

	d.On(d.Event(Attention), func(data interface{}) {
		gobottest.Assert(t, data.(byte), byte(40))
		sem <- true
	})

	<-sem

	// CodeMeditation
	go func() {
		time.Sleep(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x05, 60, 0x00}))
	}()

	d.On(d.Event(Meditation), func(data interface{}) {
		gobottest.Assert(t, data.(byte), byte(60))
		sem <- true
	})

	<-sem

	// CodeBlink
	go func() {
		time.Sleep(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x16, 150, 0x00}))
	}()

	d.On(d.Event(Blink), func(data interface{}) {
		gobottest.Assert(t, data.(byte), byte(150))
		sem <- true
	})

	<-sem

	// CodeWave
	go func() {
		time.Sleep(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 4, 0x80, 0x00, 0x40, 0x11, 0x00}))
	}()

	d.On(d.Event(Wave), func(data interface{}) {
		gobottest.Assert(t, data.(int16), int16(16401))
		sem <- true
	})

	<-sem

	// CodeAsicEEG
	go func() {
		time.Sleep(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 30, 0x83, 24, 1, 121, 89, 0,
			97, 26, 0, 30, 189, 0, 57, 1, 0, 62, 160, 0, 31, 127, 0, 18, 207, 0, 13,
			108, 0x00}))
	}()

	d.On(d.Event(EEG), func(data interface{}) {
		gobottest.Assert(t,
			data.(EEGData),
			EEGData{
				Delta:    1573241,
				Theta:    5832801,
				LoAlpha:  1703966,
				HiAlpha:  12386361,
				LoBeta:   65598,
				HiBeta:   10485791,
				LoGamma:  8323090,
				MidGamma: 13565965,
			})
		sem <- true
	})
	<-sem
}
