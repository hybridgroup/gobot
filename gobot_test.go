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
	g.Robots = []*Robot{
		NewTestRobot("Robot 1"),
		NewTestRobot("Robot 2"),
		NewTestRobot("Robot 3"),
	}
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
	Expect(t, g.Robot("Robot 1").Device("Device 4"), (*device)(nil))
	Expect(t, g.Robot("Robot 1").Device("Device 1").Name, "Device 1")
	Expect(t, len(g.Robot("Robot 1").Devices()), 3)
	Expect(t, g.Robot("Robot 1").Connection("Connection 4"), (*connection)(nil))
	Expect(t, len(g.Robot("Robot 1").Connections()), 3)
}
