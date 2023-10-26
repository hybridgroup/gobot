package microbit

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*MagnetometerDriver)(nil)

func initTestMagnetometerDriver() *MagnetometerDriver {
	d := NewMagnetometerDriver(NewBleTestAdaptor())
	return d
}

func TestMagnetometerDriver(t *testing.T) {
	d := initTestMagnetometerDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit Magnetometer"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestMagnetometerDriverStartAndHalt(t *testing.T) {
	d := initTestMagnetometerDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestMagnetometerDriverReadData(t *testing.T) {
	sem := make(chan bool)
	a := NewBleTestAdaptor()
	d := NewMagnetometerDriver(a)
	_ = d.Start()
	_ = d.On(Magnetometer, func(data interface{}) {
		assert.Equal(t, float32(8.738), data.(*MagnetometerData).X)
		assert.Equal(t, float32(8.995), data.(*MagnetometerData).Y)
		assert.Equal(t, float32(9.252), data.(*MagnetometerData).Z)
		sem <- true
	})

	a.TestReceiveNotification([]byte{0x22, 0x22, 0x23, 0x23, 0x24, 0x24}, nil)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Microbit Event \"Magnetometer\" was not published")
	}
}
