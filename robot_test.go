package gobot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRobotConnectionEach(t *testing.T) {
	r := newTestRobot("Robot1")

	i := 0
	r.Connections().Each(func(conn Connection) {
		i++
	})
	assert.Equal(t, i, r.Connections().Len())
}

func TestRobotToJSON(t *testing.T) {
	r := newTestRobot("Robot99")
	r.AddCommand("test_function", func(params map[string]interface{}) interface{} {
		return nil
	})
	json := NewJSONRobot(r)
	assert.Equal(t, r.Devices().Len(), len(json.Devices))
	assert.Equal(t, len(r.Commands()), len(json.Commands))
}

func TestRobotDevicesToJSON(t *testing.T) {
	r := newTestRobot("Robot99")
	json := NewJSONRobot(r)
	assert.Equal(t, r.Devices().Len(), len(json.Devices))
	assert.Equal(t, "Device1", json.Devices[0].Name)
	assert.Equal(t, "*gobot.testDriver", json.Devices[0].Driver)
	assert.Equal(t, "Connection1", json.Devices[0].Connection)
	assert.Equal(t, 1, len(json.Devices[0].Commands))
}

func TestRobotStart(t *testing.T) {
	r := newTestRobot("Robot99")
	assert.NoError(t, r.Start())
	assert.NoError(t, r.Stop())
	assert.False(t, r.Running())
}

func TestRobotStartAutoRun(t *testing.T) {
	adaptor1 := newTestAdaptor("Connection1", "/dev/null")
	driver1 := newTestDriver(adaptor1, "Device1", "0")
	// work := func() {}
	r := NewRobot("autorun",
		[]Connection{adaptor1},
		[]Device{driver1},
		// work,
	)

	go func() {
		assert.NoError(t, r.Start())
	}()

	time.Sleep(10 * time.Millisecond)
	assert.True(t, r.Running())

	// stop it
	assert.NoError(t, r.Stop())
	assert.False(t, r.Running())
}
