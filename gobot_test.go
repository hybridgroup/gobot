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
	g.AddRobot(NewTestRobot("Robot 1"))
	g.AddRobot(NewTestRobot("Robot 2"))
	g.AddRobot(NewTestRobot("Robot 3"))
	return g
}

func TestGobotStart(t *testing.T) {
	g := initTestGobot()
	g.Start()
}

func TestGobotRobot(t *testing.T) {
	g := initTestGobot()
	Assert(t, g.Robot("Robot 1").Name, "Robot 1")
	Assert(t, g.Robot("Robot 4"), (*Robot)(nil))
	Assert(t, g.Robot("Robot 1").Device("Device 4"), (Device)(nil))
	Assert(t, g.Robot("Robot 1").Device("Device 1").Name(), "Device 1")
	Assert(t, g.Robot("Robot 1").Devices().Len(), 3)
	Assert(t, g.Robot("Robot 1").Connection("Connection 4"), (Connection)(nil))
	Assert(t, g.Robot("Robot 1").Connections().Len(), 3)
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
