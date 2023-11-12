package gobot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Len(t, json.Devices, r.Devices().Len())
	assert.Len(t, json.Commands, len(r.Commands()))
}

func TestRobotDevicesToJSON(t *testing.T) {
	r := newTestRobot("Robot99")
	json := NewJSONRobot(r)
	assert.Len(t, json.Devices, r.Devices().Len())
	assert.Equal(t, "Device1", json.Devices[0].Name)
	assert.Equal(t, "*gobot.testDriver", json.Devices[0].Driver)
	assert.Equal(t, "Connection1", json.Devices[0].Connection)
	assert.Len(t, json.Devices[0].Commands, 1)
}

func TestRobotStart(t *testing.T) {
	r := newTestRobot("Robot99")
	require.NoError(t, r.Start())
	require.NoError(t, r.Stop())
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
		require.NoError(t, r.Start())
	}()

	time.Sleep(10 * time.Millisecond)
	assert.True(t, r.Running())

	// stop it
	require.NoError(t, r.Stop())
	assert.False(t, r.Running())
}
