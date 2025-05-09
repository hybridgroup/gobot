//nolint:forcetypeassert // ok here
package neurosky

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/serial/testutil"
)

var _ gobot.Driver = (*MindWaveDriver)(nil)

func initTestNeuroskyDriver() *MindWaveDriver {
	a := testutil.NewSerialTestAdaptor()
	_ = a.Connect()
	return NewMindWaveDriver(a)
}

func TestNeuroskyDriver(t *testing.T) {
	d := initTestNeuroskyDriver()
	assert.NotNil(t, d.Connection())
}

func TestNeuroskyDriverName(t *testing.T) {
	d := initTestNeuroskyDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "MindWave"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestNeuroskyDriverStart(t *testing.T) {
	sem := make(chan bool)

	a := testutil.NewSerialTestAdaptor()
	_ = a.Connect()
	a.SetSimulateReadError(true)

	d := NewMindWaveDriver(a)
	e := errors.New("read error")
	_ = d.Once(d.Event(ErrorEvent), func(data interface{}) {
		assert.Equal(t, e, data.(error))
		sem <- true
	})

	require.NoError(t, d.Start())
	time.Sleep(50 * time.Millisecond)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		{
			require.Fail(t, "error was not emitted")
		}

	}
}

func TestNeuroskyDriverHalt(t *testing.T) {
	d := initTestNeuroskyDriver()
	require.NoError(t, d.Halt())
}

func TestNeuroskyDriverParse(t *testing.T) {
	sem := make(chan bool)
	d := initTestNeuroskyDriver()

	// mindWaveCodeEx
	go func() {
		time.Sleep(5 * time.Millisecond)
		_ = d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 1, 0x55, 0x00}))
	}()

	_ = d.On(d.Event(ExtendedEvent), func(data interface{}) {
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "Event \"extended\" was not published")
	}

	// mindWaveCodeSignalQuality
	go func() {
		time.Sleep(5 * time.Millisecond)
		_ = d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x02, 100, 0x00}))
	}()

	_ = d.On(d.Event(SignalEvent), func(data interface{}) {
		assert.Equal(t, byte(100), data.(byte))
		sem <- true
	})

	<-sem

	// mindWaveCodeAttention
	go func() {
		time.Sleep(5 * time.Millisecond)
		_ = d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x04, 40, 0x00}))
	}()

	_ = d.On(d.Event(AttentionEvent), func(data interface{}) {
		assert.Equal(t, byte(40), data.(byte))
		sem <- true
	})

	<-sem

	// mindWaveCodeMeditation
	go func() {
		time.Sleep(5 * time.Millisecond)
		_ = d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x05, 60, 0x00}))
	}()

	_ = d.On(d.Event(MeditationEvent), func(data interface{}) {
		assert.Equal(t, byte(60), data.(byte))
		sem <- true
	})

	<-sem

	// mindWaveCodeBlink
	go func() {
		time.Sleep(5 * time.Millisecond)
		_ = d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 2, 0x16, 150, 0x00}))
	}()

	_ = d.On(d.Event(BlinkEvent), func(data interface{}) {
		assert.Equal(t, byte(150), data.(byte))
		sem <- true
	})

	<-sem

	// mindWaveCodeWave
	go func() {
		time.Sleep(5 * time.Millisecond)
		_ = d.parse(bytes.NewBuffer([]byte{0xAA, 0xAA, 4, 0x80, 0x00, 0x40, 0x11, 0x00}))
	}()

	_ = d.On(d.Event(WaveEvent), func(data interface{}) {
		assert.Equal(t, int16(16401), data.(int16))
		sem <- true
	})

	<-sem

	// mindWaveCodeAsicEEG
	go func() {
		time.Sleep(5 * time.Millisecond)
		_ = d.parse(bytes.NewBuffer([]byte{
			0xAA, 0xAA, 30, 0x83, 24, 1, 121, 89, 0,
			97, 26, 0, 30, 189, 0, 57, 1, 0, 62, 160, 0, 31, 127, 0, 18, 207, 0, 13,
			108, 0x00,
		}))
	}()

	_ = d.On(d.Event(EEGEvent), func(data interface{}) {
		assert.Equal(t,
			MindWaveEEGData{
				Delta:    1573241,
				Theta:    5832801,
				LoAlpha:  1703966,
				HiAlpha:  12386361,
				LoBeta:   65598,
				HiBeta:   10485791,
				LoGamma:  8323090,
				MidGamma: 13565965,
			},
			data.(MindWaveEEGData))
		sem <- true
	})
	<-sem
}
