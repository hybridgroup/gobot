package gobot

import (
	"os"
	"os/signal"
	"runtime"
)

type Master struct {
	Robots []Robot
	NumCPU int
}

func GobotMaster() *Master {
	m := new(Master)
	m.NumCPU = runtime.NumCPU()
	return m
}

func (m *Master) Start() {
	runtime.GOMAXPROCS(m.NumCPU)

	for s := range m.Robots {
		go m.Robots[s].startRobot()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	for _ = range c {
		for r := range m.Robots {
			m.Robots[r].finalizeConnections()
		}
		break
	}
}

func (m *Master) FindRobot(name string) *Robot {
	for s := range m.Robots {
		if m.Robots[s].Name == name {
			return &m.Robots[s]
		}
	}
	return nil
}

func (m *Master) FindRobotDevice(name string, device string) *device {
	for r := range m.Robots {
		if m.Robots[r].Name == name {
			for d := range m.Robots[r].devices {
				if m.Robots[r].devices[d].Name == device {
					return m.Robots[r].devices[d]
				}
			}
		}
	}
	return nil
}

func (m *Master) FindRobotConnection(name string, connection string) *connection {
	for r := range m.Robots {
		if m.Robots[r].Name == name {
			for c := range m.Robots[r].connections {
				if m.Robots[r].connections[c].Name == connection {
					return m.Robots[r].connections[c]
				}
			}
		}
	}
	return nil
}
