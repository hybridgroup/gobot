package gobot

import "time"

type Gobot struct {
	Robots []Robot
}

func NewGobot() *Gobot {
	g := new(Gobot)
	return g
}

func (g *Gobot) Start() {
	for s := range g.Robots {
		go g.Robots[s].Start()
	}

	for {
		time.Sleep(10 * time.Millisecond)
	}
}

func (g *Gobot) FindRobot(name string) *Robot {
	for s := range g.Robots {
		if g.Robots[s].Name == name {
			return &g.Robots[s]
		}
	}
	return nil
}
func (g *Gobot) FindRobotDevice(name string, device string) *Device {
	for r := range g.Robots {
		if g.Robots[r].Name == name {
			for d := range g.Robots[r].devices {
				if g.Robots[r].devices[d].Name == device {
					return g.Robots[r].devices[d]
				}
			}
		}
	}
	return nil
}
func (g *Gobot) FindRobotConnection(name string, connection string) *Connection {
	for r := range g.Robots {
		if g.Robots[r].Name == name {
			for c := range g.Robots[r].connections {
				if g.Robots[r].connections[c].Name == connection {
					return g.Robots[r].connections[c]
				}
			}
		}
	}
	return nil
}
