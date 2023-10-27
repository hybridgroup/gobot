package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*BuzzerDriver)(nil)

func initTestBuzzerDriver(conn DigitalWriter) *BuzzerDriver {
	return NewBuzzerDriver(conn, "1")
}

func TestBuzzerDriverDefaultName(t *testing.T) {
	g := initTestBuzzerDriver(newGpioTestAdaptor())
	assert.True(t, strings.HasPrefix(g.Name(), "Buzzer"))
}

func TestBuzzerDriverSetName(t *testing.T) {
	g := initTestBuzzerDriver(newGpioTestAdaptor())
	g.SetName("mybot")
	assert.Equal(t, "mybot", g.Name())
}

func TestBuzzerDriverStart(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	assert.NoError(t, d.Start())
}

func TestBuzzerDriverHalt(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	assert.NoError(t, d.Halt())
}

func TestBuzzerDriverToggle(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	_ = d.Off()
	_ = d.Toggle()
	assert.True(t, d.State())
	_ = d.Toggle()
	assert.False(t, d.State())
}

func TestBuzzerDriverTone(t *testing.T) {
	d := initTestBuzzerDriver(newGpioTestAdaptor())
	assert.NoError(t, d.Tone(100, 0.01))
}

func TestBuzzerDriverOnError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.digitalWriteFunc = func(string, byte) (err error) {
		return errors.New("write error")
	}

	assert.ErrorContains(t, d.On(), "write error")
}

func TestBuzzerDriverOffError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.digitalWriteFunc = func(string, byte) (err error) {
		return errors.New("write error")
	}

	assert.ErrorContains(t, d.Off(), "write error")
}

func TestBuzzerDriverToneError(t *testing.T) {
	a := newGpioTestAdaptor()
	d := initTestBuzzerDriver(a)
	a.digitalWriteFunc = func(string, byte) (err error) {
		return errors.New("write error")
	}

	assert.ErrorContains(t, d.Tone(100, 0.01), "write error")
}
