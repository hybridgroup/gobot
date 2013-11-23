package gobot

import "time"

type Gobot struct {
	//  Master *Master
	Robots []Robot
}

func NewGobot() *Gobot {
	g := new(Gobot)
	//  g.Master = new(Master)
	return g
}

func (g *Gobot) Start() {
	//  g.Master.robots = g.Robots
	for s := range g.Robots {
		go g.Robots[s].Start()
	}

	for {
		time.Sleep(10 * time.Millisecond)
	}
}

func (m *Gobot) FindRobot(name string) *Robot {
	for s := range m.Robots {
		if m.Robots[s].Name == name {
			return &m.Robots[s]
		}
	}
	return nil
}
func (m *Gobot) FindRobotDevice(name string, device string) *Device {
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
func (m *Gobot) FindRobotConnection(name string, connection string) *Connection {
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
