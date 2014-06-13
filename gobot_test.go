package gobot

import (
	"log"
	"os"
	"testing"
)

var g *Gobot

func init() {
	log.SetOutput(new(Null))
	g = NewGobot()
	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}
	g.Robots = []*Robot{
		NewTestRobot("Robot 1"),
		NewTestRobot("Robot 2"),
		NewTestRobot("Robot 3"),
	}
}

func TestStart(t *testing.T) {
	g.Start()
}

func TestRobot(t *testing.T) {
	Expect(t, g.Robot("Robot 1").Name, "Robot 1")
	Expect(t, g.Robot("Robot 4"), (*Robot)(nil))
	Expect(t, g.Robot("Robot 1").Device("Device 4"), (*device)(nil))
	Expect(t, g.Robot("Robot 1").Device("Device 1").Name, "Device 1")
	Expect(t, len(g.Robot("Robot 1").Devices()), 3)
	Expect(t, g.Robot("Robot 1").Connection("Connection 4"), (*connection)(nil))
	Expect(t, len(g.Robot("Robot 1").Connections()), 3)
}
