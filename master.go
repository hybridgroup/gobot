package gobot

import (
	"os"
	"os/signal"
	"runtime"
)

type Master struct {
	Robots []*Robot
	NumCPU int
}

func GobotMaster() *Master {
	m := new(Master)
	m.NumCPU = runtime.NumCPU()
	return m
}

var trap = func(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
}

func (m *Master) Start() {
	runtime.GOMAXPROCS(m.NumCPU)

	for s := range m.Robots {
		m.Robots[s].startRobot()
	}

	var c = make(chan os.Signal, 1)
	trap(c)

	for _ = range c {
		for r := range m.Robots {
			m.Robots[r].haltDevices()
			m.Robots[r].finalizeConnections()
		}
		break
	}
}

func (m *Master) FindRobot(name string) *Robot {
	for _, robot := range m.Robots {
		if robot.Name == name {
			return robot
		}
	}
	return nil
}

func (m *Master) FindRobotDevice(name string, device string) *device {
	robot := m.FindRobot(name)
	if robot != nil {
		return robot.GetDevice(device)
	}
	return nil
}

func (m *Master) FindRobotConnection(name string, connection string) *connection {
	robot := m.FindRobot(name)
	if robot != nil {
		return robot.GetConnection(connection)
	}
	return nil
}
