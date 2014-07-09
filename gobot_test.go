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
	Expect(t, g.Robot("Robot 1").Name, "Robot 1")
	Expect(t, g.Robot("Robot 4"), (*Robot)(nil))
	Expect(t, g.Robot("Robot 1").Device("Device 4"), (Device)(nil))
	Expect(t, g.Robot("Robot 1").Device("Device 1").Name(), "Device 1")
	Expect(t, g.Robot("Robot 1").Devices().Len(), 3)
	Expect(t, g.Robot("Robot 1").Connection("Connection 4"), (Connection)(nil))
	Expect(t, g.Robot("Robot 1").Connections().Len(), 3)
}
