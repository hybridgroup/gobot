package microbit

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*TemperatureDriver)(nil)

func initTestTemperatureDriver() *TemperatureDriver {
	d := NewTemperatureDriver(NewBleTestAdaptor())
	return d
}

func TestTemperatureDriver(t *testing.T) {
	d := initTestTemperatureDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit Temperature"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestTemperatureDriverStartAndHalt(t *testing.T) {
	d := initTestTemperatureDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestTemperatureDriverReadData(t *testing.T) {
	sem := make(chan bool)
	a := NewBleTestAdaptor()
	d := NewTemperatureDriver(a)
	_ = d.Start()
	_ = d.On(Temperature, func(data interface{}) {
		assert.Equal(t, int8(0x22), data)
		sem <- true
	})

	a.TestReceiveNotification([]byte{0x22}, nil)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Microbit Event \"Temperature\" was not published")
	}
}
