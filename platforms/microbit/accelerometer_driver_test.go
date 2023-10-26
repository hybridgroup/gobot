package microbit

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*AccelerometerDriver)(nil)

func initTestAccelerometerDriver() *AccelerometerDriver {
	d := NewAccelerometerDriver(NewBleTestAdaptor())
	return d
}

func TestAccelerometerDriver(t *testing.T) {
	d := initTestAccelerometerDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit Accelerometer"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestAccelerometerDriverStartAndHalt(t *testing.T) {
	d := initTestAccelerometerDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestAccelerometerDriverReadData(t *testing.T) {
	sem := make(chan bool)
	a := NewBleTestAdaptor()
	d := NewAccelerometerDriver(a)
	_ = d.Start()
	_ = d.On(Accelerometer, func(data interface{}) {
		assert.Equal(t, float32(8.738), data.(*AccelerometerData).X)
		assert.Equal(t, float32(8.995), data.(*AccelerometerData).Y)
		assert.Equal(t, float32(9.252), data.(*AccelerometerData).Z)
		sem <- true
	})

	a.TestReceiveNotification([]byte{0x22, 0x22, 0x23, 0x23, 0x24, 0x24}, nil)

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Microbit Event \"Accelerometer\" was not published")
	}
}
