package gobot

import (
	"testing"
	"time"

	"gobot.io/x/gobot/gobottest"
)

func TestRobotConnectionEach(t *testing.T) {
	r := newTestRobot("Robot1")

	i := 0
	r.Connections().Each(func(conn Connection) {
		i++
	})
	gobottest.Assert(t, r.Connections().Len(), i)
}

func TestRobotToJSON(t *testing.T) {
	r := newTestRobot("Robot99")
	r.AddCommand("test_function", func(params map[string]interface{}) interface{} {
		return nil
	})
	json := NewJSONRobot(r)
	gobottest.Assert(t, len(json.Devices), r.Devices().Len())
	gobottest.Assert(t, len(json.Commands), len(r.Commands()))
}

func TestRobotDevicesToJSON(t *testing.T) {
	r := newTestRobot("Robot99")
	json := NewJSONRobot(r)
	gobottest.Assert(t, len(json.Devices), r.Devices().Len())
	gobottest.Assert(t, json.Devices[0].Name, "Device1")
	gobottest.Assert(t, json.Devices[0].Driver, "*gobot.testDriver")
	gobottest.Assert(t, json.Devices[0].Connection, "Connection1")
	gobottest.Assert(t, len(json.Devices[0].Commands), 1)
}

func TestRobotStart(t *testing.T) {
	r := newTestRobot("Robot99")
	gobottest.Assert(t, r.Start(), nil)
	gobottest.Assert(t, r.Stop(), nil)
	gobottest.Assert(t, r.Running(), false)
}

func TestRobotStartAutoRun(t *testing.T) {
	adaptor1 := newTestAdaptor("Connection1", "/dev/null")
	driver1 := newTestDriver(adaptor1, "Device1", "0")
	//work := func() {}
	r := NewRobot("autorun",
		[]Connection{adaptor1},
		[]Device{driver1},
		//work,
	)

	go func() {
		gobottest.Assert(t, r.Start(), nil)
	}()

	time.Sleep(10 * time.Millisecond)
	gobottest.Assert(t, r.Running(), true)

	// stop it
	gobottest.Assert(t, r.Stop(), nil)
	gobottest.Assert(t, r.Running(), false)
}
