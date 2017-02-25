package gobot

import (
	"testing"

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
