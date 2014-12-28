package neurosky

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
)

func initTestNeuroskyDriver() *NeuroskyDriver {
	a := NewNeuroskyAdaptor("bot", "/dev/null")
	a.connect = func(n *NeuroskyAdaptor) (io.ReadWriteCloser, error) {
		return &NullReadWriteCloser{}, nil
	}
	a.Connect()
	return NewNeuroskyDriver(a, "bot")
}

func TestNeuroskyDriver(t *testing.T) {
	d := initTestNeuroskyDriver()
	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Connection().Name(), "bot")
}
func TestNeuroskyDriverStart(t *testing.T) {
	sem := make(chan bool, 0)
	d := initTestNeuroskyDriver()
	gobot.Once(d.Event("error"), func(data interface{}) {
		gobot.Assert(t, data.(error), errors.New("read error"))
		sem <- true
	})

	gobot.Assert(t, len(d.Start()), 0)

	<-time.After(50 * time.Millisecond)
	readError = errors.New("read error")

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
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestNeuroskyDriverParse(t *testing.T) {
	sem := make(chan bool)
	d := initTestNeuroskyDriver()

	// CodeEx
	go func() {
		<-time.After(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 1, 0x55, 0x00}))
	}()

	gobot.On(d.Event("extended"), func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Event \"extended\" was not published")
	}

	// CodeSignalQuality
	go func() {
		<-time.After(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x02, 100, 0x00}))
	}()

	gobot.On(d.Event("signal"), func(data interface{}) {
		gobot.Assert(t, data.(byte), byte(100))
		sem <- true
	})

	<-sem

	// CodeAttention
	go func() {
		<-time.After(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x04, 40, 0x00}))
	}()

	gobot.On(d.Event("attention"), func(data interface{}) {
		gobot.Assert(t, data.(byte), byte(40))
		sem <- true
	})

	<-sem

	// CodeMeditation
	go func() {
		<-time.After(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x05, 60, 0x00}))
	}()

	gobot.On(d.Event("meditation"), func(data interface{}) {
		gobot.Assert(t, data.(byte), byte(60))
		sem <- true
	})

	<-sem

	// CodeBlink
	go func() {
		<-time.After(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x16, 150, 0x00}))
	}()

	gobot.On(d.Event("blink"), func(data interface{}) {
		gobot.Assert(t, data.(byte), byte(150))
		sem <- true
	})

	<-sem

	// CodeWave
	go func() {
		<-time.After(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 4, 0x80, 0x00, 0x40, 0x11, 0x00}))
	}()

	gobot.On(d.Event("wave"), func(data interface{}) {
		gobot.Assert(t, data.(int16), int16(16401))
		sem <- true
	})

	<-sem

	// CodeAsicEEG
	go func() {
		<-time.After(5 * time.Millisecond)
		d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 30, 0x83, 24, 1, 121, 89, 0,
			97, 26, 0, 30, 189, 0, 57, 1, 0, 62, 160, 0, 31, 127, 0, 18, 207, 0, 13,
			108, 0x00}))
	}()

	gobot.On(d.Event("eeg"), func(data interface{}) {
		gobot.Assert(t,
			data.(EEG),
			EEG{
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
