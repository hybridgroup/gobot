package gobot

import (
	"log"
	"os"
	"testing"
)

func initTestGobot() *Gobot {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewGobot()
	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}
	g.AddRobot(NewTestRobot("Robot1"))
	g.AddRobot(NewTestRobot("Robot2"))
	g.AddRobot(NewTestRobot("Robot3"))
	return g
}

func TestGobotStart(t *testing.T) {
	g := initTestGobot()
	g.Start()
}

func TestGobotRobot(t *testing.T) {
	g := initTestGobot()
	Assert(t, g.Robot("Robot1").Name, "Robot1")
	Assert(t, g.Robot("Robot4"), (*Robot)(nil))
	Assert(t, g.Robot("Robot1").Device("Device4"), (Device)(nil))
	Assert(t, g.Robot("Robot1").Device("Device1").Name(), "Device1")
	Assert(t, g.Robot("Robot1").Devices().Len(), 3)
	Assert(t, g.Robot("Robot1").Connection("Connection4"), (Connection)(nil))
	Assert(t, g.Robot("Robot1").Connections().Len(), 3)
}

func TestGobotToJSON(t *testing.T) {
	g := initTestGobot()
	g.AddCommand("test_function", func(params map[string]interface{}) interface{} {
		return nil
	})
	json := g.ToJSON()
	Assert(t, len(json.Robots), g.Robots().Len())
	Assert(t, len(json.Commands), len(g.Commands()))
}
